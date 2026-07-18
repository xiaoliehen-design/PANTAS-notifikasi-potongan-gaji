package httpapi

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/bcpriok/pantas/internal/auth"
	"github.com/jackc/pgx/v5"
)

func (a *App) appeals(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	rows, err := a.pool.Query(request.Context(), `
		select ap.id::text, rp.id::text, rp.label, ap.status, ap.submitted_at, ap.finalized_at,
		       ai.id::text, ar.work_date, ai.original_deduction_rate, rc.label,
		       ai.explanation, ai.supervisor_status, coalesce(ai.supervisor_comment, ''),
		       ai.admin_status, coalesce(ai.admin_comment, ''),
		       coalesce(ai.adjusted_deduction_rate, ar.deduction_rate),
		       (select count(*) from public.appeal_documents ad where ad.appeal_item_id = ai.id)
		from public.appeals ap
		join public.reporting_periods rp on rp.id = ap.period_id
		join public.appeal_items ai on ai.appeal_id = ap.id
		join public.attendance_records ar on ar.id = ai.attendance_record_id
		join public.appeal_reason_categories rc on rc.id = ai.reason_category_id
		where ap.user_id = $1
		order by rp.period_end desc, ar.work_date`, principal.ID)
	if err != nil {
		a.internalError(response, "appeals", err)
		return
	}
	defer rows.Close()
	byAppeal := map[string]map[string]any{}
	order := []string{}
	for rows.Next() {
		var appealID, periodID, periodLabel, status, itemID, reason, explanation, supervisorStatus, supervisorComment, adminStatus, adminComment string
		var submitted time.Time
		var finalized *time.Time
		var workDate time.Time
		var original, adjusted float64
		var documents int
		if err := rows.Scan(&appealID, &periodID, &periodLabel, &status, &submitted, &finalized, &itemID, &workDate, &original, &reason, &explanation, &supervisorStatus, &supervisorComment, &adminStatus, &adminComment, &adjusted, &documents); err != nil {
			a.internalError(response, "appeals scan", err)
			return
		}
		appeal, ok := byAppeal[appealID]
		if !ok {
			appeal = map[string]any{"id": appealID, "period_id": periodID, "period_label": periodLabel, "status": status, "submitted_at": submitted, "finalized_at": finalized, "items": []map[string]any{}}
			byAppeal[appealID] = appeal
			order = append(order, appealID)
		}
		items := appeal["items"].([]map[string]any)
		items = append(items, map[string]any{"id": itemID, "date": workDate.Format("2006-01-02"), "original": original, "reason": reason, "explanation": explanation, "supervisor_status": supervisorStatus, "supervisor_comment": supervisorComment, "admin_status": adminStatus, "admin_comment": adminComment, "adjusted": adjusted, "document_count": documents})
		appeal["items"] = items
	}
	result := make([]map[string]any, 0, len(order))
	for _, id := range order {
		result = append(result, byAppeal[id])
	}
	writeJSON(response, http.StatusOK, map[string]any{"items": result})
}

func (a *App) appealOptions(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	periodID := request.URL.Query().Get("period_id")
	if periodID == "" {
		period, err := a.currentPeriod(request)
		if errors.Is(err, pgx.ErrNoRows) {
			writeJSON(response, http.StatusOK, map[string]any{"period": nil, "days": []any{}, "reasons": []any{}})
			return
		}
		if err != nil {
			a.internalError(response, "appeal options period", err)
			return
		}
		periodID = period.ID
	}
	if !validUUID(periodID) {
		writeError(response, http.StatusBadRequest, "Periode tidak valid.", "invalid_period")
		return
	}
	rows, err := a.pool.Query(request.Context(), `
		select ea.id, ea.work_date, ea.deduction_rate, ea.deduction_components,
		       coalesce(ea.late_code, ''), coalesce(ea.early_leave_code, ''), coalesce(ea.notes, '')
		from public.effective_attendance ea
		left join public.appeal_items ai on ai.attendance_record_id = ea.id
		where ea.user_id = $1 and ea.period_id = $2 and ea.deduction_rate > 0 and ai.id is null
		order by ea.work_date`, principal.ID, periodID)
	if err != nil {
		a.internalError(response, "appeal options", err)
		return
	}
	days := []map[string]any{}
	for rows.Next() {
		var id int64
		var date time.Time
		var rate float64
		var components []byte
		var late, early, notes string
		if err := rows.Scan(&id, &date, &rate, &components, &late, &early, &notes); err != nil {
			rows.Close()
			a.internalError(response, "appeal option scan", err)
			return
		}
		days = append(days, map[string]any{"attendance_id": id, "date": date.Format("2006-01-02"), "rate": rate, "components": jsonRaw(components), "late": late, "early_leave": early, "notes": notes})
	}
	rows.Close()
	reasonRows, err := a.pool.Query(request.Context(), `select id::text, code, label, coalesce(description, '') from public.appeal_reason_categories where is_active order by sort_order, label`)
	if err != nil {
		a.internalError(response, "appeal reasons", err)
		return
	}
	defer reasonRows.Close()
	reasons := []map[string]any{}
	for reasonRows.Next() {
		var id, code, label, description string
		if err := reasonRows.Scan(&id, &code, &label, &description); err != nil {
			a.internalError(response, "appeal reason scan", err)
			return
		}
		reasons = append(reasons, map[string]any{"id": id, "code": code, "label": label, "description": description})
	}
	writeJSON(response, http.StatusOK, map[string]any{"period_id": periodID, "days": days, "reasons": reasons})
}

func (a *App) createAppeal(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	var input struct {
		PeriodID string `json:"period_id"`
		Items    []struct {
			AttendanceID int64  `json:"attendance_id"`
			ReasonID     string `json:"reason_id"`
			Explanation  string `json:"explanation"`
		} `json:"items"`
	}
	if !decodeJSON(response, request, &input) {
		return
	}
	if !validUUID(input.PeriodID) || len(input.Items) == 0 || len(input.Items) > 62 {
		writeError(response, http.StatusUnprocessableEntity, "Pilih sedikitnya satu dan maksimal 62 hari potongan.", "invalid_appeal")
		return
	}
	seen := map[int64]struct{}{}
	for _, item := range input.Items {
		item.Explanation = strings.TrimSpace(item.Explanation)
		if item.AttendanceID <= 0 || !validUUID(item.ReasonID) || len([]rune(item.Explanation)) < 10 || len([]rune(item.Explanation)) > 3000 {
			writeError(response, http.StatusUnprocessableEntity, "Setiap hari memerlukan kategori dan penjelasan 10–3.000 karakter.", "invalid_appeal_item")
			return
		}
		if _, duplicate := seen[item.AttendanceID]; duplicate {
			writeError(response, http.StatusUnprocessableEntity, "Hari potongan tidak boleh diduplikasi.", "duplicate_appeal_item")
			return
		}
		seen[item.AttendanceID] = struct{}{}
	}

	tx, err := a.pool.BeginTx(request.Context(), pgx.TxOptions{})
	if err != nil {
		a.internalError(response, "appeal begin", err)
		return
	}
	defer tx.Rollback(request.Context())
	var periodLabel string
	if err := tx.QueryRow(request.Context(), `select label from public.reporting_periods where id = $1 and published_batch_id is not null`, input.PeriodID).Scan(&periodLabel); err != nil {
		writeError(response, http.StatusUnprocessableEntity, "Periode tidak tersedia.", "invalid_period")
		return
	}
	var appealID string
	err = tx.QueryRow(request.Context(), `
		insert into public.appeals (user_id, period_id, status)
		values ($1, $2, 'supervisor_review')
		on conflict (user_id, period_id) do nothing
		returning id::text`, principal.ID, input.PeriodID).Scan(&appealID)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(response, http.StatusConflict, "Banding untuk periode ini sudah pernah diajukan.", "appeal_exists")
		return
	}
	if err != nil {
		a.internalError(response, "appeal insert", err)
		return
	}
	supervisors, err := findSupervisors(request.Context(), tx, principal)
	if err != nil {
		a.internalError(response, "appeal supervisor", err)
		return
	}
	directAdmin := len(supervisors) == 0
	createdItems := []map[string]any{}
	for _, item := range input.Items {
		var rate float64
		var workDate time.Time
		err := tx.QueryRow(request.Context(), `
			select deduction_rate, work_date
			from public.effective_attendance
			where id = $1 and user_id = $2 and period_id = $3 and deduction_rate > 0
			  and appeal_item_id is null`, item.AttendanceID, principal.ID, input.PeriodID).Scan(&rate, &workDate)
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(response, http.StatusUnprocessableEntity, "Salah satu hari tidak lagi dapat dibanding.", "appeal_day_unavailable")
			return
		}
		if err != nil {
			a.internalError(response, "appeal attendance", err)
			return
		}
		var reasonExists bool
		if err := tx.QueryRow(request.Context(), `select exists(select 1 from public.appeal_reason_categories where id = $1 and is_active)`, item.ReasonID).Scan(&reasonExists); err != nil || !reasonExists {
			writeError(response, http.StatusUnprocessableEntity, "Kategori alasan tidak tersedia.", "invalid_reason")
			return
		}
		var itemID string
		supervisorStatus := "pending"
		var supervisorComment any
		if directAdmin {
			supervisorStatus = "accepted"
			supervisorComment = "Diteruskan langsung ke administrator karena tidak ada atasan pada tingkat berikutnya."
		}
		err = tx.QueryRow(request.Context(), `
			insert into public.appeal_items (
			  appeal_id, attendance_record_id, reason_category_id, explanation, original_deduction_rate,
			  supervisor_status, supervisor_comment, supervisor_reviewed_at
			) values ($1, $2, $3, $4, $5, $6, $7,
			  case when $6 = 'accepted' then now() else null end)
			returning id::text`, appealID, item.AttendanceID, item.ReasonID, strings.TrimSpace(item.Explanation), rate, supervisorStatus, supervisorComment).Scan(&itemID)
		if err != nil {
			a.internalError(response, "appeal item insert", err)
			return
		}
		createdItems = append(createdItems, map[string]any{"id": itemID, "attendance_id": item.AttendanceID, "date": workDate.Format("2006-01-02")})
	}
	if directAdmin {
		if _, err := tx.Exec(request.Context(), `update public.appeals set status = 'admin_review' where id = $1`, appealID); err != nil {
			a.internalError(response, "appeal status", err)
			return
		}
	}
	if _, err := tx.Exec(request.Context(), `
		insert into public.audit_logs (actor_id, action, entity_type, entity_id, metadata, ip_address)
		values ($1, 'appeal.submit', 'appeal', $2, jsonb_build_object('period', $3, 'items', $4), nullif($5, '')::inet)`,
		principal.ID, appealID, periodLabel, len(input.Items), auth.ClientIP(request, a.cfg.TrustProxy)); err != nil {
		a.internalError(response, "appeal audit", err)
		return
	}
	for _, supervisor := range supervisors {
		if _, err := tx.Exec(request.Context(), `
			insert into public.notifications (user_id, kind, title, body, action_url)
			values ($1, 'appeal_submitted', 'Banding menunggu verifikasi', $2, '/#reviews')`, supervisor.ID, principal.Name+" mengajukan banding periode "+periodLabel+"."); err != nil {
			a.internalError(response, "appeal notify", err)
			return
		}
		if supervisor.Email != "" && supervisor.EmailVerified {
			payload := map[string]any{"name": supervisor.Name, "period": periodLabel}
			if err := queueTx(request.Context(), tx, supervisor.ID, supervisor.Email, "appeal_submitted", payload); err != nil {
				a.internalError(response, "appeal email", err)
				return
			}
		}
	}
	if err := tx.Commit(request.Context()); err != nil {
		a.internalError(response, "appeal commit", err)
		return
	}
	writeJSON(response, http.StatusCreated, map[string]any{"id": appealID, "status": map[bool]string{true: "admin_review", false: "supervisor_review"}[directAdmin], "items": createdItems})
}

type supervisorContact struct {
	ID, Name, Email string
	EmailVerified   bool
}

func findSupervisors(ctx context.Context, tx pgx.Tx, principal auth.Principal) ([]supervisorContact, error) {
	query := ""
	args := []any{}
	switch principal.PositionRole {
	case "staff":
		query = `select id::text, name, coalesce(email, ''), email_verified_at is not null from public.users where unit_id = $1 and position_role = 'section_head' and is_active and deleted_at is null`
		args = append(args, principal.UnitID)
	case "section_head":
		query = `select u.id::text, u.name, coalesce(u.email, ''), u.email_verified_at is not null from public.users u join public.units child on child.parent_id = u.unit_id where child.id = $1 and u.position_role = 'division_head' and u.is_active and u.deleted_at is null`
		args = append(args, principal.UnitID)
	case "division_head", "functional":
		query = `select id::text, name, coalesce(email, ''), email_verified_at is not null from public.users where position_role = 'office_head' and is_active and deleted_at is null`
	case "office_head":
		return nil, nil
	default:
		return nil, nil
	}
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []supervisorContact{}
	for rows.Next() {
		var item supervisorContact
		if err := rows.Scan(&item.ID, &item.Name, &item.Email, &item.EmailVerified); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (a *App) uploadAppealDocument(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	itemID := request.PathValue("id")
	if !validUUID(itemID) {
		writeError(response, http.StatusBadRequest, "Item banding tidak valid.", "invalid_id")
		return
	}
	var appealID string
	var owner bool
	err := a.pool.QueryRow(request.Context(), `
		select ap.id::text, ap.user_id = $2 and ai.admin_status = 'pending'
		from public.appeal_items ai join public.appeals ap on ap.id = ai.appeal_id
		where ai.id = $1`, itemID, principal.ID).Scan(&appealID, &owner)
	if err != nil || !owner {
		writeError(response, http.StatusForbidden, "Dokumen tidak dapat ditambahkan setelah keputusan final.", "forbidden")
		return
	}
	request.Body = http.MaxBytesReader(response, request.Body, a.cfg.MaxDocumentBytes)
	data, err := io.ReadAll(request.Body)
	if err != nil || len(data) == 0 {
		writeError(response, http.StatusRequestEntityTooLarge, "Dokumen kosong atau melebihi 5 MB.", "invalid_document")
		return
	}
	detected := http.DetectContentType(data)
	if semicolon := strings.IndexByte(detected, ';'); semicolon >= 0 {
		detected = detected[:semicolon]
	}
	extension := map[string]string{"application/pdf": ".pdf", "image/jpeg": ".jpg", "image/png": ".png"}[detected]
	if extension == "" {
		writeError(response, http.StatusUnprocessableEntity, "Dokumen harus PDF, JPG, atau PNG.", "invalid_document_type")
		return
	}
	original := sanitizeFilename(request.Header.Get("X-Filename"), extension)
	randomPart := make([]byte, 16)
	if _, err := rand.Read(randomPart); err != nil {
		a.internalError(response, "document random", err)
		return
	}
	storagePath := fmt.Sprintf("appeals/%s/%s/%s%s", appealID, itemID, hex.EncodeToString(randomPart), extension)
	if err := a.storage.Upload(request.Context(), storagePath, detected, data); err != nil {
		a.internalError(response, "document upload", err)
		return
	}
	hash := sha256.Sum256(data)
	var documentID string
	err = a.pool.QueryRow(request.Context(), `
		insert into public.appeal_documents (appeal_item_id, storage_path, original_filename, mime_type, size_bytes, sha256, uploaded_by)
		values ($1, $2, $3, $4, $5, $6, $7) returning id::text`,
		itemID, storagePath, original, detected, len(data), hex.EncodeToString(hash[:]), principal.ID).Scan(&documentID)
	if err != nil {
		_ = a.storage.Delete(request.Context(), []string{storagePath})
		a.internalError(response, "document metadata", err)
		return
	}
	writeJSON(response, http.StatusCreated, map[string]any{"id": documentID, "filename": original, "mime_type": detected, "size": len(data)})
}

func (a *App) supervisorQueue(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	if principal.IsAdmin {
		writeError(response, http.StatusForbidden, "Administrator hanya memberikan keputusan final banding.", "forbidden")
		return
	}
	items, err := a.reviewQueue(request, principal, false)
	if err != nil {
		a.internalError(response, "supervisor queue", err)
		return
	}
	writeJSON(response, http.StatusOK, map[string]any{"items": items})
}

func (a *App) adminReviewQueue(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	items, err := a.reviewQueue(request, principal, true)
	if err != nil {
		a.internalError(response, "admin queue", err)
		return
	}
	writeJSON(response, http.StatusOK, map[string]any{"items": items})
}

func (a *App) reviewQueue(request *http.Request, principal auth.Principal, adminQueue bool) ([]map[string]any, error) {
	filter := reviewScopeSQL(principal)
	statusFilter := "ai.supervisor_status = 'pending'"
	args := []any{principal.ID, principal.UnitID}
	if adminQueue {
		filter = "true"
		statusFilter = "ai.supervisor_status <> 'pending' and ai.admin_status = 'pending'"
		args = nil
	}
	query := fmt.Sprintf(`
		select ai.id::text, ap.id::text, u.id::text, u.nip, u.name, un.name, rp.label,
		       ar.work_date, ai.original_deduction_rate, rc.label, ai.explanation,
		       ai.supervisor_status, coalesce(ai.supervisor_comment, ''), ai.admin_status,
		       (select count(*) from public.appeal_documents ad where ad.appeal_item_id = ai.id)
		from public.appeal_items ai
		join public.appeals ap on ap.id = ai.appeal_id
		join public.users u on u.id = ap.user_id
		join public.units un on un.id = u.unit_id
		join public.reporting_periods rp on rp.id = ap.period_id
		join public.attendance_records ar on ar.id = ai.attendance_record_id
		join public.appeal_reason_categories rc on rc.id = ai.reason_category_id
		where %s and %s
		order by ap.submitted_at, ar.work_date`, statusFilter, filter)
	rows, err := a.pool.Query(request.Context(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		var itemID, appealID, userID, nip, name, unit, period, reason, explanation, supervisorStatus, supervisorComment, adminStatus string
		var date time.Time
		var rate float64
		var documents int
		if err := rows.Scan(&itemID, &appealID, &userID, &nip, &name, &unit, &period, &date, &rate, &reason, &explanation, &supervisorStatus, &supervisorComment, &adminStatus, &documents); err != nil {
			return nil, err
		}
		items = append(items, map[string]any{"id": itemID, "appeal_id": appealID, "user_id": userID, "nip": nip, "name": name, "unit": unit, "period": period, "date": date.Format("2006-01-02"), "rate": rate, "reason": reason, "explanation": explanation, "supervisor_status": supervisorStatus, "supervisor_comment": supervisorComment, "admin_status": adminStatus, "document_count": documents})
	}
	return items, rows.Err()
}

func (a *App) supervisorDecision(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	if principal.IsAdmin {
		writeError(response, http.StatusForbidden, "Administrator hanya memberikan keputusan final banding.", "forbidden")
		return
	}
	itemID := request.PathValue("id")
	var input struct {
		Decision string `json:"decision"`
		Comment  string `json:"comment"`
	}
	if !validUUID(itemID) || !decodeJSON(response, request, &input) {
		return
	}
	input.Comment = strings.TrimSpace(input.Comment)
	if (input.Decision != "accepted" && input.Decision != "rejected") || len([]rune(input.Comment)) > 2000 {
		writeError(response, http.StatusUnprocessableEntity, "Keputusan atau komentar tidak valid.", "invalid_review")
		return
	}
	allowed, err := a.canReviewItem(request.Context(), principal, itemID)
	if err != nil {
		a.internalError(response, "supervisor permission", err)
		return
	}
	if !allowed {
		writeError(response, http.StatusForbidden, "Item banding bukan dalam lingkup verifikasi Anda.", "forbidden")
		return
	}
	tx, err := a.pool.BeginTx(request.Context(), pgx.TxOptions{})
	if err != nil {
		a.internalError(response, "supervisor review begin", err)
		return
	}
	defer tx.Rollback(request.Context())
	var appealID string
	err = tx.QueryRow(request.Context(), `
		update public.appeal_items
		set supervisor_status = $2, supervisor_by = $3, supervisor_comment = nullif($4, ''), supervisor_reviewed_at = now()
		where id = $1 and supervisor_status = 'pending'
		returning appeal_id::text`, itemID, input.Decision, principal.ID, input.Comment).Scan(&appealID)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(response, http.StatusConflict, "Item banding sudah diverifikasi.", "already_reviewed")
		return
	}
	if err != nil {
		a.internalError(response, "supervisor review", err)
		return
	}
	var pending int
	if err := tx.QueryRow(request.Context(), `select count(*) from public.appeal_items where appeal_id = $1 and supervisor_status = 'pending'`, appealID).Scan(&pending); err != nil {
		a.internalError(response, "supervisor pending", err)
		return
	}
	if pending == 0 {
		if _, err := tx.Exec(request.Context(), `update public.appeals set status = 'admin_review' where id = $1`, appealID); err != nil {
			a.internalError(response, "supervisor appeal status", err)
			return
		}
		if _, err := tx.Exec(request.Context(), `
			insert into public.notifications (user_id, kind, title, body, action_url)
			select account_id, 'appeal_admin_queue', 'Banding menunggu keputusan admin', 'Verifikasi atasan telah selesai.', '/#admin-reviews'
			from public.admin_accounts where is_active`); err != nil {
			a.internalError(response, "admin notify", err)
			return
		}
	}
	if err := tx.Commit(request.Context()); err != nil {
		a.internalError(response, "supervisor commit", err)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

func (a *App) adminDecision(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	itemID := request.PathValue("id")
	var input struct {
		Decision     string   `json:"decision"`
		Comment      string   `json:"comment"`
		AdjustedRate *float64 `json:"adjusted_rate"`
	}
	if !validUUID(itemID) || !decodeJSON(response, request, &input) {
		return
	}
	input.Comment = strings.TrimSpace(input.Comment)
	if (input.Decision != "approved" && input.Decision != "rejected") || len([]rune(input.Comment)) > 2000 {
		writeError(response, http.StatusUnprocessableEntity, "Keputusan atau komentar tidak valid.", "invalid_review")
		return
	}
	tx, err := a.pool.BeginTx(request.Context(), pgx.TxOptions{})
	if err != nil {
		a.internalError(response, "admin review begin", err)
		return
	}
	defer tx.Rollback(request.Context())
	var appealID, userID, periodLabel, name, email string
	var original float64
	var emailVerified bool
	err = tx.QueryRow(request.Context(), `
		select ai.appeal_id::text, ap.user_id::text, rp.label, u.name, coalesce(u.email, ''),
		       u.email_verified_at is not null, ai.original_deduction_rate
		from public.appeal_items ai
		join public.appeals ap on ap.id = ai.appeal_id
		join public.reporting_periods rp on rp.id = ap.period_id
		join public.users u on u.id = ap.user_id
		where ai.id = $1 and ai.supervisor_status <> 'pending' and ai.admin_status = 'pending'
		for update of ai`, itemID).Scan(&appealID, &userID, &periodLabel, &name, &email, &emailVerified, &original)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(response, http.StatusConflict, "Item belum siap atau sudah diputuskan.", "already_reviewed")
		return
	}
	if err != nil {
		a.internalError(response, "admin item", err)
		return
	}
	adjusted := original
	if input.Decision == "approved" {
		adjusted = 0
		if input.AdjustedRate != nil {
			adjusted = *input.AdjustedRate
		}
	}
	if adjusted < 0 || adjusted > original {
		writeError(response, http.StatusUnprocessableEntity, "Potongan hasil tidak boleh negatif atau melebihi potongan awal.", "invalid_adjustment")
		return
	}
	if _, err := tx.Exec(request.Context(), `
		update public.appeal_items
		set admin_status = $2, admin_by = $3, admin_comment = nullif($4, ''),
		    adjusted_deduction_rate = $5, admin_reviewed_at = now()
		where id = $1`, itemID, input.Decision, principal.ID, input.Comment, adjusted); err != nil {
		a.internalError(response, "admin review", err)
		return
	}
	var pending int
	if err := tx.QueryRow(request.Context(), `select count(*) from public.appeal_items where appeal_id = $1 and admin_status = 'pending'`, appealID).Scan(&pending); err != nil {
		a.internalError(response, "admin pending", err)
		return
	}
	if pending == 0 {
		if _, err := tx.Exec(request.Context(), `update public.appeals set status = 'finalized', finalized_at = now() where id = $1`, appealID); err != nil {
			a.internalError(response, "appeal finalized", err)
			return
		}
	}
	if _, err := tx.Exec(request.Context(), `
		insert into public.notifications (user_id, kind, title, body, action_url)
		values ($1, 'appeal_reviewed', 'Status banding diperbarui', $2, '/#appeals')`, userID, "Keputusan untuk salah satu hari pada periode "+periodLabel+" telah tersedia."); err != nil {
		a.internalError(response, "appeal user notify", err)
		return
	}
	if email != "" && emailVerified {
		if err := queueTx(request.Context(), tx, userID, email, "appeal_reviewed", map[string]any{"name": name, "period": periodLabel}); err != nil {
			a.internalError(response, "appeal user email", err)
			return
		}
	}
	if err := tx.Commit(request.Context()); err != nil {
		a.internalError(response, "admin review commit", err)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

func (a *App) downloadDocument(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	documentID := request.PathValue("id")
	if !validUUID(documentID) {
		writeError(response, http.StatusBadRequest, "Dokumen tidak valid.", "invalid_id")
		return
	}
	var itemID, storagePath, filename, mimeType string
	err := a.pool.QueryRow(request.Context(), `
		select ad.appeal_item_id::text, ad.storage_path, ad.original_filename, ad.mime_type
		from public.appeal_documents ad where ad.id = $1`, documentID).Scan(&itemID, &storagePath, &filename, &mimeType)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(response, http.StatusNotFound, "Dokumen tidak ditemukan.", "not_found")
		return
	}
	if err != nil {
		a.internalError(response, "document lookup", err)
		return
	}
	allowed, err := a.canAccessItem(request.Context(), principal, itemID)
	if err != nil {
		a.internalError(response, "document permission", err)
		return
	}
	if !allowed {
		writeError(response, http.StatusForbidden, "Anda tidak berhak membuka dokumen ini.", "forbidden")
		return
	}
	body, returnedType, err := a.storage.Download(request.Context(), storagePath)
	if err != nil {
		a.internalError(response, "document download", err)
		return
	}
	defer body.Close()
	if returnedType != "" {
		mimeType = strings.Split(returnedType, ";")[0]
	}
	response.Header().Set("Content-Type", mimeType)
	response.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": filename}))
	response.Header().Set("Cache-Control", "private, no-store")
	_, _ = io.Copy(response, body)
}

func (a *App) appealDocuments(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	itemID := request.PathValue("id")
	if !validUUID(itemID) {
		writeError(response, http.StatusBadRequest, "Item banding tidak valid.", "invalid_id")
		return
	}
	allowed, err := a.canAccessItem(request.Context(), principal, itemID)
	if err != nil {
		a.internalError(response, "document list permission", err)
		return
	}
	if !allowed {
		writeError(response, http.StatusForbidden, "Anda tidak berhak melihat dokumen ini.", "forbidden")
		return
	}
	rows, err := a.pool.Query(request.Context(), `
		select id::text, original_filename, mime_type, size_bytes, created_at
		from public.appeal_documents where appeal_item_id = $1 order by created_at`, itemID)
	if err != nil {
		a.internalError(response, "document list", err)
		return
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		var id, filename, mimeType string
		var size int64
		var created time.Time
		if err := rows.Scan(&id, &filename, &mimeType, &size, &created); err != nil {
			a.internalError(response, "document list scan", err)
			return
		}
		items = append(items, map[string]any{"id": id, "filename": filename, "mime_type": mimeType, "size": size, "created_at": created})
	}
	writeJSON(response, http.StatusOK, map[string]any{"items": items})
}

func (a *App) canReviewItem(ctx context.Context, principal auth.Principal, itemID string) (bool, error) {
	if principal.IsAdmin {
		return true, nil
	}
	query := fmt.Sprintf(`
		select exists(
		  select 1 from public.appeal_items ai
		  join public.appeals ap on ap.id = ai.appeal_id
		  join public.users u on u.id = ap.user_id
		  join public.units un on un.id = u.unit_id
		  where ai.id = $3 and %s
		)`, reviewScopeSQL(principal))
	var allowed bool
	err := a.pool.QueryRow(ctx, query, principal.ID, principal.UnitID, itemID).Scan(&allowed)
	return allowed, err
}

func (a *App) canAccessItem(ctx context.Context, principal auth.Principal, itemID string) (bool, error) {
	if principal.IsAdmin {
		return true, nil
	}
	var owner bool
	if err := a.pool.QueryRow(ctx, `
		select exists(select 1 from public.appeal_items ai join public.appeals ap on ap.id = ai.appeal_id where ai.id = $1 and ap.user_id = $2)`, itemID, principal.ID).Scan(&owner); err != nil {
		return false, err
	}
	if owner {
		return true, nil
	}
	return a.canReviewItem(ctx, principal, itemID)
}

func reviewScopeSQL(principal auth.Principal) string {
	switch principal.PositionRole {
	case "section_head":
		return "u.unit_id = $2 and u.id <> $1 and u.position_role = 'staff'"
	case "division_head":
		return "un.parent_id = $2 and u.position_role = 'section_head'"
	case "office_head":
		return "u.position_role in ('division_head', 'functional')"
	default:
		if principal.IsAdmin {
			return "true"
		}
		return "false"
	}
}

func queueTx(ctx context.Context, tx pgx.Tx, userID, destination, template string, payload map[string]any) error {
	encoded, err := jsonBytes(payload)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
		insert into public.notification_jobs (user_id, channel, destination, template_code, payload)
		values ($1, 'email', $2, $3, $4::jsonb)`, userID, destination, template, string(encoded))
	return err
}

func sanitizeFilename(value, fallbackExtension string) string {
	decoded, err := url.QueryUnescape(value)
	if err == nil {
		value = decoded
	}
	value = filepath.Base(strings.TrimSpace(value))
	value = strings.Map(func(char rune) rune {
		if char < 32 || strings.ContainsRune(`<>:"/\|?*`, char) {
			return -1
		}
		return char
	}, value)
	if value == "" || len(value) > 180 {
		value = "dokumen" + fallbackExtension
	} else {
		value = strings.TrimSuffix(value, filepath.Ext(value)) + fallbackExtension
	}
	return value
}

func jsonRaw(value []byte) any { return jsonMessage(value) }

// Alias kecil agar file handler tidak perlu mengekspos detail encoding JSON.
type jsonMessage []byte

func (m jsonMessage) MarshalJSON() ([]byte, error) {
	if len(m) == 0 {
		return []byte("[]"), nil
	}
	return m, nil
}

func jsonBytes(value any) ([]byte, error) {
	return json.Marshal(value)
}
