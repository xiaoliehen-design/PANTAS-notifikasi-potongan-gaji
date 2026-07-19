package mailer

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
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
		ResendAPIKey:    "re_test",
		ResendAPIURL:    "https://api.resend.com/emails",
		EmailFrom:       "PANTAS <noreply@example.go.id>",
		PhoneWebhookURL: "https://provider.example/otp",
	}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if !worker.ChannelConfigured("email") || !worker.ChannelConfigured("phone") {
		t.Fatal("configured delivery channel reported unavailable")
	}
	if worker.ChannelConfigured("unknown") {
		t.Fatal("unknown channel reported available")
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
		ResendAPIKey: "re_test",
		ResendAPIURL: provider.URL,
		EmailFrom:    "PANTAS <noreply@example.go.id>",
		AppURL:       "https://pantas.example.go.id",
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
