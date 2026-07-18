package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/bcpriok/pantas/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Worker struct {
	pool   *pgxpool.Pool
	cfg    config.Config
	client *http.Client
	log    *slog.Logger
}

type job struct {
	ID          string
	Channel     string
	Destination string
	Template    string
	Payload     map[string]any
	Attempts    int
}

type Execer interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

func New(pool *pgxpool.Pool, cfg config.Config, logger *slog.Logger) *Worker {
	return &Worker{
		pool:   pool,
		cfg:    cfg,
		client: &http.Client{Timeout: 15 * time.Second},
		log:    logger,
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.cfg.WorkerInterval)
	defer ticker.Stop()
	for {
		if err := w.process(ctx); err != nil && ctx.Err() == nil {
			w.log.Error("notification worker", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func Queue(ctx context.Context, database Execer, userID *string, channel, destination, template string, payload map[string]any) error {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = database.Exec(ctx, `
		insert into public.notification_jobs (user_id, channel, destination, template_code, payload)
		values ($1, $2, $3, $4, $5::jsonb)`, userID, channel, destination, template, string(encoded))
	return err
}

func (w *Worker) process(ctx context.Context) error {
	tx, err := w.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
		select id::text, channel, destination, template_code, payload, attempts
		from public.notification_jobs
		where status in ('pending', 'processing')
		  and next_attempt_at <= now()
		  and (status = 'pending' or locked_at < now() - interval '5 minutes')
		order by created_at
		for update skip locked
		limit 10`)
	if err != nil {
		return err
	}
	var jobs []job
	for rows.Next() {
		var item job
		var raw []byte
		if err := rows.Scan(&item.ID, &item.Channel, &item.Destination, &item.Template, &raw, &item.Attempts); err != nil {
			rows.Close()
			return err
		}
		if err := json.Unmarshal(raw, &item.Payload); err != nil {
			item.Payload = map[string]any{}
		}
		jobs = append(jobs, item)
	}
	rows.Close()
	for _, item := range jobs {
		if _, err := tx.Exec(ctx, `
			update public.notification_jobs
			set status = 'processing', locked_at = now(), attempts = attempts + 1
			where id = $1`, item.ID); err != nil {
			return err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	for _, item := range jobs {
		err := w.send(ctx, item)
		if err == nil {
			_, err = w.pool.Exec(ctx, `
				update public.notification_jobs
				set status = 'sent', sent_at = now(), locked_at = null, last_error = null
				where id = $1`, item.ID)
		} else {
			nextStatus := "pending"
			if item.Attempts+1 >= 5 {
				nextStatus = "failed"
			}
			delayMinutes := 1 << min(item.Attempts, 5)
			_, updateErr := w.pool.Exec(ctx, `
				update public.notification_jobs
				set status = $2, next_attempt_at = now() + make_interval(mins => $3),
				    locked_at = null, last_error = left($4, 2000)
				where id = $1`, item.ID, nextStatus, delayMinutes, err.Error())
			if updateErr != nil {
				w.log.Error("update failed notification job", "job_id", item.ID, "error", updateErr)
			}
		}
	}
	return nil
}

func (w *Worker) send(ctx context.Context, item job) error {
	switch item.Channel {
	case "email":
		return w.sendEmail(ctx, item)
	case "phone":
		return w.sendPhone(ctx, item)
	default:
		return fmt.Errorf("unsupported notification channel %q", item.Channel)
	}
}

func (w *Worker) sendEmail(ctx context.Context, item job) error {
	if w.cfg.ResendAPIKey == "" || w.cfg.EmailFrom == "" {
		return fmt.Errorf("email provider belum dikonfigurasi")
	}
	subject, body := renderTemplate(item.Template, item.Payload, w.cfg.AppURL)
	payload := map[string]any{
		"from":    w.cfg.EmailFrom,
		"to":      []string{item.Destination},
		"subject": subject,
		"html":    body,
	}
	encoded, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(encoded))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+w.cfg.ResendAPIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2000))
		return fmt.Errorf("email provider status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

func (w *Worker) sendPhone(ctx context.Context, item job) error {
	if w.cfg.PhoneWebhookURL == "" {
		return fmt.Errorf("webhook OTP nomor HP belum dikonfigurasi")
	}
	_, body := renderTemplate(item.Template, item.Payload, w.cfg.AppURL)
	payload, _ := json.Marshal(map[string]any{
		"to":       item.Destination,
		"message":  stripHTML(body),
		"template": item.Template,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.cfg.PhoneWebhookURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if w.cfg.PhoneWebhookToken != "" {
		req.Header.Set("Authorization", "Bearer "+w.cfg.PhoneWebhookToken)
	}
	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		response, _ := io.ReadAll(io.LimitReader(resp.Body, 2000))
		return fmt.Errorf("phone webhook status %d: %s", resp.StatusCode, strings.TrimSpace(string(response)))
	}
	return nil
}

func renderTemplate(code string, payload map[string]any, appURL string) (string, string) {
	name := html.EscapeString(value(payload, "name", "Pegawai"))
	period := html.EscapeString(value(payload, "period", "periode terbaru"))
	otp := html.EscapeString(value(payload, "otp", ""))
	link := html.EscapeString(appURL)
	switch code {
	case "password_otp":
		return "Kode reset password PANTAS", fmt.Sprintf("Halo %s, kode reset password Anda adalah <strong>%s</strong>. Kode berlaku 10 menit. Jangan berikan kode ini kepada siapa pun.", name, otp)
	case "contact_otp":
		return "Kode verifikasi kontak PANTAS", fmt.Sprintf("Halo %s, kode verifikasi kontak Anda adalah <strong>%s</strong>. Kode berlaku 10 menit.", name, otp)
	case "appeal_submitted":
		return "Banding PANTAS menunggu verifikasi", fmt.Sprintf("Banding periode %s telah diajukan dan menunggu verifikasi. Buka <a href=\"%s\">PANTAS</a> untuk melihat detail.", period, link)
	case "appeal_reviewed":
		return "Status banding PANTAS diperbarui", fmt.Sprintf("Status banding periode %s telah diperbarui. Buka <a href=\"%s\">PANTAS</a> untuk melihat hasilnya.", period, link)
	default:
		return "Data potongan PANTAS telah diperbarui", fmt.Sprintf("Halo %s, data presensi dan potongan periode %s telah tersedia. Demi privasi, detail hanya dapat dilihat setelah login di <a href=\"%s\">PANTAS</a>.", name, period, link)
	}
}

func value(payload map[string]any, key, fallback string) string {
	if text, ok := payload[key].(string); ok && text != "" {
		return text
	}
	return fallback
}

func stripHTML(value string) string {
	replacer := strings.NewReplacer("<strong>", "", "</strong>", "", "<a href=\"", "", "</a>", "", "\">", " ")
	return html.UnescapeString(replacer.Replace(value))
}
