package importer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bcpriok/pantas/internal/auth"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrMissingUsers     = errors.New("terdapat NIP yang belum terdaftar")
	ErrAppealsExist     = errors.New("periode tidak dapat diganti karena sudah memiliki banding")
	ErrAlreadyPublished = errors.New("file yang sama sudah pernah dipublikasikan")
)

type Rule struct {
	SourceField string
	Code        string
	Label       string
	Rate        float64
}

type userMatch struct {
	ID         string
	Name       string
	UnitSource string
}

type MissingUser struct {
	NIP       string `json:"nip"`
	Name      string `json:"name"`
	Placement string `json:"placement"`
}

type UnitMismatch struct {
	NIP              string `json:"nip"`
	Name             string `json:"name"`
	ExcelPlacement   string `json:"excel_placement"`
	CurrentPlacement string `json:"current_placement"`
}

type Preview struct {
	BatchID          string         `json:"batch_id,omitempty"`
	PeriodID         string         `json:"period_id,omitempty"`
	PeriodLabel      string         `json:"period_label"`
	PeriodStart      string         `json:"period_start"`
	PeriodEnd        string         `json:"period_end"`
	SheetName        string         `json:"sheet_name"`
	IntegrityStatus  string         `json:"integrity_status"`
	Rows             int            `json:"rows"`
	BlankRowsIgnored int            `json:"blank_rows_ignored"`
	Employees        int            `json:"employees"`
	DeductionDays    int            `json:"deduction_days"`
	TotalDeduction   float64        `json:"total_deduction"`
	MissingUsers     []MissingUser  `json:"missing_users"`
	UnitMismatches   []UnitMismatch `json:"unit_mismatches"`
	Warnings         []string       `json:"warnings"`
	ReadyToPublish   bool           `json:"ready_to_publish"`
}

type Service struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (s *Service) PreviewAndStage(ctx context.Context, actor auth.Principal, filename string, data []byte) (Preview, error) {
	parsed, err := ParseXLSX(data)
	if err != nil {
		return Preview{}, err
	}
	users, err := s.loadUsers(ctx)
	if err != nil {
		return Preview{}, err
	}
	rules, err := s.loadRules(ctx)
	if err != nil {
		return Preview{}, err
	}
	preview := Preview{
		PeriodLabel:      parsed.PeriodLabel,
		PeriodStart:      parsed.PeriodStart.Format("2006-01-02"),
		PeriodEnd:        parsed.PeriodEnd.Format("2006-01-02"),
		SheetName:        parsed.SheetName,
		IntegrityStatus:  parsed.IntegrityStatus,
		Rows:             len(parsed.Records),
		BlankRowsIgnored: parsed.BlankRows,
		Warnings:         append([]string(nil), parsed.Warnings...),
	}
	employeeSet := map[string]struct{}{}
	missing := map[string]MissingUser{}
	mismatch := map[string]UnitMismatch{}
	for index := range parsed.Records {
		record := &parsed.Records[index]
		employeeSet[record.NIP] = struct{}{}
		matched, ok := users[record.NIP]
		if !ok {
			missing[record.NIP] = MissingUser{NIP: record.NIP, Name: record.Name, Placement: record.SourcePlacement}
			continue
		}
		if matched.Name != record.Name {
			preview.Warnings = appendUnique(preview.Warnings, fmt.Sprintf("Nama NIP %s berbeda: aplikasi=%q, Excel=%q.", record.NIP, matched.Name, record.Name))
		}
		if normalizeUnit(matched.UnitSource) != normalizeUnit(record.SourcePlacement) {
			mismatch[record.NIP] = UnitMismatch{
				NIP: record.NIP, Name: record.Name,
				ExcelPlacement: record.SourcePlacement, CurrentPlacement: matched.UnitSource,
			}
		}
		applyRules(record, rules)
		if record.DeductionRate > 0 {
			preview.DeductionDays++
			preview.TotalDeduction += record.DeductionRate
		}
	}
	preview.Employees = len(employeeSet)
	for _, item := range missing {
		preview.MissingUsers = append(preview.MissingUsers, item)
	}
	for _, item := range mismatch {
		preview.UnitMismatches = append(preview.UnitMismatches, item)
	}
	sort.Slice(preview.MissingUsers, func(i, j int) bool { return preview.MissingUsers[i].Name < preview.MissingUsers[j].Name })
	sort.Slice(preview.UnitMismatches, func(i, j int) bool { return preview.UnitMismatches[i].Name < preview.UnitMismatches[j].Name })
	if len(preview.MissingUsers) > 0 {
		preview.ReadyToPublish = false
		return preview, ErrMissingUsers
	}
	if len(preview.UnitMismatches) > 0 {
		preview.Warnings = append(preview.Warnings, fmt.Sprintf("%d pegawai memiliki penempatan Excel yang berbeda dari profil aplikasi; data tetap dapat dipratinjau, tetapi admin sebaiknya memeriksa mutasi.", len(preview.UnitMismatches)))
	}

	batchID, periodID, err := s.stage(ctx, actor, filename, parsed, users, preview)
	if err != nil {
		return Preview{}, err
	}
	preview.BatchID = batchID
	preview.PeriodID = periodID
	preview.ReadyToPublish = true
	return preview, nil
}

func (s *Service) Publish(ctx context.Context, actor auth.Principal, batchID, ip string) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var periodID, label, status, hash string
	var currentBatchID *string
	err = tx.QueryRow(ctx, `
		select ib.period_id::text, rp.label, ib.status, ib.file_sha256, rp.published_batch_id::text
		from public.import_batches ib
		join public.reporting_periods rp on rp.id = ib.period_id
		where ib.id = $1
		for update of ib, rp`, batchID).Scan(&periodID, &label, &status, &hash, &currentBatchID)
	if errors.Is(err, pgx.ErrNoRows) {
		return errors.New("batch import tidak ditemukan")
	}
	if err != nil {
		return err
	}
	if status != "draft" {
		return errors.New("batch import tidak lagi berstatus draft")
	}
	var duplicate bool
	if err := tx.QueryRow(ctx, `select exists(select 1 from public.import_batches where file_sha256 = $1 and status = 'published' and id <> $2)`, hash, batchID).Scan(&duplicate); err != nil {
		return err
	}
	if duplicate {
		return ErrAlreadyPublished
	}
	if currentBatchID != nil {
		var appealExists bool
		if err := tx.QueryRow(ctx, `select exists(select 1 from public.appeals where period_id = $1 and status <> 'cancelled')`, periodID).Scan(&appealExists); err != nil {
			return err
		}
		if appealExists {
			return ErrAppealsExist
		}
		if _, err := tx.Exec(ctx, `update public.import_batches set status = 'superseded' where id = $1`, *currentBatchID); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(ctx, `
		update public.import_batches
		set status = 'published', published_by = $2, published_at = now()
		where id = $1`, batchID, actor.ID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `update public.reporting_periods set published_batch_id = $2 where id = $1`, periodID, batchID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		insert into public.notifications (user_id, kind, title, body, action_url)
		select id, 'period_published', 'Data periode ' || $1 || ' tersedia',
		       'Data presensi, potongan, dan ringkasan hari kerja Anda telah diperbarui.', '/#dashboard'
		from public.users where is_active and deleted_at is null`, label); err != nil {
		return err
	}
	if err := queuePublishedPeriodJobs(ctx, tx, label, batchID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		insert into public.audit_logs (actor_id, action, entity_type, entity_id, metadata, ip_address)
		values ($1, 'import.publish', 'import_batch', $2, jsonb_build_object('period', $3), nullif($4, '')::inet)`, actor.ID, batchID, label, ip); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Service) RejectDraft(ctx context.Context, actor auth.Principal, batchID, ip string) error {
	command, err := s.pool.Exec(ctx, `
		update public.import_batches set status = 'rejected'
		where id = $1 and status = 'draft'`, batchID)
	if err != nil {
		return err
	}
	if command.RowsAffected() == 0 {
		return errors.New("draft tidak ditemukan atau sudah diproses")
	}
	_, _ = s.pool.Exec(ctx, `
		insert into public.audit_logs (actor_id, action, entity_type, entity_id, ip_address)
		values ($1, 'import.reject', 'import_batch', $2, nullif($3, '')::inet)`, actor.ID, batchID, ip)
	return nil
}

func (s *Service) stage(ctx context.Context, actor auth.Principal, filename string, parsed ParseResult, users map[string]userMatch, preview Preview) (string, string, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", "", err
	}
	defer tx.Rollback(ctx)
	lockKey := parsed.PeriodStart.Format("2006-01-02") + ":" + parsed.PeriodEnd.Format("2006-01-02")
	if _, err := tx.Exec(ctx, `select pg_advisory_xact_lock(hashtext($1))`, lockKey); err != nil {
		return "", "", err
	}
	var periodID string
	err = tx.QueryRow(ctx, `
		insert into public.reporting_periods (label, period_start, period_end)
		values ($1, $2, $3)
		on conflict (period_start, period_end) do update set label = excluded.label
		returning id::text`, parsed.PeriodLabel, parsed.PeriodStart, parsed.PeriodEnd).Scan(&periodID)
	if err != nil {
		return "", "", err
	}
	var version int
	if err := tx.QueryRow(ctx, `select coalesce(max(version), 0) + 1 from public.import_batches where period_id = $1`, periodID).Scan(&version); err != nil {
		return "", "", err
	}
	warningJSON, _ := json.Marshal(map[string]any{
		"warnings":           preview.Warnings,
		"unit_mismatches":    preview.UnitMismatches,
		"blank_rows_ignored": preview.BlankRowsIgnored,
	})
	var batchID string
	err = tx.QueryRow(ctx, `
		insert into public.import_batches (
		  period_id, version, original_filename, file_sha256, file_size_bytes, sheet_name,
		  integrity_status, status, row_count, employee_count, deduction_day_count,
		  total_deduction_rate, warning_summary, created_by
		) values ($1, $2, $3, $4, $5, $6, $7, 'draft', $8, $9, $10, $11, $12::jsonb, $13)
		returning id::text`, periodID, version, filename, parsed.FileSHA256, parsed.FileSize,
		parsed.SheetName, parsed.IntegrityStatus, len(parsed.Records), preview.Employees,
		preview.DeductionDays, preview.TotalDeduction, string(warningJSON), actor.ID).Scan(&batchID)
	if err != nil {
		return "", "", err
	}

	rows := make([][]any, 0, len(parsed.Records))
	for _, record := range parsed.Records {
		components, _ := json.Marshal(record.DeductionComponents)
		var checkIn, checkOut any
		if record.CheckIn != nil {
			checkIn = pgTime(*record.CheckIn)
		}
		if record.CheckOut != nil {
			checkOut = pgTime(*record.CheckOut)
		}
		rows = append(rows, []any{
			batchID, periodID, users[record.NIP].ID, record.SourceRow, record.WorkDate,
			checkIn, checkOut, nullString(record.LateCode), nullString(record.EarlyLeaveCode),
			nullString(record.ShiftCode), nullString(record.AttendanceStatus), nullString(record.LeaveType),
			nullString(record.AssignmentType), nullString(record.SourceConfirmation), nullString(record.Notes),
			nullString(record.SourceDivision), nullString(record.SourcePlacement), record.DeductionRate, json.RawMessage(components),
		})
	}
	columns := []string{
		"batch_id", "period_id", "user_id", "source_row", "work_date", "check_in", "check_out",
		"late_code", "early_leave_code", "shift_code", "attendance_status", "leave_type",
		"assignment_type", "source_confirmation", "notes", "source_division", "source_placement",
		"deduction_rate", "deduction_components",
	}
	count, err := tx.CopyFrom(ctx, pgx.Identifier{"public", "attendance_records"}, columns, pgx.CopyFromRows(rows))
	if err != nil {
		return "", "", fmt.Errorf("simpan staging import: %w", err)
	}
	if int(count) != len(parsed.Records) {
		return "", "", fmt.Errorf("jumlah baris tersimpan tidak sesuai: %d dari %d", count, len(parsed.Records))
	}
	if _, err := tx.Exec(ctx, `
		insert into public.audit_logs (actor_id, action, entity_type, entity_id, metadata)
		values ($1, 'import.stage', 'import_batch', $2, jsonb_build_object('rows', $3, 'period', $4))`,
		actor.ID, batchID, len(parsed.Records), parsed.PeriodLabel); err != nil {
		return "", "", err
	}
	if err := tx.Commit(ctx); err != nil {
		return "", "", err
	}
	return batchID, periodID, nil
}

func (s *Service) loadUsers(ctx context.Context) (map[string]userMatch, error) {
	rows, err := s.pool.Query(ctx, `
		select u.nip, u.id::text, u.name, coalesce(un.source_name, un.name)
		from public.users u join public.units un on un.id = u.unit_id
		where u.is_active and u.deleted_at is null`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := map[string]userMatch{}
	for rows.Next() {
		var nip string
		var matched userMatch
		if err := rows.Scan(&nip, &matched.ID, &matched.Name, &matched.UnitSource); err != nil {
			return nil, err
		}
		result[nip] = matched
	}
	return result, rows.Err()
}

func (s *Service) loadRules(ctx context.Context) (map[string]Rule, error) {
	rows, err := s.pool.Query(ctx, `
		select source_field, code, label, rate
		from public.deduction_rules where is_active order by sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := map[string]Rule{}
	for rows.Next() {
		var rule Rule
		if err := rows.Scan(&rule.SourceField, &rule.Code, &rule.Label, &rule.Rate); err != nil {
			return nil, err
		}
		result[ruleKey(rule.SourceField, rule.Code)] = rule
	}
	return result, rows.Err()
}

func applyRules(record *RawRecord, rules map[string]Rule) {
	values := map[string]string{
		"late":        record.LateCode,
		"early_leave": record.EarlyLeaveCode,
		"leave":       record.LeaveType,
		"status":      record.AttendanceStatus,
		"shift":       record.ShiftCode,
	}
	for field, value := range values {
		if value == "" {
			continue
		}
		if rule, ok := rules[ruleKey(field, value)]; ok {
			record.DeductionRate += rule.Rate
			record.DeductionComponents = append(record.DeductionComponents, DeductionComponent{
				SourceField: field, Code: value, Label: rule.Label, Rate: rule.Rate,
			})
		}
	}
}

func ruleKey(field, value string) string {
	field = strings.ToLower(strings.TrimSpace(field))
	value = strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(value)), " "))
	if field == "status" && value == "izin tidak masuk" {
		value = "i"
	}
	return field + "\x00" + value
}

func normalizeUnit(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, "rumah tangga", "rumah tangga")
	return strings.Join(strings.Fields(value), " ")
}

func appendUnique(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	if len(values) < 50 {
		return append(values, value)
	}
	return values
}

func nullString(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}

func pgTime(value time.Time) pgtype.Time {
	microseconds := int64(value.Hour()*3600+value.Minute()*60+value.Second()) * 1_000_000
	return pgtype.Time{Microseconds: microseconds, Valid: true}
}

func FormatPeriod(start, end time.Time) string {
	return start.Format("02-01-2006") + " s.d. " + end.Format("02-01-2006")
}
