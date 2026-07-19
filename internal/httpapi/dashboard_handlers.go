package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bcpriok/pantas/internal/auth"
	"github.com/jackc/pgx/v5"
)

type periodInfo struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Start string `json:"start"`
	End   string `json:"end"`
}

func (a *App) dashboard(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	period, err := a.currentPeriod(request)
	if errors.Is(err, pgx.ErrNoRows) {
		writeJSON(response, http.StatusOK, map[string]any{
			"current_period": nil, "summary": nil, "deductions": []any{},
			"unread_notifications": 0, "can_monitor": principal.IsSupervisor(),
		})
		return
	}
	if err != nil {
		a.internalError(response, "dashboard period", err)
		return
	}
	var original, effective float64
	var deductionDays, workDays, overtimeDays, leaveDays, offDays int
	err = a.pool.QueryRow(request.Context(), `
		select coalesce(original_deduction_rate, 0), coalesce(effective_deduction_rate, 0),
		       coalesce(deduction_days, 0), coalesce(work_days, 0), coalesce(overtime_days, 0),
		       coalesce(leave_days, 0), coalesce(off_days, 0)
		from public.monthly_user_summary
		where period_id = $1 and user_id = $2`, period.ID, principal.ID).Scan(
		&original, &effective, &deductionDays, &workDays, &overtimeDays, &leaveDays, &offDays,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		original, effective = 0, 0
		deductionDays, workDays, overtimeDays, leaveDays, offDays = 0, 0, 0, 0, 0
	} else if err != nil {
		a.internalError(response, "dashboard summary", err)
		return
	}
	var unread, pendingAppeal int
	if err := a.pool.QueryRow(request.Context(), `select count(*) from public.notifications where user_id = $1 and read_at is null`, principal.ID).Scan(&unread); err != nil {
		a.internalError(response, "dashboard notifications", err)
		return
	}
	if err := a.pool.QueryRow(request.Context(), `
		select count(*)
		from public.appeal_items ai
		join public.appeals ap on ap.id = ai.appeal_id
		where ap.user_id = $1 and ap.period_id = $2 and ai.admin_status = 'pending'`, principal.ID, period.ID).Scan(&pendingAppeal); err != nil {
		a.internalError(response, "dashboard appeal", err)
		return
	}
	items, err := a.loadDeductionRows(request, principal.ID, period.ID)
	if err != nil {
		a.internalError(response, "dashboard deduction rows", err)
		return
	}
	writeJSON(response, http.StatusOK, map[string]any{
		"current_period": period,
		"summary": map[string]any{
			"original_deduction": original, "effective_deduction": effective,
			"deduction_days": deductionDays, "work_days": workDays, "overtime_days": overtimeDays,
			"leave_days": leaveDays, "off_days": offDays,
		},
		"deductions":           items,
		"can_appeal":           original > 0,
		"pending_appeal_items": pendingAppeal,
		"unread_notifications": unread,
		"can_monitor":          principal.IsSupervisor(),
	})
}

func (a *App) history(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	fromValue, toValue := request.URL.Query().Get("from"), request.URL.Query().Get("to")
	if fromValue == "" && toValue == "" {
		var latest time.Time
		if err := a.pool.QueryRow(request.Context(), `select coalesce(max(period_end), current_date) from public.reporting_periods where published_batch_id is not null`).Scan(&latest); err != nil {
			a.internalError(response, "history latest period", err)
			return
		}
		toValue = latest.Format("2006-01")
	}
	from, to, err := monthRange(fromValue, toValue)
	if err != nil {
		writeError(response, http.StatusBadRequest, err.Error(), "invalid_range")
		return
	}
	rows, err := a.pool.Query(request.Context(), `
		select rp.id::text, rp.label, rp.period_start, rp.period_end,
		       coalesce(mus.original_deduction_rate, 0), coalesce(mus.effective_deduction_rate, 0),
		       coalesce(mus.deduction_days, 0)
		from public.reporting_periods rp
		left join public.monthly_user_summary mus on mus.period_id = rp.id and mus.user_id = $1
		where rp.published_batch_id is not null and rp.period_end between $2 and $3
		order by rp.period_end`, principal.ID, from, to)
	if err != nil {
		a.internalError(response, "history", err)
		return
	}
	defer rows.Close()
	points := []map[string]any{}
	for rows.Next() {
		var id, label string
		var start, end time.Time
		var original, effective float64
		var days int
		if err := rows.Scan(&id, &label, &start, &end, &original, &effective, &days); err != nil {
			a.internalError(response, "history scan", err)
			return
		}
		points = append(points, map[string]any{
			"period_id": id, "label": label, "start": start.Format("2006-01-02"), "end": end.Format("2006-01-02"),
			"original": original, "effective": effective, "deduction_days": days,
		})
	}
	if err := rows.Err(); err != nil {
		a.internalError(response, "history rows", err)
		return
	}
	writeJSON(response, http.StatusOK, map[string]any{"points": points, "from": from.Format("2006-01"), "to": to.Format("2006-01")})
}

func (a *App) deductions(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	periodID := request.URL.Query().Get("period_id")
	if periodID == "" {
		period, err := a.currentPeriod(request)
		if errors.Is(err, pgx.ErrNoRows) {
			writeJSON(response, http.StatusOK, map[string]any{"items": []any{}})
			return
		}
		if err != nil {
			a.internalError(response, "deductions period", err)
			return
		}
		periodID = period.ID
	}
	if !validUUID(periodID) {
		writeError(response, http.StatusBadRequest, "Periode tidak valid.", "invalid_period")
		return
	}
	items, err := a.loadDeductionRows(request, principal.ID, periodID)
	if err != nil {
		a.internalError(response, "deductions", err)
		return
	}
	writeJSON(response, http.StatusOK, map[string]any{"items": items})
}

func (a *App) notifications(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	rows, err := a.pool.Query(request.Context(), `
		select id::text, kind, title, body, coalesce(action_url, ''), created_at, read_at
		from public.notifications where user_id = $1 order by created_at desc limit 50`, principal.ID)
	if err != nil {
		a.internalError(response, "notifications", err)
		return
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		var id, kind, title, body, action string
		var created time.Time
		var readAt *time.Time
		if err := rows.Scan(&id, &kind, &title, &body, &action, &created, &readAt); err != nil {
			a.internalError(response, "notification scan", err)
			return
		}
		items = append(items, map[string]any{"id": id, "kind": kind, "title": title, "body": body, "action_url": action, "created_at": created, "read": readAt != nil})
	}
	writeJSON(response, http.StatusOK, map[string]any{"items": items})
}

func (a *App) markNotificationRead(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	id := request.PathValue("id")
	if !validUUID(id) {
		writeError(response, http.StatusBadRequest, "Notifikasi tidak valid.", "invalid_id")
		return
	}
	command, err := a.pool.Exec(request.Context(), `update public.notifications set read_at = coalesce(read_at, now()) where id = $1 and user_id = $2`, id, principal.ID)
	if err != nil {
		a.internalError(response, "notification read", err)
		return
	}
	if command.RowsAffected() == 0 {
		writeError(response, http.StatusNotFound, "Notifikasi tidak ditemukan.", "not_found")
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

func (a *App) markAllNotificationsRead(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	command, err := a.pool.Exec(request.Context(), `
		update public.notifications
		set read_at = now()
		where user_id = $1 and read_at is null`, principal.ID)
	if err != nil {
		a.internalError(response, "notifications read all", err)
		return
	}
	writeJSON(response, http.StatusOK, map[string]any{"read": command.RowsAffected()})
}

func (a *App) monitoring(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	period, err := a.currentPeriod(request)
	if errors.Is(err, pgx.ErrNoRows) {
		writeJSON(response, http.StatusOK, map[string]any{"current_period": nil, "mode": "empty", "items": []any{}})
		return
	}
	if err != nil {
		a.internalError(response, "monitoring period", err)
		return
	}
	unitID := request.URL.Query().Get("unit_id")
	mode, items, err := a.monitoringData(request, principal, period.ID, unitID)
	if err != nil {
		if errors.Is(err, errForbiddenScope) {
			writeError(response, http.StatusForbidden, "Detail unit ini tidak dapat dibuka.", "forbidden")
			return
		}
		a.internalError(response, "monitoring data", err)
		return
	}
	totals, err := a.monitoringTotals(request, principal, period.ID)
	if err != nil {
		a.internalError(response, "monitoring totals", err)
		return
	}
	writeJSON(response, http.StatusOK, map[string]any{"current_period": period, "mode": mode, "items": items, "totals": totals})
}

var errForbiddenScope = errors.New("forbidden scope")

func (a *App) monitoringData(request *http.Request, principal auth.Principal, periodID, requestedUnitID string) (string, []map[string]any, error) {
	ctx := request.Context()
	peopleMode := principal.PositionRole == "section_head"
	targetUnit := principal.UnitID
	targetUnitType := principal.UnitType
	if requestedUnitID != "" {
		if !validUUID(requestedUnitID) {
			return "", nil, errForbiddenScope
		}
		var unitType string
		err := a.pool.QueryRow(ctx, `select unit_type from public.units where id = $1 and is_active`, requestedUnitID).Scan(&unitType)
		if err != nil {
			return "", nil, errForbiddenScope
		}
		switch {
		case principal.IsAdmin:
			peopleMode = true
		case principal.PositionRole == "office_head" && unitType == "functional":
			peopleMode = true
		default:
			return "", nil, errForbiddenScope
		}
		targetUnit = requestedUnitID
		targetUnitType = unitType
	}

	if peopleMode {
		rows, err := a.pool.Query(ctx, `
			select u.id::text, u.nip, u.name, u.position_role, un.name,
			       coalesce(mus.original_deduction_rate, 0), coalesce(mus.effective_deduction_rate, 0),
			       coalesce(mus.deduction_days, 0), coalesce(mus.work_days, 0), coalesce(mus.overtime_days, 0),
			       coalesce(mus.leave_days, 0), coalesce(mus.off_days, 0)
			from public.users u
			join public.units un on un.id = u.unit_id
			left join public.monthly_user_summary mus on mus.user_id = u.id and mus.period_id = $1
			where u.is_active and u.deleted_at is null
			  and (
			    u.unit_id = $2
			    or ($3 = 'division' and un.parent_id = $2)
			    or ($3 = 'office' and (un.parent_id = $2 or un.id = $2 or un.parent_id in (select id from public.units where parent_id = $2)))
			  )
			order by u.position_role <> 'section_head', u.name`, periodID, targetUnit, targetUnitType)
		if err != nil {
			return "", nil, err
		}
		defer rows.Close()
		items := []map[string]any{}
		for rows.Next() {
			var id, nip, name, role, unitName string
			var original, effective float64
			var deductionDays, workDays, overtimeDays, leaveDays, offDays int
			if err := rows.Scan(&id, &nip, &name, &role, &unitName, &original, &effective, &deductionDays, &workDays, &overtimeDays, &leaveDays, &offDays); err != nil {
				return "", nil, err
			}
			items = append(items, map[string]any{"id": id, "nip": nip, "name": name, "role": role, "unit_name": unitName, "original": original, "effective": effective, "deduction_days": deductionDays, "work_days": workDays, "overtime_days": overtimeDays, "leave_days": leaveDays, "off_days": offDays})
		}
		return "people", items, rows.Err()
	}

	var unitFilter, memberFilter string
	args := []any{periodID}
	switch {
	case principal.PositionRole == "division_head":
		unitFilter = "un.parent_id = $2"
		memberFilter = "u.unit_id = un.id"
		args = append(args, principal.UnitID)
	case principal.PositionRole == "office_head" || principal.IsAdmin:
		unitFilter = "(un.parent_id = (select id from public.units where unit_type = 'office' limit 1) or un.unit_type = 'functional')"
		memberFilter = `(u.unit_id = un.id or exists (
			select 1 from public.units member_unit
			where member_unit.id = u.unit_id and member_unit.parent_id = un.id
		))`
	default:
		return "", nil, errForbiddenScope
	}
	query := fmt.Sprintf(`
		select un.id::text, un.name, un.unit_type, count(u.id),
		       coalesce(sum(mus.original_deduction_rate), 0), coalesce(sum(mus.effective_deduction_rate), 0),
		       coalesce(avg(coalesce(mus.effective_deduction_rate, 0)), 0), coalesce(sum(mus.deduction_days), 0)
		from public.units un
		left join public.users u on u.is_active and u.deleted_at is null and %s
		left join public.monthly_user_summary mus on mus.user_id = u.id and mus.period_id = $1
		where un.is_active and %s
		group by un.id, un.name, un.unit_type, un.sort_order
		order by un.sort_order, un.name`, memberFilter, unitFilter)
	rows, err := a.pool.Query(ctx, query, args...)
	if err != nil {
		return "", nil, err
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		var id, name, unitType string
		var members, deductionDays int
		var original, effective, average float64
		if err := rows.Scan(&id, &name, &unitType, &members, &original, &effective, &average, &deductionDays); err != nil {
			return "", nil, err
		}
		items = append(items, map[string]any{"unit_id": id, "unit_name": name, "unit_type": unitType, "members": members, "original": original, "effective": effective, "average": average, "deduction_days": deductionDays, "detail_allowed": principal.IsAdmin || (principal.PositionRole == "office_head" && unitType == "functional")})
	}
	return "aggregate", items, rows.Err()
}

func (a *App) monitoringTotals(request *http.Request, principal auth.Principal, periodID string) (map[string]any, error) {
	filter := "true"
	args := []any{periodID}
	switch {
	case principal.IsAdmin || principal.PositionRole == "office_head":
	case principal.PositionRole == "division_head":
		filter = "(u.unit_id = $2 or un.parent_id = $2)"
		args = append(args, principal.UnitID)
	case principal.PositionRole == "section_head":
		filter = "u.unit_id = $2"
		args = append(args, principal.UnitID)
	default:
		return nil, errForbiddenScope
	}
	query := fmt.Sprintf(`
		select count(u.id),
		       coalesce(sum(mus.original_deduction_rate), 0),
		       coalesce(sum(mus.effective_deduction_rate), 0),
		       coalesce(avg(coalesce(mus.effective_deduction_rate, 0)), 0),
		       coalesce(sum(mus.deduction_days), 0)
		from public.users u
		join public.units un on un.id = u.unit_id
		left join public.monthly_user_summary mus on mus.user_id = u.id and mus.period_id = $1
		where u.is_active and u.deleted_at is null and %s`, filter)
	var members, deductionDays int
	var original, effective, average float64
	if err := a.pool.QueryRow(request.Context(), query, args...).Scan(&members, &original, &effective, &average, &deductionDays); err != nil {
		return nil, err
	}
	return map[string]any{
		"members": members, "original": original, "effective": effective,
		"average": average, "deduction_days": deductionDays,
	}, nil
}

type warningUser struct {
	ID, Name, NIP, UnitID, UnitName, UnitType, ParentID, ParentName string
	Rates                                                           map[string]float64
}

func (a *App) warnings(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
	periods, err := a.recentPeriods(request, 12)
	if err != nil {
		a.internalError(response, "warning periods", err)
		return
	}
	if len(periods) == 0 {
		writeJSON(response, http.StatusOK, map[string]any{"individual": []any{}, "aggregate": []any{}})
		return
	}
	parameters, err := a.warningParameters(request)
	if err != nil {
		a.internalError(response, "warning parameters", err)
		return
	}
	users, err := a.warningUsers(request, principal, periods)
	if err != nil {
		a.internalError(response, "warning users", err)
		return
	}
	if !principal.IsAdmin {
		filtered := users[:0]
		for _, user := range users {
			if user.ID != principal.ID {
				filtered = append(filtered, user)
			}
		}
		users = filtered
	}
	current := periods[0].ID
	lookback := int(parameters["individual_anomaly_lookback_months"])
	habitPeriods := int(parameters["bad_habit_consecutive_periods"])
	priorMax := parameters["individual_anomaly_prior_max_rate"]
	individual := []map[string]any{}
	for _, user := range users {
		currentRate := user.Rates[current]
		priorTotal := 0.0
		for index := 1; index < len(periods) && index <= lookback; index++ {
			priorTotal += user.Rates[periods[index].ID]
		}
		if currentRate > 0 && priorTotal <= priorMax {
			individual = append(individual, map[string]any{"type": "individual_anomaly", "severity": "warning", "user_id": user.ID, "name": user.Name, "nip": user.NIP, "unit": user.UnitName, "current_rate": currentRate, "message": fmt.Sprintf("Potongan muncul setelah %d periode tanpa potongan berarti.", lookback)})
		}
		consecutive := habitPeriods > 0 && len(periods) >= habitPeriods
		for index := 0; consecutive && index < habitPeriods; index++ {
			consecutive = user.Rates[periods[index].ID] > 0
		}
		if consecutive {
			individual = append(individual, map[string]any{"type": "bad_habit", "severity": "danger", "user_id": user.ID, "name": user.Name, "nip": user.NIP, "unit": user.UnitName, "current_rate": currentRate, "message": fmt.Sprintf("Memiliki potongan %d periode berturut-turut.", habitPeriods)})
		}
	}

	type groupStats struct {
		ID, Name string
		Members  map[string]struct{}
		Totals   map[string]float64
	}
	groups := map[string]*groupStats{}
	for _, user := range users {
		groupID, groupName := warningGroup(principal, user)
		group, ok := groups[groupID]
		if !ok {
			group = &groupStats{ID: groupID, Name: groupName, Members: map[string]struct{}{}, Totals: map[string]float64{}}
			groups[groupID] = group
		}
		group.Members[user.ID] = struct{}{}
		for periodID, rate := range user.Rates {
			group.Totals[periodID] += rate
		}
	}
	spikeLookback := int(parameters["aggregate_spike_lookback_months"])
	spikeMultiplier := parameters["aggregate_spike_multiplier"]
	spikeDelta := parameters["aggregate_spike_min_delta"]
	averageThreshold := parameters["aggregate_average_threshold"]
	aggregate := []map[string]any{}
	for _, group := range groups {
		members := float64(max(1, len(group.Members)))
		currentAverage := group.Totals[current] / members
		if currentAverage > averageThreshold {
			aggregate = append(aggregate, map[string]any{"type": "average_threshold", "severity": "danger", "unit_id": group.ID, "unit": group.Name, "current_average": currentAverage, "members": len(group.Members), "message": "Rata-rata potongan unit melewati ambang parameter."})
		}
		priorAverage := 0.0
		priorCount := 0
		for index := 1; index < len(periods) && index <= spikeLookback; index++ {
			priorAverage += group.Totals[periods[index].ID] / members
			priorCount++
		}
		if priorCount > 0 {
			priorAverage /= float64(priorCount)
			if currentAverage-priorAverage >= spikeDelta && currentAverage > priorAverage*spikeMultiplier {
				aggregate = append(aggregate, map[string]any{"type": "aggregate_spike", "severity": "warning", "unit_id": group.ID, "unit": group.Name, "current_average": currentAverage, "baseline_average": priorAverage, "members": len(group.Members), "message": "Potongan unit melonjak dibanding rerata periode acuan."})
			}
		}
	}
	sort.Slice(individual, func(i, j int) bool { return individual[i]["name"].(string) < individual[j]["name"].(string) })
	sort.Slice(aggregate, func(i, j int) bool { return aggregate[i]["unit"].(string) < aggregate[j]["unit"].(string) })
	writeJSON(response, http.StatusOK, map[string]any{"period": periods[0], "individual": individual, "aggregate": aggregate, "parameters": parameters})
}

func (a *App) warningParameters(request *http.Request) (map[string]float64, error) {
	defaults := map[string]float64{"individual_anomaly_lookback_months": 6, "individual_anomaly_prior_max_rate": 0, "bad_habit_consecutive_periods": 3, "aggregate_spike_lookback_months": 6, "aggregate_spike_multiplier": 2, "aggregate_spike_min_delta": .005, "aggregate_average_threshold": .005}
	rows, err := a.pool.Query(request.Context(), `select key, value_json::text from public.parameters where category = 'warning'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var key, raw string
		if err := rows.Scan(&key, &raw); err != nil {
			return nil, err
		}
		var value float64
		if err := json.Unmarshal([]byte(raw), &value); err == nil {
			defaults[key] = value
		}
	}
	return defaults, rows.Err()
}

func (a *App) warningUsers(request *http.Request, principal auth.Principal, periods []periodInfo) ([]warningUser, error) {
	periodIDs := make([]string, len(periods))
	for index := range periods {
		periodIDs[index] = periods[index].ID
	}
	filter := "true"
	args := []any{periodIDs}
	if !principal.IsAdmin && principal.PositionRole == "section_head" {
		filter = "u.unit_id = $2"
		args = append(args, principal.UnitID)
	} else if !principal.IsAdmin && principal.PositionRole == "division_head" {
		filter = "(u.unit_id = $2 or un.parent_id = $2)"
		args = append(args, principal.UnitID)
	} else if !principal.IsAdmin && principal.PositionRole != "office_head" {
		return nil, errForbiddenScope
	}
	query := fmt.Sprintf(`
		select u.id::text, u.name, u.nip, un.id::text, un.name, un.unit_type,
		       coalesce(un.parent_id::text, ''), coalesce(parent.name, ''),
		       mus.period_id::text, coalesce(mus.effective_deduction_rate, 0)
		from public.users u
		join public.units un on un.id = u.unit_id
		left join public.units parent on parent.id = un.parent_id
		left join public.monthly_user_summary mus on mus.user_id = u.id and mus.period_id = any($1::uuid[])
		where u.is_active and u.deleted_at is null and %s
		order by u.name`, filter)
	rows, err := a.pool.Query(request.Context(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	byID := map[string]*warningUser{}
	order := []string{}
	for rows.Next() {
		var user warningUser
		var periodID *string
		var rate float64
		if err := rows.Scan(&user.ID, &user.Name, &user.NIP, &user.UnitID, &user.UnitName, &user.UnitType, &user.ParentID, &user.ParentName, &periodID, &rate); err != nil {
			return nil, err
		}
		existing, ok := byID[user.ID]
		if !ok {
			user.Rates = map[string]float64{}
			existing = &user
			byID[user.ID] = existing
			order = append(order, user.ID)
		}
		if periodID != nil {
			existing.Rates[*periodID] = rate
		}
	}
	result := make([]warningUser, 0, len(order))
	for _, id := range order {
		result = append(result, *byID[id])
	}
	return result, rows.Err()
}

func warningGroup(principal auth.Principal, user warningUser) (string, string) {
	if principal.PositionRole == "section_head" && !principal.IsAdmin {
		return user.UnitID, user.UnitName
	}
	if principal.PositionRole == "division_head" && !principal.IsAdmin {
		if user.ParentID != "" {
			return user.UnitID, user.UnitName
		}
		return user.UnitID, user.UnitName
	}
	if user.UnitType == "section" && user.ParentID != "" {
		return user.ParentID, user.ParentName
	}
	return user.UnitID, user.UnitName
}

func (a *App) loadDeductionRows(request *http.Request, userID, periodID string) ([]map[string]any, error) {
	rows, err := a.pool.Query(request.Context(), `
		select id, work_date, coalesce(to_char(check_in, 'HH24:MI'), ''), coalesce(to_char(check_out, 'HH24:MI'), ''),
		       coalesce(late_code, ''), coalesce(early_leave_code, ''),
		       coalesce(shift_code, ''), coalesce(attendance_status, ''), coalesce(leave_type, ''),
		       coalesce(assignment_type, ''), coalesce(notes, ''), deduction_rate,
		       effective_deduction_rate, deduction_components, coalesce(supervisor_status, ''), coalesce(admin_status, '')
		from public.effective_attendance
		where user_id = $1 and period_id = $2 and deduction_rate > 0
		order by work_date`, userID, periodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		var id int64
		var date time.Time
		var checkIn, checkOut string
		var late, early, shift, status, leave, assignment, notes, supervisorStatus, adminStatus string
		var original, effective float64
		var components []byte
		if err := rows.Scan(&id, &date, &checkIn, &checkOut, &late, &early, &shift, &status, &leave, &assignment, &notes, &original, &effective, &components, &supervisorStatus, &adminStatus); err != nil {
			return nil, err
		}
		items = append(items, map[string]any{"id": id, "date": date.Format("2006-01-02"), "check_in": checkIn, "check_out": checkOut, "late": late, "early_leave": early, "shift": shift, "status": status, "leave": leave, "assignment": assignment, "notes": notes, "original": original, "effective": effective, "components": json.RawMessage(components), "supervisor_status": supervisorStatus, "admin_status": adminStatus})
	}
	return items, rows.Err()
}

func (a *App) currentPeriod(request *http.Request) (periodInfo, error) {
	var result periodInfo
	var start, end time.Time
	err := a.pool.QueryRow(request.Context(), `
		select id::text, label, period_start, period_end
		from public.reporting_periods where published_batch_id is not null
		order by period_end desc, updated_at desc limit 1`).Scan(&result.ID, &result.Label, &start, &end)
	result.Start, result.End = start.Format("2006-01-02"), end.Format("2006-01-02")
	return result, err
}

func (a *App) recentPeriods(request *http.Request, limit int) ([]periodInfo, error) {
	rows, err := a.pool.Query(request.Context(), `
		select id::text, label, period_start, period_end
		from public.reporting_periods where published_batch_id is not null
		order by period_end desc limit $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []periodInfo{}
	for rows.Next() {
		var item periodInfo
		var start, end time.Time
		if err := rows.Scan(&item.ID, &item.Label, &start, &end); err != nil {
			return nil, err
		}
		item.Start, item.End = start.Format("2006-01-02"), end.Format("2006-01-02")
		result = append(result, item)
	}
	return result, rows.Err()
}

func monthRange(fromValue, toValue string) (time.Time, time.Time, error) {
	now := time.Now().UTC()
	if toValue == "" {
		toValue = now.Format("2006-01")
	}
	toMonth, err := time.Parse("2006-01", toValue)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("bulan akhir tidak valid")
	}
	if fromValue == "" {
		fromValue = toMonth.AddDate(0, -11, 0).Format("2006-01")
	}
	fromMonth, err := time.Parse("2006-01", fromValue)
	if err != nil || fromMonth.After(toMonth) || toMonth.Sub(fromMonth) > 62*31*24*time.Hour {
		return time.Time{}, time.Time{}, errors.New("rentang bulan tidak valid atau lebih dari 60 bulan")
	}
	toEnd := toMonth.AddDate(0, 1, -1)
	return fromMonth, toEnd, nil
}

func validUUID(value string) bool {
	parts := strings.Split(value, "-")
	if len(parts) != 5 || len(parts[0]) != 8 || len(parts[1]) != 4 || len(parts[2]) != 4 || len(parts[3]) != 4 || len(parts[4]) != 12 {
		return false
	}
	for _, part := range parts {
		if _, err := strconv.ParseUint(part, 16, 64); err != nil {
			return false
		}
	}
	return true
}

func (a *App) internalError(response http.ResponseWriter, operation string, err error) {
	a.log.Error(operation, "error", err)
	writeError(response, http.StatusInternalServerError, "Data belum dapat diproses.", "internal_error")
}
