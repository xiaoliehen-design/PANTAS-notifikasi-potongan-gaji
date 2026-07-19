package mailer

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log/slog"
	"mime"
	"mime/quotedprintable"
	"net"
	"net/http"
	mailaddress "net/mail"
	"net/smtp"
	"net/url"
	"strconv"
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

type QueryExecer interface {
	Execer
	QueryRow(context.Context, string, ...any) pgx.Row
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

// QueueImmediate reserves a job for delivery by the request that created it.
// If the process stops before delivery, the background worker safely retries it
// after the processing lock expires.
func QueueImmediate(ctx context.Context, database QueryExecer, userID *string, channel, destination, template string, payload map[string]any) (string, error) {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	var id string
	err = database.QueryRow(ctx, `
		insert into public.notification_jobs
		  (user_id, channel, destination, template_code, payload, status, attempts, locked_at)
		values ($1, $2, $3, $4, $5::jsonb, 'processing', 1, now())
		returning id::text`, userID, channel, destination, template, string(encoded)).Scan(&id)
	return id, err
}

func (w *Worker) ChannelConfigured(channel string) bool {
	switch channel {
	case "email":
		return w.emailProvider() != ""
	case "phone":
		return w.phoneProvider() != ""
	default:
		return false
	}
}

func (w *Worker) emailProvider() string {
	switch strings.ToLower(w.cfg.EmailProvider) {
	case "smtp":
		if w.cfg.SMTPHost != "" && w.cfg.SMTPPort != "" && w.cfg.SMTPUsername != "" && w.cfg.SMTPPassword != "" && w.cfg.EmailFrom != "" {
			return "smtp"
		}
	case "resend":
		if w.cfg.ResendAPIKey != "" && w.cfg.ResendAPIURL != "" && w.cfg.EmailFrom != "" {
			return "resend"
		}
	case "", "auto":
		if w.cfg.SMTPHost != "" && w.cfg.SMTPPort != "" && w.cfg.SMTPUsername != "" && w.cfg.SMTPPassword != "" && w.cfg.EmailFrom != "" {
			return "smtp"
		}
		if w.cfg.ResendAPIKey != "" && w.cfg.ResendAPIURL != "" && w.cfg.EmailFrom != "" {
			return "resend"
		}
	}
	return ""
}

func (w *Worker) phoneProvider() string {
	switch strings.ToLower(w.cfg.PhoneProvider) {
	case "twilio":
		if w.twilioConfigured() {
			return "twilio"
		}
	case "webhook":
		if w.cfg.PhoneWebhookURL != "" {
			return "webhook"
		}
	case "", "auto":
		if w.twilioConfigured() {
			return "twilio"
		}
		if w.cfg.PhoneWebhookURL != "" {
			return "webhook"
		}
	}
	return ""
}

func (w *Worker) twilioConfigured() bool {
	credentialOK := (w.cfg.TwilioAPIKey != "" && w.cfg.TwilioAPISecret != "") || w.cfg.TwilioAuthToken != ""
	return w.cfg.TwilioAccountSID != "" && credentialOK && (w.cfg.TwilioMessagingSID != "" || w.cfg.TwilioFrom != "")
}

// DeliverClaimed sends a job created with QueueImmediate. Provider failures are
// retained in notification_jobs and scheduled for the normal retry worker.
func (w *Worker) DeliverClaimed(ctx context.Context, id string) error {
	var item job
	var raw []byte
	err := w.pool.QueryRow(ctx, `
		select id::text, channel, destination, template_code, payload, attempts
		from public.notification_jobs
		where id = $1 and status = 'processing'`, id).Scan(
		&item.ID, &item.Channel, &item.Destination, &item.Template, &raw, &item.Attempts,
	)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, &item.Payload); err != nil {
		item.Payload = map[string]any{}
	}

	deliveryErr := w.send(ctx, item)
	if deliveryErr == nil {
		_, err = w.pool.Exec(ctx, `
			update public.notification_jobs
			set status = 'sent', sent_at = now(), locked_at = null, last_error = null
			where id = $1`, item.ID)
		if err != nil {
			w.log.Error("record immediate notification success", "job_id", item.ID, "error", err)
		}
		return nil
	}

	nextStatus := "pending"
	if item.Attempts >= 5 {
		nextStatus = "failed"
	}
	delayMinutes := 1 << min(max(item.Attempts-1, 0), 5)
	if _, err := w.pool.Exec(ctx, `
		update public.notification_jobs
		set status = $2, next_attempt_at = now() + make_interval(mins => $3),
		    locked_at = null, last_error = left($4, 2000)
		where id = $1`, item.ID, nextStatus, delayMinutes, deliveryErr.Error()); err != nil {
		w.log.Error("record immediate notification failure", "job_id", item.ID, "error", err)
	}
	w.log.Error("immediate notification delivery", "job_id", item.ID, "channel", item.Channel, "error", deliveryErr)
	return deliveryErr
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
	switch w.emailProvider() {
	case "smtp":
		return w.sendEmailSMTP(ctx, item)
	case "resend":
		return w.sendEmailResend(ctx, item)
	default:
		return fmt.Errorf("email provider belum dikonfigurasi")
	}
}

func (w *Worker) sendEmailResend(ctx context.Context, item job) error {
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
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.cfg.ResendAPIURL, bytes.NewReader(encoded))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+w.cfg.ResendAPIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "PANTAS/1.0")
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

func (w *Worker) sendEmailSMTP(ctx context.Context, item job) error {
	from, err := mailaddress.ParseAddress(w.cfg.EmailFrom)
	if err != nil {
		return fmt.Errorf("EMAIL_FROM tidak valid: %w", err)
	}
	port, err := strconv.Atoi(w.cfg.SMTPPort)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("SMTP_PORT tidak valid")
	}
	address := net.JoinHostPort(w.cfg.SMTPHost, strconv.Itoa(port))
	dialer := &net.Dialer{Timeout: 15 * time.Second}
	tlsConfig := &tls.Config{ServerName: w.cfg.SMTPHost, MinVersion: tls.VersionTLS12}

	var connection net.Conn
	if w.cfg.SMTPTLSMode == "implicit" {
		connection, err = tls.DialWithDialer(dialer, "tcp", address, tlsConfig)
	} else {
		connection, err = dialer.DialContext(ctx, "tcp", address)
	}
	if err != nil {
		return fmt.Errorf("koneksi SMTP gagal: %w", err)
	}
	defer connection.Close()
	deadline := time.Now().Add(20 * time.Second)
	if contextDeadline, ok := ctx.Deadline(); ok && contextDeadline.Before(deadline) {
		deadline = contextDeadline
	}
	_ = connection.SetDeadline(deadline)

	client, err := smtp.NewClient(connection, w.cfg.SMTPHost)
	if err != nil {
		return fmt.Errorf("inisialisasi SMTP gagal: %w", err)
	}
	defer client.Close()
	if w.cfg.SMTPTLSMode != "implicit" {
		if ok, _ := client.Extension("STARTTLS"); !ok {
			return fmt.Errorf("server SMTP tidak mendukung STARTTLS")
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("STARTTLS SMTP gagal: %w", err)
		}
	}
	if err := client.Auth(smtp.PlainAuth("", w.cfg.SMTPUsername, w.cfg.SMTPPassword, w.cfg.SMTPHost)); err != nil {
		return fmt.Errorf("autentikasi SMTP gagal: %w", err)
	}
	if err := client.Mail(from.Address); err != nil {
		return fmt.Errorf("alamat pengirim SMTP ditolak: %w", err)
	}
	if err := client.Rcpt(item.Destination); err != nil {
		return fmt.Errorf("alamat penerima SMTP ditolak: %w", err)
	}
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA ditolak: %w", err)
	}
	subject, body := renderTemplate(item.Template, item.Payload, w.cfg.AppURL)
	message, err := smtpMessage(from, item.Destination, subject, body)
	if err != nil {
		_ = writer.Close()
		return err
	}
	if _, err := writer.Write(message); err != nil {
		_ = writer.Close()
		return fmt.Errorf("penulisan email SMTP gagal: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("pengiriman email SMTP gagal: %w", err)
	}
	if err := client.Quit(); err != nil {
		return fmt.Errorf("penutupan SMTP gagal: %w", err)
	}
	return nil
}

func smtpMessage(from *mailaddress.Address, destination, subject, body string) ([]byte, error) {
	to := (&mailaddress.Address{Address: destination}).String()
	var message bytes.Buffer
	fmt.Fprintf(&message, "From: %s\r\n", from.String())
	fmt.Fprintf(&message, "To: %s\r\n", to)
	fmt.Fprintf(&message, "Subject: %s\r\n", mime.QEncoding.Encode("UTF-8", subject))
	fmt.Fprintf(&message, "Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	message.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
	encoded := quotedprintable.NewWriter(&message)
	if _, err := encoded.Write([]byte(body)); err != nil {
		return nil, err
	}
	if err := encoded.Close(); err != nil {
		return nil, err
	}
	return message.Bytes(), nil
}

func (w *Worker) sendPhone(ctx context.Context, item job) error {
	switch w.phoneProvider() {
	case "twilio":
		return w.sendPhoneTwilio(ctx, item)
	case "webhook":
		return w.sendPhoneWebhook(ctx, item)
	default:
		return fmt.Errorf("provider OTP nomor HP belum dikonfigurasi")
	}
}

func (w *Worker) sendPhoneWebhook(ctx context.Context, item job) error {
	if w.cfg.PhoneWebhookURL == "" {
		return fmt.Errorf("webhook OTP nomor HP belum dikonfigurasi")
	}
	payload, _ := json.Marshal(map[string]any{
		"to":       item.Destination,
		"message":  renderPhoneTemplate(item.Template, item.Payload, w.cfg.AppURL),
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

func (w *Worker) sendPhoneTwilio(ctx context.Context, item job) error {
	if !w.twilioConfigured() {
		return fmt.Errorf("Twilio SMS belum dikonfigurasi")
	}
	form := url.Values{}
	form.Set("To", item.Destination)
	form.Set("Body", renderPhoneTemplate(item.Template, item.Payload, w.cfg.AppURL))
	if w.cfg.TwilioMessagingSID != "" {
		form.Set("MessagingServiceSid", w.cfg.TwilioMessagingSID)
	} else {
		form.Set("From", w.cfg.TwilioFrom)
	}
	endpoint := fmt.Sprintf("%s/Accounts/%s/Messages.json", w.cfg.TwilioAPIBaseURL, url.PathEscape(w.cfg.TwilioAccountSID))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	username, password := w.cfg.TwilioAPIKey, w.cfg.TwilioAPISecret
	if username == "" {
		username, password = w.cfg.TwilioAccountSID, w.cfg.TwilioAuthToken
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "PANTAS/1.0")
	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	responseBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2000))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Twilio status %d: %s", resp.StatusCode, strings.TrimSpace(string(responseBody)))
	}
	var result struct {
		SID    string `json:"sid"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return fmt.Errorf("respons Twilio tidak valid: %w", err)
	}
	if result.SID == "" {
		return fmt.Errorf("Twilio tidak mengembalikan SID pesan")
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

func renderPhoneTemplate(code string, payload map[string]any, appURL string) string {
	otp := value(payload, "otp", "")
	switch code {
	case "password_otp":
		return fmt.Sprintf("PANTAS: kode reset password %s. Berlaku 10 menit. Jangan bagikan kode ini.", otp)
	case "contact_otp":
		return fmt.Sprintf("PANTAS: kode verifikasi nomor HP %s. Berlaku 10 menit. Jangan bagikan kode ini.", otp)
	default:
		_, body := renderTemplate(code, payload, appURL)
		return stripHTML(body)
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
