package httpapi

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/bcpriok/pantas/internal/auth"
	"github.com/jackc/pgx/v5"
)

type treasuryRecapRow struct {
	Number        int
	Name          string
	NIP           string
	Unit          string
	EffectiveRate float64
}

type treasuryRecapData struct {
	Period periodInfo
	Rows   []treasuryRecapRow
}

func (a *App) treasuryRecap(response http.ResponseWriter, request *http.Request, _ auth.Principal) {
	data, err := a.loadTreasuryRecap(request)
	if errors.Is(err, pgx.ErrNoRows) {
		writeJSON(response, http.StatusOK, map[string]any{"current_period": nil, "items": []any{}, "total": 0, "with_deduction": 0})
		return
	}
	if err != nil {
		a.internalError(response, "treasury recap", err)
		return
	}
	items := make([]map[string]any, 0, len(data.Rows))
	withDeduction := 0
	totalRate := 0.0
	for _, row := range data.Rows {
		if row.EffectiveRate > 0 {
			withDeduction++
		}
		totalRate += row.EffectiveRate
		items = append(items, map[string]any{
			"number": row.Number, "name": row.Name, "nip": row.NIP,
			"unit": row.Unit, "effective": row.EffectiveRate,
		})
	}
	writeJSON(response, http.StatusOK, map[string]any{
		"current_period": data.Period, "items": items, "total": len(items),
		"with_deduction": withDeduction, "total_effective": totalRate,
	})
}

func (a *App) treasuryRecapXLSX(response http.ResponseWriter, request *http.Request, _ auth.Principal) {
	data, err := a.loadTreasuryRecap(request)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(response, http.StatusNotFound, "Belum ada periode berjalan yang dipublikasikan.", "period_not_found")
		return
	}
	if err != nil {
		a.internalError(response, "treasury export data", err)
		return
	}
	file, err := buildTreasuryRecapXLSX(data, time.Now().UTC())
	if err != nil {
		a.internalError(response, "treasury export xlsx", err)
		return
	}
	filename := "Rekap_Potongan_Efektif_" + strings.ReplaceAll(data.Period.Label, " ", "_") + ".xlsx"
	response.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	response.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": filename}))
	response.Header().Set("Cache-Control", "no-store")
	response.Header().Set("Content-Length", fmt.Sprint(len(file)))
	response.WriteHeader(http.StatusOK)
	_, _ = response.Write(file)
}

func (a *App) loadTreasuryRecap(request *http.Request) (treasuryRecapData, error) {
	period, err := a.currentPeriod(request)
	if err != nil {
		return treasuryRecapData{}, err
	}
	rows, err := a.pool.Query(request.Context(), `
		select u.name, u.nip, un.name, coalesce(mus.effective_deduction_rate, 0)
		from public.users u
		join public.units un on un.id = u.unit_id
		left join public.monthly_user_summary mus on mus.user_id = u.id and mus.period_id = $1
		where u.is_active and u.deleted_at is null
		order by un.sort_order, un.name, u.name`, period.ID)
	if err != nil {
		return treasuryRecapData{}, err
	}
	defer rows.Close()
	result := treasuryRecapData{Period: period, Rows: []treasuryRecapRow{}}
	for rows.Next() {
		var row treasuryRecapRow
		row.Number = len(result.Rows) + 1
		if err := rows.Scan(&row.Name, &row.NIP, &row.Unit, &row.EffectiveRate); err != nil {
			return treasuryRecapData{}, err
		}
		result.Rows = append(result.Rows, row)
	}
	return result, rows.Err()
}
