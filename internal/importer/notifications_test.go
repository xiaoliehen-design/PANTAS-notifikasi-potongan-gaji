package importer

import (
	"context"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

type periodNotificationRecorder struct {
	query     string
	arguments []any
}

func (r *periodNotificationRecorder) Exec(_ context.Context, query string, arguments ...any) (pgconn.CommandTag, error) {
	r.query = query
	r.arguments = arguments
	return pgconn.CommandTag{}, nil
}

func TestQueuePublishedPeriodJobsTargetsBothVerifiedContacts(t *testing.T) {
	database := &periodNotificationRecorder{}
	if err := queuePublishedPeriodJobs(context.Background(), database, "Juli 2026", "batch-123"); err != nil {
		t.Fatalf("queuePublishedPeriodJobs() error = %v", err)
	}

	for _, expected := range []string{
		"'email'::text",
		"u.email_verified_at",
		"'phone'::text",
		"u.phone_verified_at",
		"contact.destination is not null",
		"contact.verified_at is not null",
		"u.is_active",
		"u.deleted_at is null",
		"'period_published'",
	} {
		if !strings.Contains(database.query, expected) {
			t.Errorf("notification query does not contain %q", expected)
		}
	}
	if len(database.arguments) != 2 || database.arguments[0] != "Juli 2026" || database.arguments[1] != "batch-123" {
		t.Fatalf("arguments = %#v, want period label and batch ID", database.arguments)
	}
}
