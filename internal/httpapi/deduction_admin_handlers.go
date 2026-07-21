package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bcpriok/pantas/internal/auth"
	"github.com/jackc/pgx/v5"
)

type manualReason struct {
	Code   string
	Detail *string
}

type manualDeductionInput struct {
	PeriodID     string `json:"period_id"`
	UserID       string `json:"user_id"`
	WorkDate     string `json:"work_date"`
	RuleID       string `json:"rule_id"`
	Notes        string `json:"notes"`
	ReasonCode   string `json:"reason_code"`
	ReasonDetail string `json:"reason_detail"`
}

func (a *App) adminDeductions(response http.ResponseWriter, request *http.Request, _ auth.Principal) {
	periods, err := a.publishedPeriods(request)
	if err != nil {
		a.internalError(response, "admin deduction periods", err)
		return
	}
	if len(periods) == 0 {
		writeJSON(response, http.StatusOK, map[string]any{"periods": periods, "current_period": nil, "items": []any{}, "page": 1, "limit": 50, "total": 0})
		return
	}
	periodID := strings.TrimSpace(request.URL.Query().Get("period_id"))
	if periodID == "" {
		periodID = periods[0].ID
	}
	if !validUUID(periodID) {
		writeError(response, http.StatusBadRequest, "Periode tidak valid.", "invalid_period")
		return
	}
	var selected *periodInfo
	for index := range periods {
		if periods[index].ID == periodID {
			selected = &periods[index]
			break
		}
	}
	if selected == nil {
		writeError(response, http.StatusNotFound, "Periode tidak ditemukan atau belum dipublikasikan.", "period_not_found")
		return
	}
	queryText := strings.TrimSpace(request.URL.Query().Get("q"))
	page := parsePositiveInt(request.URL.Query().Get("page"), 1)
	limit := min(parsePositiveInt(request.URL.Query().Get("limit"), 50), 100)
	offset := (page - 1) * limit
	rows, err := a.pool.Query(request.Context(), `
		select ea.id, ea.work_date, u.id::text, u.nip, u.name, un.name,
		       coalesce(ea.late_code, ''), coalesce(ea.early_leave_code, ''),
		       coalesce(ea.shift_code, ''), coalesce(ea.attendance_status, ''),
		       coalesce(ea.leave_type, ''), coalesce(ea.notes, ''),
		       ea.deduction_rate, ea.effective_deduction_rate, ea.deduction_components,
		       ea.record_source, coalesce(ea.last_manual_reason_code, ''),
		       coalesce(ea.last_manual_reason_detail, ''), ea.updated_at,
		       exists(select 1 from public.appeal_items ai where ai.attendance_record_id = ea.id),
		       count(*) over()
		from public.effective_attendance ea
		join public.users u on u.id = ea.user_id
		join public.units un on un.id = u.unit_id
		where ea.period_id = $1 and (ea.deduction_rate > 0 or ea.record_source <> 'import')
		  and ($2 = '' or u.nip ilike '%' || $2 || '%' or u.name ilike '%' || $2 || '%' or un.name ilike '%' || $2 || '%')
		order by u.name, ea.work_date
		limit $3 offset $4`, periodID, queryText, limit, offset)
	if err != nil {
		a.internalError(response, "admin deductions", err)
		return
	}
	defer rows.Close()
	items := []map[string]any{}
	total := 0
	for rows.Next() {
		var id int64
		var workDate, updatedAt time.Time
		var userID, nip, name, unit, late, early, shift, attendanceStatus, leaveType, notes string
		var source, reasonCode, reasonDetail string
		var original, effective float64
		var components []byte
		var hasAppeal bool
		if err := rows.Scan(&id, &workDate, &userID, &nip, &name, &unit, &late, &early, &shift, &attendanceStatus, &leaveType, &notes, &original, &effective, &components, &source, &reasonCode, &reasonDetail, &updatedAt, &hasAppeal, &total); err != nil {
			a.internalError(response, "admin deduction scan", err)
			return
		}
		items = append(items, map[string]any{
			"id": id, "date": workDate.Format("2006-01-02"), "user_id": userID,
			"nip": nip, "name": name, "unit": unit, "late": late, "early_leave": early,
			"shift": shift, "attendance_status": attendanceStatus, "leave_type": leaveType,
			"notes": notes, "original": original, "effective": effective,
			"components": json.RawMessage(components), "record_source": source,
			"reason_code": reasonCode, "reason_detail": reasonDetail,
			"updated_at": updatedAt, "has_appeal": hasAppeal,
		})
	}
	if err := rows.Err(); err != nil {
		a.internalError(response, "admin deduction rows", err)
		return
	}
	writeJSON(response, http.StatusOK, map[string]any{
		"periods": periods, "current_period": selected, "items": items,
		"page": page, "limit": limit, "total": total,
	})
}

func (a *App) adminDeductionOptions(response http.ResponseWriter, request *http.Request, _ auth.Principal) {
	periods, err := a.publishedPeriods(request)
	if err != nil {
		a.internalError(response, "manual deduction periods", err)
		return
	}
	userRows, err := a.pool.Query(request.Context(), `
		select u.id::text, u.nip, u.name, un.name
		from public.users u join public.units un on un.id = u.unit_id
		where u.is_active and u.deleted_at is null
		order by u.name`)
	if err != nil {
		a.internalError(response, "manual deduction users", err)
		return
	}
	users := []map[string]any{}
	for userRows.Next() {
		var id, nip, name, unit string
		if err := userRows.Scan(&id, &nip, &name, &unit); err != nil {
			userRows.Close()
			a.internalError(response, "manual deduction user scan", err)
			return
		}
		users = append(users, map[string]any{"id": id, "nip": nip, "name": name, "unit": unit})
	}
	if err := userRows.Err(); err != nil {
		userRows.Close()
		a.internalError(response, "manual deduction user rows", err)
		return
	}
	userRows.Close()

	ruleRows, err := a.pool.Query(request.Context(), `
		select id::text, source_field, code, label, rate
		from public.deduction_rules where is_active
		order by sort_order, source_field, code`)
	if err != nil {
		a.internalError(response, "manual deduction rules", err)
		return
	}
	defer ruleRows.Close()
	rules := []map[string]any{}
	for ruleRows.Next() {
		var id, source, code, label string
		var rate float64
		if err := ruleRows.Scan(&id, &source, &code, &label, &rate); err != nil {
			a.internalError(response, "manual deduction rule scan", err)
			return
		}
		rules = append(rules, map[string]any{"id": id, "source": source, "code": code, "label": label, "rate": rate})
	}
	if err := ruleRows.Err(); err != nil {
		a.internalError(response, "manual deduction rule rows", err)
		return
	}
	writeJSON(response, http.StatusOK, map[string]any{"periods": periods, "users": users, "rules": rules})
}

func (a *App) adminCreateManualDeduction(response http.ResponseWriter, request *http.Request, actor auth.Principal) {
	var input manualDeductionInput
	if !decodeJSON(response, request, &input) {
		return
	}
	reason, err := validateManualReason(input.ReasonCode, input.ReasonDetail)
	if err != nil || !validUUID(input.PeriodID) || !validUUID(input.UserID) || !validUUID(input.RuleID) || len([]rune(strings.TrimSpace(input.Notes))) > 1000 {
		message := "Periode, pegawai, tanggal, kategori, atau alasan koreksi tidak valid."
		if err != nil {
			message = err.Error()
		}
		writeError(response, http.StatusUnprocessableEntity, message, "invalid_manual_deduction")
		return
	}
	workDate, dateErr := time.Parse("2006-01-02", strings.TrimSpace(input.WorkDate))
	if dateErr != nil {
		writeError(response, http.StatusUnprocessableEntity, "Tanggal potongan tidak valid.", "invalid_work_date")
		return
	}

	tx, err := a.pool.BeginTx(request.Context(), pgx.TxOptions{})
	if err != nil {
		a.internalError(response, "manual deduction begin", err)
		return
	}
	defer tx.Rollback(request.Context())
	var batchID, periodLabel string
	var periodStart, periodEnd time.Time
	err = tx.QueryRow(request.Context(), `
		select published_batch_id::text, label, period_start, period_end
		from public.reporting_periods
		where id = $1 and published_batch_id is not null
		for update`, input.PeriodID).Scan(&batchID, &periodLabel, &periodStart, &periodEnd)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(response, http.StatusNotFound, "Periode belum dipublikasikan.", "period_not_found")
		return
	}
	if err != nil {
		a.internalError(response, "manual deduction period", err)
		return
	}
	if workDate.Before(periodStart) || workDate.After(periodEnd) {
		writeError(response, http.StatusUnprocessableEntity, "Tanggal harus berada di dalam rentang periode.", "date_outside_period")
		return
	}

	var ruleSource, ruleCode, ruleLabel string
	var ruleRate float64
	err = tx.QueryRow(request.Context(), `
		select source_field, code, label, rate
		from public.deduction_rules where id = $1 and is_active`, input.RuleID).Scan(&ruleSource, &ruleCode, &ruleLabel, &ruleRate)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(response, http.StatusUnprocessableEntity, "Kategori potongan tidak tersedia.", "rule_unavailable")
		return
	}
	if err != nil {
		a.internalError(response, "manual deduction rule", err)
		return
	}
	var employeeName, employeeNIP, sourceDivision, sourcePlacement string
	err = tx.QueryRow(request.Context(), `
		select u.name, u.nip, coalesce(parent.name, un.name), coalesce(un.source_name, un.name)
		from public.users u
		join public.units un on un.id = u.unit_id
		left join public.units parent on parent.id = un.parent_id
		where u.id = $1 and u.is_active and u.deleted_at is null`, input.UserID).Scan(&employeeName, &employeeNIP, &sourceDivision, &sourcePlacement)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(response, http.StatusUnprocessableEntity, "Pegawai tidak aktif atau tidak ditemukan.", "user_unavailable")
		return
	}
	if err != nil {
		a.internalError(response, "manual deduction user", err)
		return
	}

	componentBytes, _ := json.Marshal([]map[string]any{{"source_field": ruleSource, "code": ruleCode, "label": ruleLabel, "rate": ruleRate}})
	late, early, shift, attendanceStatus, leaveType := manualRuleColumns(ruleSource, ruleCode)
	notes := strings.TrimSpace(input.Notes)
	var recordID int64
	var previousJSON []byte
	err = tx.QueryRow(request.Context(), `
		select id, to_jsonb(ar)
		from public.attendance_records ar
		where batch_id = $1 and user_id = $2 and work_date = $3
		for update`, batchID, input.UserID, workDate).Scan(&recordID, &previousJSON)
	action := "manual_create"
	created := false
	if errors.Is(err, pgx.ErrNoRows) {
		created = true
		var sourceRow int
		if err := tx.QueryRow(request.Context(), `select coalesce(max(source_row), 4) + 1 from public.attendance_records where batch_id = $1`, batchID).Scan(&sourceRow); err != nil {
			a.internalError(response, "manual deduction source row", err)
			return
		}
		err = tx.QueryRow(request.Context(), `
			insert into public.attendance_records (
			  batch_id, period_id, user_id, source_row, work_date,
			  late_code, early_leave_code, shift_code, attendance_status, leave_type,
			  notes, source_division, source_placement, deduction_rate, deduction_components,
			  record_source, last_manual_reason_code, last_manual_reason_detail, updated_by, updated_at
			) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			  nullif($11, ''), $12, $13, $14, $15::jsonb,
			  'manual', $16, $17, $18, now())
			returning id`, batchID, input.PeriodID, input.UserID, sourceRow, workDate,
			late, early, shift, attendanceStatus, leaveType, notes, sourceDivision, sourcePlacement,
			ruleRate, string(componentBytes), reason.Code, reason.Detail, actor.ID).Scan(&recordID)
	} else if err != nil {
		a.internalError(response, "manual deduction existing record", err)
		return
	} else {
		action = "manual_override"
		if conflict, conflictErr := attendanceHasAppeal(request, tx, recordID); conflictErr != nil {
			a.internalError(response, "manual deduction appeal check", conflictErr)
			return
		} else if conflict {
			writeError(response, http.StatusConflict, "Data tidak dapat diubah karena sudah memiliki banding.", "appeal_exists")
			return
		}
		_, err = tx.Exec(request.Context(), `
			update public.attendance_records
			set late_code = $2, early_leave_code = $3, shift_code = $4,
			    attendance_status = $5, leave_type = $6,
			    notes = coalesce(nullif($7, ''), notes), deduction_rate = $8, deduction_components = $9::jsonb,
			    record_source = case when record_source = 'manual' then 'manual' else 'import_adjusted' end,
			    last_manual_reason_code = $10, last_manual_reason_detail = $11,
			    updated_by = $12, updated_at = now()
			where id = $1`, recordID, late, early, shift, attendanceStatus, leaveType,
			notes, ruleRate, string(componentBytes), reason.Code, reason.Detail, actor.ID)
	}
	if err != nil {
		a.internalError(response, "manual deduction save", err)
		return
	}
	var newJSON []byte
	if err := tx.QueryRow(request.Context(), `select to_jsonb(ar) from public.attendance_records ar where id = $1`, recordID).Scan(&newJSON); err != nil {
		a.internalError(response, "manual deduction snapshot", err)
		return
	}
	if _, err := tx.Exec(request.Context(), `
		insert into public.attendance_record_changes (
		  attendance_record_id, action, reason_code, reason_detail,
		  previous_values, new_values, changed_by
		) values ($1, $2, $3, $4, $5::jsonb, $6::jsonb, $7)`,
		recordID, action, reason.Code, reason.Detail, nullableJSON(previousJSON), string(newJSON), actor.ID); err != nil {
		a.internalError(response, "manual deduction history", err)
		return
	}
	if _, err := tx.Exec(request.Context(), `
		insert into public.audit_logs (actor_id, action, entity_type, entity_id, metadata, ip_address)
		values ($1, $2, 'attendance_record', $3, jsonb_build_object(
		  'period', $4, 'nip', $5, 'date', $6, 'rate', $7, 'reason', $8
		), nullif($9, '')::inet)`, actor.ID, "attendance."+action, strconv.FormatInt(recordID, 10), periodLabel,
		employeeNIP, workDate.Format("2006-01-02"), ruleRate, reason.Code, auth.ClientIP(request, a.cfg.TrustProxy)); err != nil {
		a.internalError(response, "manual deduction audit", err)
		return
	}
	if _, err := tx.Exec(request.Context(), `
		insert into public.notifications (user_id, kind, title, body, action_url)
		values ($1, 'deduction_adjusted', 'Data potongan diperbarui', $2, '/#deductions')`,
		input.UserID, fmt.Sprintf("Data %s tanggal %s diperbarui oleh administrator.", periodLabel, workDate.Format("02-01-2006"))); err != nil {
		a.internalError(response, "manual deduction notify", err)
		return
	}
	if err := tx.Commit(request.Context()); err != nil {
		a.internalError(response, "manual deduction commit", err)
		return
	}
	statusCode := http.StatusOK
	if created {
		statusCode = http.StatusCreated
	}
	writeJSON(response, statusCode, map[string]any{
		"id": recordID, "created": created,
		"message": fmt.Sprintf("Data %s untuk %s berhasil disimpan.", ruleLabel, employeeName),
	})
}

func (a *App) adminEditDeduction(response http.ResponseWriter, request *http.Request, actor auth.Principal) {
	recordID, err := strconv.ParseInt(request.PathValue("id"), 10, 64)
	if err != nil || recordID <= 0 {
		writeError(response, http.StatusBadRequest, "Data potongan tidak valid.", "invalid_id")
		return
	}
	var input struct {
		Rate         float64 `json:"rate"`
		ReasonCode   string  `json:"reason_code"`
		ReasonDetail string  `json:"reason_detail"`
	}
	if !decodeJSON(response, request, &input) {
		return
	}
	reason, reasonErr := validateManualReason(input.ReasonCode, input.ReasonDetail)
	if reasonErr != nil || input.Rate < 0 || input.Rate > 1 {
		message := "Potongan harus berada pada rentang 0% sampai 100%."
		if reasonErr != nil {
			message = reasonErr.Error()
		}
		writeError(response, http.StatusUnprocessableEntity, message, "invalid_manual_edit")
		return
	}
	tx, err := a.pool.BeginTx(request.Context(), pgx.TxOptions{})
	if err != nil {
		a.internalError(response, "manual edit begin", err)
		return
	}
	defer tx.Rollback(request.Context())
	var userID, nip, name, periodLabel string
	var workDate time.Time
	var previousJSON, components []byte
	err = tx.QueryRow(request.Context(), `
		select ar.user_id::text, u.nip, u.name, rp.label, ar.work_date,
		       to_jsonb(ar), ar.deduction_components
		from public.attendance_records ar
		join public.reporting_periods rp on rp.id = ar.period_id and rp.published_batch_id = ar.batch_id
		join public.users u on u.id = ar.user_id
		where ar.id = $1
		for update of ar`, recordID).Scan(&userID, &nip, &name, &periodLabel, &workDate, &previousJSON, &components)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(response, http.StatusNotFound, "Data potongan tidak ditemukan pada periode aktif.", "deduction_not_found")
		return
	}
	if err != nil {
		a.internalError(response, "manual edit record", err)
		return
	}
	if conflict, conflictErr := attendanceHasAppeal(request, tx, recordID); conflictErr != nil {
		a.internalError(response, "manual edit appeal check", conflictErr)
		return
	} else if conflict {
		writeError(response, http.StatusConflict, "Data tidak dapat diubah karena sudah memiliki banding.", "appeal_exists")
		return
	}
	adjustedComponents := manualAdjustedComponents(components, input.Rate, manualReasonLabel(reason.Code))
	_, err = tx.Exec(request.Context(), `
		update public.attendance_records
		set deduction_rate = $2, deduction_components = $3::jsonb,
		    record_source = case when record_source = 'manual' then 'manual' else 'import_adjusted' end,
		    last_manual_reason_code = $4, last_manual_reason_detail = $5,
		    updated_by = $6, updated_at = now()
		where id = $1`, recordID, input.Rate, string(adjustedComponents), reason.Code, reason.Detail, actor.ID)
	if err != nil {
		a.internalError(response, "manual edit save", err)
		return
	}
	var newJSON []byte
	if err := tx.QueryRow(request.Context(), `select to_jsonb(ar) from public.attendance_records ar where id = $1`, recordID).Scan(&newJSON); err != nil {
		a.internalError(response, "manual edit snapshot", err)
		return
	}
	if _, err := tx.Exec(request.Context(), `
		insert into public.attendance_record_changes (
		  attendance_record_id, action, reason_code, reason_detail,
		  previous_values, new_values, changed_by
		) values ($1, 'manual_edit', $2, $3, $4::jsonb, $5::jsonb, $6)`,
		recordID, reason.Code, reason.Detail, string(previousJSON), string(newJSON), actor.ID); err != nil {
		a.internalError(response, "manual edit history", err)
		return
	}
	if _, err := tx.Exec(request.Context(), `
		insert into public.audit_logs (actor_id, action, entity_type, entity_id, metadata, ip_address)
		values ($1, 'attendance.manual_edit', 'attendance_record', $2,
		  jsonb_build_object('period', $3, 'nip', $4, 'date', $5, 'new_rate', $6, 'reason', $7),
		  nullif($8, '')::inet)`, actor.ID, strconv.FormatInt(recordID, 10), periodLabel, nip,
		workDate.Format("2006-01-02"), input.Rate, reason.Code, auth.ClientIP(request, a.cfg.TrustProxy)); err != nil {
		a.internalError(response, "manual edit audit", err)
		return
	}
	if _, err := tx.Exec(request.Context(), `
		insert into public.notifications (user_id, kind, title, body, action_url)
		values ($1, 'deduction_adjusted', 'Potongan dikoreksi administrator', $2, '/#deductions')`,
		userID, fmt.Sprintf("Potongan tanggal %s pada periode %s telah dikoreksi.", workDate.Format("02-01-2006"), periodLabel)); err != nil {
		a.internalError(response, "manual edit notify", err)
		return
	}
	if err := tx.Commit(request.Context()); err != nil {
		a.internalError(response, "manual edit commit", err)
		return
	}
	writeJSON(response, http.StatusOK, map[string]any{"message": fmt.Sprintf("Potongan %s berhasil dikoreksi.", name)})
}

func (a *App) adminDeletePeriodDeductions(response http.ResponseWriter, request *http.Request, actor auth.Principal) {
	periodID := request.PathValue("id")
	if !validUUID(periodID) {
		writeError(response, http.StatusBadRequest, "Periode tidak valid.", "invalid_period")
		return
	}
	var input struct {
		Confirmation string `json:"confirmation"`
		Reason       string `json:"reason"`
	}
	if !decodeJSON(response, request, &input) {
		return
	}
	input.Reason = strings.TrimSpace(input.Reason)
	if len([]rune(input.Reason)) < 3 || len([]rune(input.Reason)) > 500 {
		writeError(response, http.StatusUnprocessableEntity, "Alasan penghapusan harus 3-500 karakter.", "invalid_delete_reason")
		return
	}
	tx, err := a.pool.BeginTx(request.Context(), pgx.TxOptions{})
	if err != nil {
		a.internalError(response, "delete period begin", err)
		return
	}
	defer tx.Rollback(request.Context())
	var label, batchID string
	err = tx.QueryRow(request.Context(), `
		select label, published_batch_id::text
		from public.reporting_periods
		where id = $1 and published_batch_id is not null
		for update`, periodID).Scan(&label, &batchID)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(response, http.StatusNotFound, "Periode tidak ditemukan atau sudah dihapus.", "period_not_found")
		return
	}
	if err != nil {
		a.internalError(response, "delete period lookup", err)
		return
	}
	expected := periodDeleteConfirmation(label)
	if strings.TrimSpace(input.Confirmation) != expected {
		writeError(response, http.StatusUnprocessableEntity, "Teks konfirmasi tidak sesuai. Ketik persis: "+expected, "confirmation_mismatch")
		return
	}
	pathRows, err := tx.Query(request.Context(), `
		select ad.storage_path
		from public.appeal_documents ad
		join public.appeal_items ai on ai.id = ad.appeal_item_id
		join public.appeals ap on ap.id = ai.appeal_id
		where ap.period_id = $1`, periodID)
	if err != nil {
		a.internalError(response, "delete period documents", err)
		return
	}
	storagePaths := []string{}
	for pathRows.Next() {
		var path string
		if err := pathRows.Scan(&path); err != nil {
			pathRows.Close()
			a.internalError(response, "delete period document scan", err)
			return
		}
		storagePaths = append(storagePaths, path)
	}
	if err := pathRows.Err(); err != nil {
		pathRows.Close()
		a.internalError(response, "delete period document rows", err)
		return
	}
	pathRows.Close()
	var attendanceCount, appealCount, batchCount int
	if err := tx.QueryRow(request.Context(), `select count(*) from public.attendance_records where period_id = $1`, periodID).Scan(&attendanceCount); err != nil {
		a.internalError(response, "delete period attendance count", err)
		return
	}
	if err := tx.QueryRow(request.Context(), `select count(*) from public.appeals where period_id = $1`, periodID).Scan(&appealCount); err != nil {
		a.internalError(response, "delete period appeal count", err)
		return
	}
	if err := tx.QueryRow(request.Context(), `select count(*) from public.import_batches where period_id = $1`, periodID).Scan(&batchCount); err != nil {
		a.internalError(response, "delete period batch count", err)
		return
	}
	if _, err := tx.Exec(request.Context(), `
		update public.notification_jobs set status = 'cancelled'
		where template_code = 'period_published'
		  and (payload->>'batch_id' = $1 or payload->>'period' = $2)
		  and status in ('pending', 'processing', 'failed')`, batchID, label); err != nil {
		a.internalError(response, "delete period jobs", err)
		return
	}
	if _, err := tx.Exec(request.Context(), `
		delete from public.notifications
		where (kind = 'period_published' and title = 'Data periode ' || $1 || ' tersedia')
		   or (kind = 'deduction_adjusted' and body ilike '%' || $1 || '%')`, label); err != nil {
		a.internalError(response, "delete period notifications", err)
		return
	}
	if _, err := tx.Exec(request.Context(), `delete from public.appeals where period_id = $1`, periodID); err != nil {
		a.internalError(response, "delete period appeals", err)
		return
	}
	if _, err := tx.Exec(request.Context(), `update public.reporting_periods set published_batch_id = null, updated_at = now() where id = $1`, periodID); err != nil {
		a.internalError(response, "delete period unpublish", err)
		return
	}
	if _, err := tx.Exec(request.Context(), `delete from public.attendance_records where period_id = $1`, periodID); err != nil {
		a.internalError(response, "delete period attendance", err)
		return
	}
	if _, err := tx.Exec(request.Context(), `delete from public.import_batches where period_id = $1`, periodID); err != nil {
		a.internalError(response, "delete period batches", err)
		return
	}
	if _, err := tx.Exec(request.Context(), `
		insert into public.audit_logs (actor_id, action, entity_type, entity_id, metadata, ip_address)
		values ($1, 'period.delete_deductions', 'reporting_period', $2,
		  jsonb_build_object('period', $3, 'reason', $4, 'attendance_rows', $5, 'appeals', $6, 'batches', $7),
		  nullif($8, '')::inet)`, actor.ID, periodID, label, input.Reason, attendanceCount, appealCount,
		batchCount, auth.ClientIP(request, a.cfg.TrustProxy)); err != nil {
		a.internalError(response, "delete period audit", err)
		return
	}
	if err := tx.Commit(request.Context()); err != nil {
		a.internalError(response, "delete period commit", err)
		return
	}
	storageWarning := ""
	if err := a.storage.Delete(request.Context(), storagePaths); err != nil {
		storageWarning = "Metadata sudah dihapus, tetapi sebagian berkas banding perlu dibersihkan dari Storage secara manual."
		a.log.Error("delete period storage cleanup", "period_id", periodID, "error", err)
	}
	writeJSON(response, http.StatusOK, map[string]any{
		"message":         fmt.Sprintf("Seluruh data potongan periode %s berhasil dihapus.", label),
		"attendance_rows": attendanceCount, "appeals": appealCount, "batches": batchCount,
		"storage_warning": storageWarning,
	})
}

func (a *App) publishedPeriods(request *http.Request) ([]periodInfo, error) {
	rows, err := a.pool.Query(request.Context(), `
		select id::text, label, period_start, period_end
		from public.reporting_periods
		where published_batch_id is not null
		order by period_end desc, updated_at desc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	periods := []periodInfo{}
	for rows.Next() {
		var item periodInfo
		var start, end time.Time
		if err := rows.Scan(&item.ID, &item.Label, &start, &end); err != nil {
			return nil, err
		}
		item.Start, item.End = start.Format("2006-01-02"), end.Format("2006-01-02")
		periods = append(periods, item)
	}
	return periods, rows.Err()
}

func validateManualReason(code, detail string) (manualReason, error) {
	code = strings.ToLower(strings.TrimSpace(code))
	detail = strings.TrimSpace(detail)
	if code != "input_error" && code != "system_error" && code != "other" {
		return manualReason{}, errors.New("Pilih alasan: Kesalahan input, Kesalahan sistem, atau Lainnya")
	}
	result := manualReason{Code: code}
	if code == "other" {
		if len([]rune(detail)) < 3 || len([]rune(detail)) > 500 {
			return manualReason{}, errors.New("Alasan lainnya harus diisi 3-500 karakter")
		}
		result.Detail = &detail
	}
	return result, nil
}

func manualReasonLabel(code string) string {
	switch code {
	case "input_error":
		return "Kesalahan input"
	case "system_error":
		return "Kesalahan sistem"
	default:
		return "Lainnya"
	}
}

func manualRuleColumns(source, code string) (any, any, any, any, any) {
	var late, early, shift, attendanceStatus, leaveType any
	switch source {
	case "late":
		late = code
	case "early_leave":
		early = code
	case "shift":
		shift = code
	case "status":
		attendanceStatus = code
	case "leave":
		leaveType = code
	}
	return late, early, shift, attendanceStatus, leaveType
}

func attendanceHasAppeal(request *http.Request, tx pgx.Tx, recordID int64) (bool, error) {
	var exists bool
	err := tx.QueryRow(request.Context(), `select exists(select 1 from public.appeal_items where attendance_record_id = $1)`, recordID).Scan(&exists)
	return exists, err
}

func manualAdjustedComponents(raw []byte, rate float64, reasonLabel string) []byte {
	var components []map[string]any
	if err := json.Unmarshal(raw, &components); err != nil || len(components) == 0 {
		components = []map[string]any{{"source_field": "manual", "code": "KOREKSI", "label": "Koreksi manual: " + reasonLabel, "rate": rate}}
	} else {
		for index := range components {
			components[index]["rate"] = 0.0
		}
		components[0]["rate"] = rate
		components[0]["label"] = fmt.Sprint(components[0]["label"]) + " (dikoreksi admin)"
	}
	result, _ := json.Marshal(components)
	return result
}

func nullableJSON(value []byte) any {
	if len(value) == 0 {
		return nil
	}
	return string(value)
}

func periodDeleteConfirmation(label string) string {
	return "HAPUS DATA " + strings.ToUpper(strings.TrimSpace(label))
}
