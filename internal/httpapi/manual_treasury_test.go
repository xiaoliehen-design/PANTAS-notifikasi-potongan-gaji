package httpapi

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"strings"
	"testing"
	"time"
)

func TestValidateManualReason(t *testing.T) {
	for _, code := range []string{"input_error", "system_error"} {
		reason, err := validateManualReason(code, "")
		if err != nil || reason.Code != code || reason.Detail != nil {
			t.Fatalf("valid reason %q rejected: %#v, %v", code, reason, err)
		}
	}
	reason, err := validateManualReason("other", "Dokumen sumber salah versi")
	if err != nil || reason.Detail == nil || *reason.Detail != "Dokumen sumber salah versi" {
		t.Fatalf("valid other reason rejected: %#v, %v", reason, err)
	}
	for _, test := range []struct{ code, detail string }{{"", ""}, {"other", ""}, {"other", "x"}, {"unknown", "cukup panjang"}} {
		if _, err := validateManualReason(test.code, test.detail); err == nil {
			t.Fatalf("invalid reason accepted: %#v", test)
		}
	}
}

func TestManualAdjustedComponentsKeepsCodesAndReconcilesRate(t *testing.T) {
	raw := []byte(`[{"source_field":"late","code":"TL1","label":"Terlambat","rate":0.005},{"source_field":"early_leave","code":"PSW1","label":"Pulang cepat","rate":0.005}]`)
	result := manualAdjustedComponents(raw, 0.025, "Kesalahan sistem")
	var components []map[string]any
	if err := json.Unmarshal(result, &components); err != nil {
		t.Fatal(err)
	}
	if len(components) != 2 || components[0]["code"] != "TL1" || components[1]["code"] != "PSW1" {
		t.Fatalf("codes were not preserved: %#v", components)
	}
	if components[0]["rate"].(float64) != 0.025 || components[1]["rate"].(float64) != 0 {
		t.Fatalf("rates do not reconcile: %#v", components)
	}
}

func TestPeriodDeleteConfirmation(t *testing.T) {
	if actual := periodDeleteConfirmation("Juli 2026"); actual != "HAPUS DATA JULI 2026" {
		t.Fatalf("confirmation = %q", actual)
	}
}

func TestBuildTreasuryRecapXLSX(t *testing.T) {
	data := treasuryRecapData{
		Period: periodInfo{Label: "Juli 2026", Start: "2026-06-16", End: "2026-07-15"},
		Rows:   []treasuryRecapRow{{Number: 1, Name: "Pegawai Contoh", NIP: "199001012020011001", Unit: "Seksi Contoh", EffectiveRate: 0.025}},
	}
	content, err := buildTreasuryRecapXLSX(data, time.Date(2026, 7, 21, 10, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	reader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		t.Fatalf("invalid XLSX zip: %v", err)
	}
	parts := map[string]string{}
	for _, file := range reader.File {
		stream, err := file.Open()
		if err != nil {
			t.Fatal(err)
		}
		data, err := io.ReadAll(stream)
		stream.Close()
		if err != nil {
			t.Fatal(err)
		}
		parts[file.Name] = string(data)
	}
	for _, required := range []string{"[Content_Types].xml", "_rels/.rels", "xl/workbook.xml", "xl/_rels/workbook.xml.rels", "xl/styles.xml", "xl/worksheets/sheet1.xml"} {
		if _, ok := parts[required]; !ok {
			t.Errorf("missing XLSX part %s", required)
		}
	}
	sheet := parts["xl/worksheets/sheet1.xml"]
	if !strings.Contains(sheet, "Pegawai Contoh") || !strings.Contains(sheet, "199001012020011001") || !strings.Contains(sheet, `<c r="E6" s="4"><v>0.025</v></c>`) {
		t.Fatalf("worksheet is missing expected values: %s", sheet)
	}
	decoder := xml.NewDecoder(strings.NewReader(sheet))
	for {
		if _, err := decoder.Token(); err == io.EOF {
			break
		} else if err != nil {
			t.Fatalf("worksheet XML invalid: %v", err)
		}
	}
}
