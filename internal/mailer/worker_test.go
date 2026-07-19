package mailer

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/mail"
	"strings"
	"testing"

	"github.com/bcpriok/pantas/internal/config"
	"github.com/jackc/pgx/v5/pgconn"
)

type recordingExecer struct {
	arguments []any
}

func (r *recordingExecer) Exec(_ context.Context, _ string, arguments ...any) (pgconn.CommandTag, error) {
	r.arguments = arguments
	return pgconn.CommandTag{}, nil
}

func TestQueueSendsJSONAsText(t *testing.T) {
	database := &recordingExecer{}
	userID := "82a18cff-5f52-4e48-9a07-82536db96aa4"
	if err := Queue(context.Background(), database, &userID, "email", "pegawai@example.go.id", "contact_otp", map[string]any{
		"name": "Pegawai PANTAS",
		"otp":  "123456",
	}); err != nil {
		t.Fatalf("Queue() error = %v", err)
	}
	if len(database.arguments) != 5 {
		t.Fatalf("argument count = %d, want 5", len(database.arguments))
	}
	payload, ok := database.arguments[4].(string)
	if !ok {
		t.Fatalf("payload type = %T, want string", database.arguments[4])
	}
	var decoded map[string]any
	if err := json.Unmarshal([]byte(payload), &decoded); err != nil {
		t.Fatalf("payload is not valid JSON: %v", err)
	}
	if decoded["otp"] != "123456" {
		t.Fatalf("payload otp = %v, want 123456", decoded["otp"])
	}
}

func TestChannelConfigured(t *testing.T) {
	worker := New(nil, config.Config{
		EmailProvider:   "resend",
		ResendAPIKey:    "re_test",
		ResendAPIURL:    "https://api.resend.com/emails",
		EmailFrom:       "PANTAS <noreply@example.go.id>",
		PhoneProvider:   "webhook",
		PhoneWebhookURL: "https://provider.example/otp",
	}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if !worker.ChannelConfigured("email") || !worker.ChannelConfigured("phone") {
		t.Fatal("configured delivery channel reported unavailable")
	}
	if worker.ChannelConfigured("unknown") {
		t.Fatal("unknown channel reported available")
	}
}

func TestChannelConfiguredForSMTPAndTwilio(t *testing.T) {
	worker := New(nil, config.Config{
		EmailProvider:      "smtp",
		EmailFrom:          "PANTAS <pantas@example.com>",
		SMTPHost:           "smtp.example.com",
		SMTPPort:           "587",
		SMTPUsername:       "pantas@example.com",
		SMTPPassword:       "app-password",
		PhoneProvider:      "twilio",
		TwilioAccountSID:   "AC00000000000000000000000000000000",
		TwilioAPIKey:       "SK00000000000000000000000000000000",
		TwilioAPISecret:    "secret",
		TwilioMessagingSID: "MG00000000000000000000000000000000",
	}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if !worker.ChannelConfigured("email") || !worker.ChannelConfigured("phone") {
		t.Fatal("SMTP/Twilio delivery channel reported unavailable")
	}
}

func TestSendEmailUsesConfiguredProvider(t *testing.T) {
	var received map[string]any
	provider := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Authorization") != "Bearer re_test" {
			t.Errorf("Authorization = %q", request.Header.Get("Authorization"))
		}
		if err := json.NewDecoder(request.Body).Decode(&received); err != nil {
			t.Errorf("decode provider payload: %v", err)
		}
		response.WriteHeader(http.StatusOK)
	}))
	defer provider.Close()

	worker := New(nil, config.Config{
		EmailProvider: "resend",
		ResendAPIKey:  "re_test",
		ResendAPIURL:  provider.URL,
		EmailFrom:     "PANTAS <noreply@example.go.id>",
		AppURL:        "https://pantas.example.go.id",
	}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	err := worker.sendEmail(context.Background(), job{
		Channel: "email", Destination: "pegawai@example.go.id", Template: "contact_otp",
		Payload: map[string]any{"name": "Pegawai PANTAS", "otp": "123456"},
	})
	if err != nil {
		t.Fatalf("sendEmail() error = %v", err)
	}
	if received["subject"] != "Kode verifikasi kontak PANTAS" {
		t.Fatalf("subject = %v", received["subject"])
	}
}

func TestSendEmailUsesBrevoHTTPSProvider(t *testing.T) {
	var received struct {
		Sender struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"sender"`
		To []struct {
			Email string `json:"email"`
		} `json:"to"`
		Subject     string `json:"subject"`
		HTMLContent string `json:"htmlContent"`
	}
	provider := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", request.Method)
		}
		if request.Header.Get("api-key") != "brevo-test-key" {
			t.Error("Brevo API key header is missing")
		}
		if err := json.NewDecoder(request.Body).Decode(&received); err != nil {
			t.Errorf("decode Brevo payload: %v", err)
		}
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusCreated)
		_, _ = response.Write([]byte(`{"messageId":"test-message-id"}`))
	}))
	defer provider.Close()

	worker := New(nil, config.Config{
		EmailProvider: "brevo",
		BrevoAPIKey:   "brevo-test-key",
		BrevoAPIURL:   provider.URL,
		EmailFrom:     "PANTAS <noreply@example.go.id>",
		AppURL:        "https://pantas.example.go.id",
	}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	err := worker.sendEmail(context.Background(), job{
		Channel: "email", Destination: "pegawai@example.go.id", Template: "contact_otp",
		Payload: map[string]any{"name": "Pegawai PANTAS", "otp": "123456"},
	})
	if err != nil {
		t.Fatalf("sendEmail() error = %v", err)
	}
	if received.Sender.Name != "PANTAS" || received.Sender.Email != "noreply@example.go.id" {
		t.Fatalf("sender = %#v", received.Sender)
	}
	if len(received.To) != 1 || received.To[0].Email != "pegawai@example.go.id" {
		t.Fatalf("recipients = %#v", received.To)
	}
	if received.Subject != "Kode verifikasi kontak PANTAS" || !strings.Contains(received.HTMLContent, "123456") {
		t.Fatalf("unexpected message: subject=%q body=%q", received.Subject, received.HTMLContent)
	}
}

func TestPublicDeliveryErrorExplainsRenderSMTPBlockWithoutLeakingDetail(t *testing.T) {
	raw := "koneksi SMTP gagal: dial tcp 192.0.2.1:587: i/o timeout"
	message := PublicDeliveryError(errors.New(raw))
	if !strings.Contains(message, "Render Free") || !strings.Contains(message, "EMAIL_PROVIDER=brevo") {
		t.Fatalf("message is not actionable: %q", message)
	}
	if strings.Contains(message, "192.0.2.1") {
		t.Fatalf("message leaked provider detail: %q", message)
	}
}

func TestSMTPMessageContainsSafeMIMEHeaders(t *testing.T) {
	from, err := mailAddress("PANTAS <pantas@example.com>")
	if err != nil {
		t.Fatal(err)
	}
	message, err := smtpMessage(from, "pegawai@example.com", "Kode verifikasi PANTAS", "Kode Anda <strong>123456</strong>")
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := mail.ReadMessage(strings.NewReader(string(message)))
	if err != nil {
		t.Fatalf("parse SMTP message: %v", err)
	}
	parsedFrom, err := mail.ParseAddress(parsed.Header.Get("From"))
	if err != nil || parsedFrom.Name != "PANTAS" || parsedFrom.Address != "pantas@example.com" {
		t.Fatalf("From header = %q, parsed=%#v, error=%v", parsed.Header.Get("From"), parsedFrom, err)
	}
	parsedTo, err := mail.ParseAddress(parsed.Header.Get("To"))
	if err != nil || parsedTo.Address != "pegawai@example.com" {
		t.Fatalf("To header = %q, parsed=%#v, error=%v", parsed.Header.Get("To"), parsedTo, err)
	}
	if !strings.HasPrefix(parsed.Header.Get("Content-Type"), "text/html") {
		t.Fatalf("Content-Type = %q, want text/html", parsed.Header.Get("Content-Type"))
	}
	body, err := io.ReadAll(parsed.Body)
	if err != nil {
		t.Fatalf("read SMTP body: %v", err)
	}
	if !strings.Contains(string(body), "123456") {
		t.Fatalf("SMTP body does not contain OTP: %q", string(body))
	}
}

func TestSendPhoneUsesTwilio(t *testing.T) {
	var receivedTo, receivedBody, receivedService string
	provider := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		username, password, ok := request.BasicAuth()
		if !ok || username != "SK_test" || password != "api-secret" {
			t.Errorf("unexpected Twilio credentials: %q %q", username, password)
		}
		if err := request.ParseForm(); err != nil {
			t.Errorf("parse Twilio form: %v", err)
		}
		receivedTo = request.Form.Get("To")
		receivedBody = request.Form.Get("Body")
		receivedService = request.Form.Get("MessagingServiceSid")
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusCreated)
		_, _ = response.Write([]byte(`{"sid":"SM00000000000000000000000000000000","status":"queued"}`))
	}))
	defer provider.Close()

	worker := New(nil, config.Config{
		PhoneProvider:      "twilio",
		TwilioAPIBaseURL:   provider.URL,
		TwilioAccountSID:   "AC00000000000000000000000000000000",
		TwilioAPIKey:       "SK_test",
		TwilioAPISecret:    "api-secret",
		TwilioMessagingSID: "MG00000000000000000000000000000000",
	}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	err := worker.sendPhone(context.Background(), job{
		Channel: "phone", Destination: "+6281234567890", Template: "contact_otp",
		Payload: map[string]any{"name": "Pegawai PANTAS", "otp": "123456"},
	})
	if err != nil {
		t.Fatalf("sendPhone() error = %v", err)
	}
	if receivedTo != "+6281234567890" || receivedService == "" || !strings.Contains(receivedBody, "123456") {
		t.Fatalf("unexpected Twilio form: to=%q service=%q body=%q", receivedTo, receivedService, receivedBody)
	}
}

func TestPeriodPublishedTemplatesContainPeriodAndSafeLink(t *testing.T) {
	payload := map[string]any{"name": "Pegawai PANTAS", "period": "Juli 2026"}
	appURL := "https://pantas-notifikasi-potongan-gaji.onrender.com"
	subject, emailBody := renderTemplate("period_published", payload, appURL)
	phoneBody := renderPhoneTemplate("period_published", payload, appURL)

	for label, text := range map[string]string{
		"subject": subject,
		"email":   emailBody,
		"phone":   phoneBody,
	} {
		if !strings.Contains(text, "Juli 2026") {
			t.Errorf("%s template does not contain period: %q", label, text)
		}
	}
	if !strings.Contains(emailBody, appURL) || !strings.Contains(phoneBody, appURL) {
		t.Fatal("period-published templates do not contain the configured application URL")
	}
	if strings.Contains(phoneBody, "<") || strings.Contains(phoneBody, ">") {
		t.Fatalf("phone template contains HTML: %q", phoneBody)
	}
	if len(phoneBody) > 160 {
		t.Fatalf("phone template length = %d, want at most one GSM-7 SMS segment", len(phoneBody))
	}
}

func mailAddress(value string) (*mail.Address, error) {
	return mail.ParseAddress(value)
}
