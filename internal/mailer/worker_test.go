package mailer

import (
	"context"
	"encoding/json"
	"testing"

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
