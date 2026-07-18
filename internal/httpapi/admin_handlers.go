package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bcpriok/pantas/internal/auth"
	"github.com/bcpriok/pantas/internal/importer"
	"github.com/jackc/pgx/v5"
)

func (a *App) adminUsers(response http.ResponseWriter, request *http.Request, _ auth.Principal) {
	queryText := strings.TrimSpace(request.URL.Query().Get("q"))
	page := parsePositiveInt(request.URL.Query().Get("page"), 1)
	limit := min(parsePositiveInt(request.URL.Query().Get("limit"), 50), 100)
	offset := (page - 1) * limit
	rows, err := a.pool.Query(request.Context(), `
		select u.id::text, u.nip, u.name, u.position_role, u.is_active,
		       u.must_change_password, coalesce(u.email, ''), u.email_verified_at is not null,
		       coalesce(u.phone_e164, ''), u.phone_verified_at is not null,
		       un.id::text, un.name, un.unit_type, count(*) over()
		from public.users u join public.units un on un.id = u.unit_id
		where u.deleted_at is null and ($1 = '' or u.nip ilike '%' || $1 || '%' or u.name ilike '%' || $1 || '%' or un.name ilike '%' || $1 || '%')
		order by u.is_active desc, u.name
		limit $2 offset $3`, queryText, limit, offset)
	if err != nil {
		a.internalError(response, "admin users", err)
		return
	}
	defer rows.Close()
	items := []map[string]any{}
	total := 0
	for rows.Next() {
		var id, nip, name, role, email, phone, unitID, unitName, unitType string
		var active, mustChange, emailVerified, phoneVerified bool
		if err := rows.Scan(&id, &nip, &name, &role, &active, &mustChange, &email, &emailVerified, &phone, &phoneVerified, &unitID, &unitName, &unitType, &total); err != nil {
			a.internalError(response, "admin user scan", err)
			return
		}
		items = append(items, map[string]any{"id": id, "nip": nip, "name": name, "role": role, "is_active": active, "must_change_password": mustChange, "email": email, "email_verified": emailVerified, "phone": phone, "phone_verified": phoneVerified, "unit_id": unitID, "unit_name": unitName, "unit_type": unitType})
	}
	writeJSON(response, http.StatusOK, map[string]any{"items": items, "page": page, "limit": limit, "total": total})
}

func (a *App) adminCreateUser(response http.ResponseWriter, request *http.Request, actor auth.Principal) {
	var input struct {
		NIP    string `json:"nip"`
		Name   string `json:"name"`
		UnitID string `json:"unit_id"`
		Role   string `json:"role"`
	}
	if !decodeJSON(response, request, &input) {
		return
	}
	input.NIP, input.Name = strings.TrimSpace(input.NIP), strings.TrimSpace(input.Name)
	if !validNIPInput(input.NIP) || len([]rune(input.Name)) < 2 || len([]rune(input.Name)) > 200 || !validUUID(input.UnitID) || !validRole(input.Role) {
		writeError(response, http.StatusUnprocessableEntity, "Nama, NIP, unit, atau jabatan tidak valid.", "invalid_user")
		return
	}
	if err := a.validateRoleUnit(request, input.Role, input.UnitID); err != nil {
		writeError(response, http.StatusUnprocessableEntity, err.Error(), "invalid_role_unit")
		return
	}
	var id string
	err := a.pool.QueryRow(request.Context(), `
		insert into public.users (nip, name, unit_id, position_role, password_hash, must_change_password)
		values ($1, $2, $3, $4, null, true) returning id::text`, input.NIP, input.Name, input.UnitID, input.Role).Scan(&id)
	if err != nil {
		if isConstraintError(err) {
			writeError(response, http.StatusConflict, "NIP atau posisi kepala unit sudah digunakan.", "user_conflict")
			return
		}
		a.internalError(response, "admin create user", err)
		return
	}
	a.audit(request, actor.ID, "user.create", "user", id, map[string]any{"nip": input.NIP, "role": input.Role, "unit_id": input.UnitID})
	writeJSON(response, http.StatusCreated, map[string]any{"id": id, "message": "Pengguna dibuat. Password awal adalah NIP dan wajib diganti saat login pertama."})
}

func (a *App) adminUpdateUser(response http.ResponseWriter, request *http.Request, actor auth.Principal) {
	userID := request.PathValue("id")
	if !validUUID(userID) {
		writeError(response, http.StatusBadRequest, "Pengguna tidak valid.", "invalid_id")
		return
	}
	var input struct {
		Name     *string `json:"name"`
		UnitID   *string `json:"unit_id"`
		Role     *string `json:"role"`
		IsActive *bool   `json:"is_active"`
		Reason   string  `json:"reason"`
	}
	if !decodeJSON(response, request, &input) {
		return
	}
	tx, err := a.pool.BeginTx(request.Context(), pgx.TxOptions{})
	if err != nil {
		a.internalError(response, "admin user begin", err)
		return
	}
	defer tx.Rollback(request.Context())
	var oldName, oldUnit, oldRole string
	var oldActive bool
	err = tx.QueryRow(request.Context(), `select name, unit_id::text, position_role, is_active from public.users where id = $1 and deleted_at is null for update`, userID).Scan(&oldName, &oldUnit, &oldRole, &oldActive)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(response, http.StatusNotFound, "Pengguna tidak ditemukan.", "not_found")
		return
	}
	if err != nil {
		a.internalError(response, "admin user lookup", err)
		return
	}
	name, unitID, role, active := oldName, oldUnit, oldRole, oldActive
	if input.Name != nil {
		name = strings.TrimSpace(*input.Name)
	}
	if input.UnitID != nil {
		unitID = *input.UnitID
	}
	if input.Role != nil {
		role = *input.Role
	}
	if input.IsActive != nil {
		active = *input.IsActive
	}
	if len([]rune(name)) < 2 || len([]rune(name)) > 200 || !validUUID(unitID) || !validRole(role) {
		writeError(response, http.StatusUnprocessableEntity, "Data pengguna tidak valid.", "invalid_user")
		return
	}
	if err := a.validateRoleUnitTx(request, tx, role, unitID); err != nil {
		writeError(response, http.StatusUnprocessableEntity, err.Error(), "invalid_role_unit")
		return
	}
	_, err = tx.Exec(request.Context(), `
		update public.users set name = $2, unit_id = $3, position_role = $4, is_admin = false, is_active = $5
		where id = $1`, userID, name, unitID, role, active)
	if err != nil {
		if isConstraintError(err) {
			writeError(response, http.StatusConflict, "Posisi kepala pada unit tersebut sudah terisi.", "user_conflict")
			return
		}
		a.internalError(response, "admin user update", err)
		return
	}
	if oldUnit != unitID || oldRole != role {
		if _, err := tx.Exec(request.Context(), `
			insert into public.user_assignment_history (user_id, previous_unit_id, new_unit_id, previous_role, new_role, changed_by, reason)
			values ($1, $2, $3, $4, $5, $6, nullif($7, ''))`, userID, oldUnit, unitID, oldRole, role, actor.ID, strings.TrimSpace(input.Reason)); err != nil {
			a.internalError(response, "admin assignment history", err)
			return
		}
	}
	if !active {
		if _, err := tx.Exec(request.Context(), `update public.sessions set revoked_at = now() where user_id = $1 and revoked_at is null`, userID); err != nil {
			a.internalError(response, "admin revoke sessions", err)
			return
		}
	}
	if err := tx.Commit(request.Context()); err != nil {
		a.internalError(response, "admin user commit", err)
		return
	}
	a.audit(request, actor.ID, "user.update", "user", userID, map[string]any{"unit_before": oldUnit, "unit_after": unitID, "role_before": oldRole, "role_after": role, "active": active})
	writeJSON(response, http.StatusOK, map[string]any{"message": "Data pengguna diperbarui."})
}

func (a *App) adminDeleteUser(response http.ResponseWriter, request *http.Request, actor auth.Principal) {
	userID := request.PathValue("id")
	if !validUUID(userID) {
		writeError(response, http.StatusBadRequest, "Pengguna tidak valid.", "invalid_id")
		return
	}
	tx, err := a.pool.BeginTx(request.Context(), pgx.TxOptions{})
	if err != nil {
		a.internalError(response, "delete user begin", err)
		return
	}
	defer tx.Rollback(request.Context())
	var nip string
	err = tx.QueryRow(request.Context(), `select nip from public.users where id = $1 and deleted_at is null for update`, userID).Scan(&nip)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(response, http.StatusNotFound, "Pengguna tidak ditemukan.", "not_found")
		return
	}
	if err != nil {
		a.internalError(response, "delete user lookup", err)
		return
	}
	if _, err := tx.Exec(request.Context(), `update public.users set is_active = false, deleted_at = now(), email = null, phone_e164 = null where id = $1`, userID); err != nil {
		a.internalError(response, "delete user", err)
		return
	}
	if _, err := tx.Exec(request.Context(), `update public.sessions set revoked_at = now() where user_id = $1 and revoked_at is null`, userID); err != nil {
		a.internalError(response, "delete sessions", err)
		return
	}
	if err := tx.Commit(request.Context()); err != nil {
		a.internalError(response, "delete user commit", err)
		return
	}
	a.audit(request, actor.ID, "user.delete", "user", userID, map[string]any{"nip": nip, "mode": "soft_delete"})
	response.WriteHeader(http.StatusNoContent)
}

func (a *App) adminResetUserPassword(response http.ResponseWriter, request *http.Request, actor auth.Principal) {
	userID := request.PathValue("id")
	if !validUUID(userID) {
		writeError(response, http.StatusBadRequest, "Pengguna tidak valid.", "invalid_id")
		return
	}
	command, err := a.pool.Exec(request.Context(), `
		update public.users set password_hash = null, must_change_password = true
		where id = $1 and deleted_at is null`, userID)
	if err != nil {
		a.internalError(response, "reset user password", err)
		return
	}
	if command.RowsAffected() == 0 {
		writeError(response, http.StatusNotFound, "Pengguna tidak ditemukan.", "not_found")
		return
	}
	_, _ = a.pool.Exec(request.Context(), `update public.sessions set revoked_at = now() where user_id = $1 and revoked_at is null`, userID)
	a.audit(request, actor.ID, "user.password_reset", "user", userID, map[string]any{})
	writeJSON(response, http.StatusOK, map[string]any{"message": "Password dikembalikan menjadi NIP dan wajib diganti saat login."})
}

func (a *App) adminUnits(response http.ResponseWriter, request *http.Request, _ auth.Principal) {
	rows, err := a.pool.Query(request.Context(), `
		select u.id::text, u.code, u.name, coalesce(u.source_name, ''), u.unit_type,
		       coalesce(u.parent_id::text, ''), coalesce(p.name, ''), u.sort_order, u.is_active,
		       (select count(*) from public.users x where x.unit_id = u.id and x.is_active and x.deleted_at is null)
		from public.units u left join public.units p on p.id = u.parent_id
		order by coalesce(p.sort_order, u.sort_order), u.unit_type, u.sort_order, u.name`)
	if err != nil {
		a.internalError(response, "admin units", err)
		return
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		var id, code, name, source, unitType, parentID, parentName string
		var order, members int
		var active bool
		if err := rows.Scan(&id, &code, &name, &source, &unitType, &parentID, &parentName, &order, &active, &members); err != nil {
			a.internalError(response, "admin unit scan", err)
			return
		}
		items = append(items, map[string]any{"id": id, "code": code, "name": name, "source_name": source, "type": unitType, "parent_id": parentID, "parent_name": parentName, "sort_order": order, "is_active": active, "members": members})
	}
	writeJSON(response, http.StatusOK, map[string]any{"items": items})
}

func (a *App) adminParameters(response http.ResponseWriter, request *http.Request, _ auth.Principal) {
	rows, err := a.pool.Query(request.Context(), `select key, category, label, coalesce(description, ''), value_json, value_type, updated_at from public.parameters order by category, label`)
	if err != nil {
		a.internalError(response, "admin parameters", err)
		return
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		var key, category, label, description, valueType string
		var value []byte
		var updated time.Time
		if err := rows.Scan(&key, &category, &label, &description, &value, &valueType, &updated); err != nil {
			a.internalError(response, "parameter scan", err)
			return
		}
		items = append(items, map[string]any{"key": key, "category": category, "label": label, "description": description, "value": json.RawMessage(value), "value_type": valueType, "updated_at": updated})
	}
	writeJSON(response, http.StatusOK, map[string]any{"items": items})
}

func (a *App) adminUpdateParameter(response http.ResponseWriter, request *http.Request, actor auth.Principal) {
	key := request.PathValue("key")
	var input struct {
		Value json.RawMessage `json:"value"`
	}
	if !decodeJSON(response, request, &input) {
		return
	}
	if len(input.Value) == 0 || !json.Valid(input.Value) {
		writeError(response, http.StatusUnprocessableEntity, "Nilai parameter bukan JSON valid.", "invalid_parameter")
		return
	}
	var valueType string
	if err := a.pool.QueryRow(request.Context(), `select value_type from public.parameters where key = $1`, key).Scan(&valueType); errors.Is(err, pgx.ErrNoRows) {
		writeError(response, http.StatusNotFound, "Parameter tidak ditemukan.", "not_found")
		return
	} else if err != nil {
		a.internalError(response, "parameter lookup", err)
		return
	}
	if err := validateParameterValue(valueType, input.Value); err != nil {
		writeError(response, http.StatusUnprocessableEntity, err.Error(), "invalid_parameter")
		return
	}
	_, err := a.pool.Exec(request.Context(), `update public.parameters set value_json = $2::jsonb, updated_by = $3 where key = $1`, key, string(input.Value), actor.ID)
	if err != nil {
		a.internalError(response, "parameter update", err)
		return
	}
	a.audit(request, actor.ID, "parameter.update", "parameter", key, map[string]any{"value": json.RawMessage(input.Value)})
	response.WriteHeader(http.StatusNoContent)
}

func (a *App) adminRules(response http.ResponseWriter, request *http.Request, _ auth.Principal) {
	rows, err := a.pool.Query(request.Context(), `select id::text, source_field, code, label, rate, is_active, sort_order from public.deduction_rules order by sort_order, source_field, code`)
	if err != nil {
		a.internalError(response, "admin rules", err)
		return
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		var id, source, code, label string
		var rate float64
		var active bool
		var order int
		if err := rows.Scan(&id, &source, &code, &label, &rate, &active, &order); err != nil {
			a.internalError(response, "rule scan", err)
			return
		}
		items = append(items, map[string]any{"id": id, "source": source, "code": code, "label": label, "rate": rate, "is_active": active, "sort_order": order})
	}
	writeJSON(response, http.StatusOK, map[string]any{"items": items})
}

func (a *App) adminUpdateRule(response http.ResponseWriter, request *http.Request, actor auth.Principal) {
	id := request.PathValue("id")
	var input struct {
		Label    string  `json:"label"`
		Rate     float64 `json:"rate"`
		IsActive bool    `json:"is_active"`
	}
	if !validUUID(id) || !decodeJSON(response, request, &input) {
		return
	}
	input.Label = strings.TrimSpace(input.Label)
	if len([]rune(input.Label)) < 2 || len([]rune(input.Label)) > 200 || input.Rate < 0 || input.Rate > 1 {
		writeError(response, http.StatusUnprocessableEntity, "Label atau persentase tidak valid.", "invalid_rule")
		return
	}
	command, err := a.pool.Exec(request.Context(), `update public.deduction_rules set label = $2, rate = $3, is_active = $4 where id = $1`, id, input.Label, input.Rate, input.IsActive)
	if err != nil {
		a.internalError(response, "rule update", err)
		return
	}
	if command.RowsAffected() == 0 {
		writeError(response, http.StatusNotFound, "Aturan tidak ditemukan.", "not_found")
		return
	}
	a.audit(request, actor.ID, "deduction_rule.update", "deduction_rule", id, map[string]any{"label": input.Label, "rate": input.Rate, "active": input.IsActive})
	response.WriteHeader(http.StatusNoContent)
}

func (a *App) adminReasons(response http.ResponseWriter, request *http.Request, _ auth.Principal) {
	rows, err := a.pool.Query(request.Context(), `select id::text, code, label, coalesce(description, ''), is_active, sort_order from public.appeal_reason_categories order by sort_order, label`)
	if err != nil {
		a.internalError(response, "admin reasons", err)
		return
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		var id, code, label, description string
		var active bool
		var order int
		if err := rows.Scan(&id, &code, &label, &description, &active, &order); err != nil {
			a.internalError(response, "reason scan", err)
			return
		}
		items = append(items, map[string]any{"id": id, "code": code, "label": label, "description": description, "is_active": active, "sort_order": order})
	}
	writeJSON(response, http.StatusOK, map[string]any{"items": items})
}

func (a *App) adminCreateReason(response http.ResponseWriter, request *http.Request, actor auth.Principal) {
	var input struct {
		Code        string `json:"code"`
		Label       string `json:"label"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
	}
	if !decodeJSON(response, request, &input) {
		return
	}
	input.Code = strings.ToLower(strings.TrimSpace(input.Code))
	input.Label = strings.TrimSpace(input.Label)
	if !validCode(input.Code) || len([]rune(input.Label)) < 2 || len([]rune(input.Description)) > 1000 {
		writeError(response, http.StatusUnprocessableEntity, "Kode, label, atau deskripsi tidak valid.", "invalid_reason")
		return
	}
	var id string
	err := a.pool.QueryRow(request.Context(), `insert into public.appeal_reason_categories (code, label, description, sort_order) values ($1, $2, nullif($3, ''), $4) returning id::text`, input.Code, input.Label, strings.TrimSpace(input.Description), input.SortOrder).Scan(&id)
	if err != nil {
		if isConstraintError(err) {
			writeError(response, http.StatusConflict, "Kode alasan sudah digunakan.", "reason_conflict")
			return
		}
		a.internalError(response, "reason create", err)
		return
	}
	a.audit(request, actor.ID, "appeal_reason.create", "appeal_reason", id, map[string]any{"code": input.Code})
	writeJSON(response, http.StatusCreated, map[string]any{"id": id})
}

func (a *App) adminUpdateReason(response http.ResponseWriter, request *http.Request, actor auth.Principal) {
	id := request.PathValue("id")
	var input struct {
		Label, Description string
		IsActive           bool `json:"is_active"`
		SortOrder          int  `json:"sort_order"`
	}
	if !validUUID(id) || !decodeJSON(response, request, &input) {
		return
	}
	input.Label = strings.TrimSpace(input.Label)
	if len([]rune(input.Label)) < 2 || len([]rune(input.Description)) > 1000 {
		writeError(response, http.StatusUnprocessableEntity, "Label atau deskripsi tidak valid.", "invalid_reason")
		return
	}
	command, err := a.pool.Exec(request.Context(), `update public.appeal_reason_categories set label=$2, description=nullif($3,''), is_active=$4, sort_order=$5 where id=$1`, id, input.Label, strings.TrimSpace(input.Description), input.IsActive, input.SortOrder)
	if err != nil {
		a.internalError(response, "reason update", err)
		return
	}
	if command.RowsAffected() == 0 {
		writeError(response, http.StatusNotFound, "Alasan tidak ditemukan.", "not_found")
		return
	}
	a.audit(request, actor.ID, "appeal_reason.update", "appeal_reason", id, map[string]any{"active": input.IsActive})
	response.WriteHeader(http.StatusNoContent)
}

func (a *App) adminImports(response http.ResponseWriter, request *http.Request, _ auth.Principal) {
	rows, err := a.pool.Query(request.Context(), `
		select ib.id::text, rp.label, rp.period_start, rp.period_end, ib.version, ib.original_filename,
		       ib.integrity_status, ib.status, ib.row_count, ib.employee_count, ib.deduction_day_count,
		       ib.total_deduction_rate, ib.created_at, ib.published_at,
		       coalesce(creator_user.name, creator_admin.name, 'Akun tidak tersedia')
		from public.import_batches ib
		join public.reporting_periods rp on rp.id = ib.period_id
		left join public.users creator_user on creator_user.id = ib.created_by
		left join public.admin_accounts creator_admin on creator_admin.account_id = ib.created_by
		order by ib.created_at desc limit 100`)
	if err != nil {
		a.internalError(response, "admin imports", err)
		return
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		var id, label, filename, integrity, status, creator string
		var start, end time.Time
		var version, rowCount, employees, deductionDays int
		var total float64
		var created time.Time
		var published *time.Time
		if err := rows.Scan(&id, &label, &start, &end, &version, &filename, &integrity, &status, &rowCount, &employees, &deductionDays, &total, &created, &published, &creator); err != nil {
			a.internalError(response, "import scan", err)
			return
		}
		items = append(items, map[string]any{"id": id, "label": label, "start": start.Format("2006-01-02"), "end": end.Format("2006-01-02"), "version": version, "filename": filename, "integrity": integrity, "status": status, "rows": rowCount, "employees": employees, "deduction_days": deductionDays, "total": total, "created_at": created, "published_at": published, "created_by": creator})
	}
	writeJSON(response, http.StatusOK, map[string]any{"items": items})
}

func (a *App) adminImportPreview(response http.ResponseWriter, request *http.Request, actor auth.Principal) {
	request.Body = http.MaxBytesReader(response, request.Body, a.cfg.MaxExcelBytes+1<<20)
	if err := request.ParseMultipartForm(a.cfg.MaxExcelBytes); err != nil {
		writeError(response, http.StatusRequestEntityTooLarge, "File melebihi batas atau form tidak valid.", "invalid_upload")
		return
	}
	file, header, err := request.FormFile("file")
	if err != nil {
		writeError(response, http.StatusBadRequest, "Pilih file Excel.", "file_required")
		return
	}
	defer file.Close()
	if strings.ToLower(filepath.Ext(header.Filename)) != ".xlsx" {
		writeError(response, http.StatusUnprocessableEntity, "File harus berformat .xlsx.", "invalid_file_type")
		return
	}
	data, err := io.ReadAll(io.LimitReader(file, a.cfg.MaxExcelBytes+1))
	if err != nil || int64(len(data)) > a.cfg.MaxExcelBytes {
		writeError(response, http.StatusRequestEntityTooLarge, "File Excel terlalu besar.", "file_too_large")
		return
	}
	preview, err := a.imports.PreviewAndStage(request.Context(), actor, sanitizeUploadName(header.Filename), data)
	if err != nil {
		var validation *importer.ValidationError
		switch {
		case errors.As(err, &validation):
			writeJSON(response, http.StatusUnprocessableEntity, map[string]any{"error": map[string]any{"code": "invalid_excel", "message": validation.Error(), "problems": validation.Problems}})
		case errors.Is(err, importer.ErrMissingUsers):
			writeJSON(response, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "missing_users", "message": err.Error()}, "preview": preview})
		default:
			a.internalError(response, "import preview", err)
		}
		return
	}
	writeJSON(response, http.StatusCreated, map[string]any{"preview": preview})
}

func (a *App) adminImportPublish(response http.ResponseWriter, request *http.Request, actor auth.Principal) {
	id := request.PathValue("id")
	if !validUUID(id) {
		writeError(response, http.StatusBadRequest, "Batch tidak valid.", "invalid_id")
		return
	}
	err := a.imports.Publish(request.Context(), actor, id, auth.ClientIP(request, a.cfg.TrustProxy))
	if err != nil {
		switch {
		case errors.Is(err, importer.ErrAlreadyPublished):
			writeError(response, http.StatusConflict, err.Error(), "duplicate_import")
		case errors.Is(err, importer.ErrAppealsExist):
			writeError(response, http.StatusConflict, err.Error(), "appeals_exist")
		default:
			a.internalError(response, "import publish", err)
		}
		return
	}
	writeJSON(response, http.StatusOK, map[string]any{"message": "Data periode berhasil dipublikasikan dan notifikasi dijadwalkan."})
}

func (a *App) adminImportReject(response http.ResponseWriter, request *http.Request, actor auth.Principal) {
	id := request.PathValue("id")
	if !validUUID(id) {
		writeError(response, http.StatusBadRequest, "Batch tidak valid.", "invalid_id")
		return
	}
	if err := a.imports.RejectDraft(request.Context(), actor, id, auth.ClientIP(request, a.cfg.TrustProxy)); err != nil {
		writeError(response, http.StatusConflict, err.Error(), "invalid_batch")
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

func (a *App) validateRoleUnit(request *http.Request, role, unitID string) error {
	var unitType string
	if err := a.pool.QueryRow(request.Context(), `select unit_type from public.units where id=$1 and is_active`, unitID).Scan(&unitType); err != nil {
		return errors.New("unit tidak ditemukan")
	}
	return roleUnitError(role, unitType)
}

func (a *App) validateRoleUnitTx(request *http.Request, tx pgx.Tx, role, unitID string) error {
	var unitType string
	if err := tx.QueryRow(request.Context(), `select unit_type from public.units where id=$1 and is_active`, unitID).Scan(&unitType); err != nil {
		return errors.New("unit tidak ditemukan")
	}
	return roleUnitError(role, unitType)
}

func roleUnitError(role, unitType string) error {
	valid := map[string]string{"section_head": "section", "division_head": "division", "office_head": "office", "functional": "functional"}
	if expected, ok := valid[role]; ok && unitType != expected {
		return fmt.Errorf("jabatan %s harus ditempatkan pada unit bertipe %s", role, expected)
	}
	if role == "staff" && unitType == "functional" {
		return errors.New("pegawai pada unit Fungsional harus memakai jabatan fungsional")
	}
	return nil
}

func validateParameterValue(valueType string, raw json.RawMessage) error {
	switch valueType {
	case "integer":
		var value int
		if json.Unmarshal(raw, &value) != nil || value < 0 || value > 100000000 {
			return errors.New("nilai harus bilangan bulat dalam rentang yang wajar")
		}
	case "decimal", "percent":
		var value float64
		if json.Unmarshal(raw, &value) != nil || value < 0 || value > 1000000 {
			return errors.New("nilai harus angka non-negatif")
		}
		if valueType == "percent" && value > 1 {
			return errors.New("persentase disimpan sebagai pecahan 0 sampai 1")
		}
	case "boolean":
		var value bool
		if json.Unmarshal(raw, &value) != nil {
			return errors.New("nilai harus true atau false")
		}
	case "json":
		var value any
		if json.Unmarshal(raw, &value) != nil {
			return errors.New("nilai JSON tidak valid")
		}
	default:
		return errors.New("tipe parameter tidak didukung")
	}
	return nil
}

func (a *App) audit(request *http.Request, actorID, action, entityType, entityID string, metadata map[string]any) {
	encoded, _ := json.Marshal(metadata)
	_, err := a.pool.Exec(request.Context(), `insert into public.audit_logs(actor_id,action,entity_type,entity_id,metadata,ip_address) values($1,$2,$3,nullif($4,''),$5::jsonb,nullif($6,'')::inet)`, actorID, action, entityType, entityID, string(encoded), auth.ClientIP(request, a.cfg.TrustProxy))
	if err != nil {
		a.log.Error("audit write", "action", action, "error", err)
	}
}

func parsePositiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return fallback
	}
	return parsed
}
func validNIPInput(value string) bool {
	if len(value) != 18 {
		return false
	}
	for _, char := range value {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}
func validRole(value string) bool {
	switch value {
	case "staff", "section_head", "division_head", "office_head", "functional":
		return true
	}
	return false
}
func validCode(value string) bool {
	if len(value) < 2 || len(value) > 60 {
		return false
	}
	for _, char := range value {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '_') {
			return false
		}
	}
	return true
}
func sanitizeUploadName(value string) string {
	value = filepath.Base(strings.TrimSpace(value))
	if len(value) > 200 {
		value = value[len(value)-200:]
	}
	return value
}
func isConstraintError(err error) bool {
	return strings.Contains(err.Error(), "SQLSTATE 23505") || strings.Contains(strings.ToLower(err.Error()), "duplicate key")
}
