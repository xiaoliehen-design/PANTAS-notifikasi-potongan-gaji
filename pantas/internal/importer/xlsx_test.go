package importer

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
	"strings"
	"testing"
)

func TestParseXLSXValidContainer(t *testing.T) {
	result, err := ParseXLSX(testWorkbook(t, false, false))
	if err != nil {
		t.Fatalf("ParseXLSX() error = %v", err)
	}
	if result.IntegrityStatus != "valid" {
		t.Fatalf("integrity = %q, want valid", result.IntegrityStatus)
	}
	assertTestWorkbook(t, result)
}

func TestParseXLSXRecoversTruncatedContainer(t *testing.T) {
	result, err := ParseXLSX(testWorkbook(t, true, false))
	if err != nil {
		t.Fatalf("ParseXLSX() error = %v", err)
	}
	if result.IntegrityStatus != "recovered_partial_container" {
		t.Fatalf("integrity = %q, want recovered_partial_container", result.IntegrityStatus)
	}
	assertTestWorkbook(t, result)
}

func TestParseXLSXRejectsChangedHeader(t *testing.T) {
	_, err := ParseXLSX(testWorkbook(t, false, true))
	if err == nil || !strings.Contains(err.Error(), `header O4 harus "Keterangan"`) {
		t.Fatalf("ParseXLSX() error = %v, want changed-header validation", err)
	}
}

func TestSuppliedWorkbookWhenConfigured(t *testing.T) {
	path := os.Getenv("PANTAS_SAMPLE_XLSX")
	if path == "" {
		t.Skip("set PANTAS_SAMPLE_XLSX to run the supplied-workbook integration test")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	result, err := ParseXLSX(data)
	if err != nil {
		t.Fatalf("ParseXLSX(supplied workbook) error = %v", err)
	}
	if result.SheetName != expectedSheetName || len(result.Records) != 33690 {
		t.Fatalf("sheet = %q, rows = %d; want %q and 33690", result.SheetName, len(result.Records), expectedSheetName)
	}
	employees := map[string]struct{}{}
	for _, record := range result.Records {
		employees[record.NIP] = struct{}{}
	}
	if len(employees) != 1123 {
		t.Fatalf("employees = %d, want 1123", len(employees))
	}
}

func TestApplyRulesAddsIndependentComponents(t *testing.T) {
	record := RawRecord{LateCode: "TL1", EarlyLeaveCode: "PSW1"}
	rules := map[string]Rule{
		"late\x00TL1":         {SourceField: "late", Code: "TL1", Label: "Terlambat", Rate: .01},
		"early_leave\x00PSW1": {SourceField: "early_leave", Code: "PSW1", Label: "Pulang awal", Rate: .005},
	}
	applyRules(&record, rules)
	if record.DeductionRate != .015 || len(record.DeductionComponents) != 2 {
		t.Fatalf("rate = %v, components = %d", record.DeductionRate, len(record.DeductionComponents))
	}
}

func assertTestWorkbook(t *testing.T, result ParseResult) {
	t.Helper()
	if result.SheetName != expectedSheetName || result.PeriodLabel != "Juli 2026" || len(result.Records) != 1 {
		t.Fatalf("unexpected result: sheet=%q period=%q records=%d", result.SheetName, result.PeriodLabel, len(result.Records))
	}
	record := result.Records[0]
	if record.NIP != "199001012010011001" || record.WorkDate.Format("2006-01-02") != "2026-07-01" || record.LateCode != "TL1" {
		t.Fatalf("unexpected record: %+v", record)
	}
}

func testWorkbook(t *testing.T, recovered, changedHeader bool) []byte {
	t.Helper()
	workbook := `<?xml version="1.0" encoding="UTF-8"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheets><sheet name="DETAIL WFH WFO" sheetId="1" r:id="rId1"/></sheets></workbook>`
	rels := `<?xml version="1.0" encoding="UTF-8"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/></Relationships>`
	header := append([]string(nil), expectedHeaders...)
	if changedHeader {
		header[14] = "Catatan"
	}
	row4 := make([]string, 0, len(header))
	for index, value := range header {
		row4 = append(row4, inlineCell(columnName(index+1)+"4", value))
	}
	data := []string{"46204", "Siti Uji", "199001012010011001", "Bidang Uji", "Bidang Uji - Seksi Uji", "0.3333333333", "0.7083333333", "TL1", "", "P", "", "", "", "", ""}
	row5 := make([]string, 0, len(data))
	for index, value := range data {
		cell := columnName(index+1) + "5"
		if index == 0 || index == 5 || index == 6 {
			row5 = append(row5, fmt.Sprintf(`<c r="%s"><v>%s</v></c>`, cell, value))
		} else if value != "" {
			row5 = append(row5, inlineCell(cell, value))
		}
	}
	sheet := `<?xml version="1.0" encoding="UTF-8"?><worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>` +
		`<row r="2">` + inlineCell("B2", "16 Juni 2026 s.d. 15 Juli 2026") + `</row>` +
		`<row r="3">` + inlineCell("B3", "Kamis, 16 Juli 2026") + `</row>` +
		`<row r="4">` + strings.Join(row4, "") + `</row>` +
		`<row r="5">` + strings.Join(row5, "") + `</row>` +
		`</sheetData></worksheet>`
	entries := map[string][]byte{
		"xl/workbook.xml":            []byte(workbook),
		"xl/_rels/workbook.xml.rels": []byte(rels),
		"xl/worksheets/sheet1.xml":   []byte(sheet),
	}
	if recovered {
		return localEntriesWithoutCentralDirectory(t, entries)
	}
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	for _, name := range []string{"xl/workbook.xml", "xl/_rels/workbook.xml.rels", "xl/worksheets/sheet1.xml"} {
		file, err := writer.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := file.Write(entries[name]); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
}

func inlineCell(reference, value string) string {
	value = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(value)
	return fmt.Sprintf(`<c r="%s" t="inlineStr"><is><t>%s</t></is></c>`, reference, value)
}

func localEntriesWithoutCentralDirectory(t *testing.T, entries map[string][]byte) []byte {
	t.Helper()
	var buffer bytes.Buffer
	for _, name := range []string{"xl/workbook.xml", "xl/_rels/workbook.xml.rels", "xl/worksheets/sheet1.xml"} {
		data := entries[name]
		header := make([]byte, 30)
		binary.LittleEndian.PutUint32(header[0:4], 0x04034b50)
		binary.LittleEndian.PutUint16(header[4:6], 20)
		binary.LittleEndian.PutUint16(header[8:10], zip.Store)
		binary.LittleEndian.PutUint32(header[14:18], crc32.ChecksumIEEE(data))
		binary.LittleEndian.PutUint32(header[18:22], uint32(len(data)))
		binary.LittleEndian.PutUint32(header[22:26], uint32(len(data)))
		binary.LittleEndian.PutUint16(header[26:28], uint16(len(name)))
		if _, err := buffer.Write(header); err != nil {
			t.Fatal(err)
		}
		buffer.WriteString(name)
		buffer.Write(data)
	}
	return buffer.Bytes()
}
