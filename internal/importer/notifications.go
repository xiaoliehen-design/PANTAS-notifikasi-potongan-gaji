package importer

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
)

type notificationExecer interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

// queuePublishedPeriodJobs creates one delivery job for every verified contact
// channel. Keeping both channels in one INSERT makes the publish transaction
// all-or-nothing: the period cannot become visible without its notification
// jobs also being recorded.
func queuePublishedPeriodJobs(ctx context.Context, database notificationExecer, periodLabel, batchID string) error {
	_, err := database.Exec(ctx, `
		insert into public.notification_jobs
			(user_id, channel, destination, template_code, payload)
		select u.id, contact.channel, contact.destination, 'period_published',
		       jsonb_build_object(
		           'name', u.name,
		           'period', $1,
		           'batch_id', $2
		       )
		from public.users u
		cross join lateral (
			values
				('email'::text, nullif(btrim(u.email), ''), u.email_verified_at),
				('phone'::text, nullif(btrim(u.phone_e164), ''), u.phone_verified_at)
		) as contact(channel, destination, verified_at)
		where u.is_active
		  and u.deleted_at is null
		  and contact.destination is not null
		  and contact.verified_at is not null`, periodLabel, batchID)
	return err
}
