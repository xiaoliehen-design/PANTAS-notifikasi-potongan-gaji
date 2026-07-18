package importer

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const expectedSheetName = "DETAIL WFH WFO"

var expectedHeaders = []string{
	"Tanggal", "Nama", "NIP", "Bidang", "Locus Penempatan", "Jam Masuk", "Jam Pulang",
	"TL", "PSW", "Shift", "Status", "Cuti", "Penugasan", "Konfirmasi", "Keterangan",
}

type RawRecord struct {
	SourceRow           int
	WorkDate            time.Time
	Name                string
	NIP                 string
	SourceDivision      string
	SourcePlacement     string
	CheckIn             *time.Time
	CheckOut            *time.Time
	LateCode            string
	EarlyLeaveCode      string
	ShiftCode           string
	AttendanceStatus    string
	LeaveType           string
	AssignmentType      string
	SourceConfirmation  string
	Notes               string
	DeductionRate       float64
	DeductionComponents []DeductionComponent
}

type DeductionComponent struct {
	SourceField string  `json:"source_field"`
	Code        string  `json:"code"`
	Label       string  `json:"label"`
	Rate        float64 `json:"rate"`
}

type ParseResult struct {
	SheetName       string
	PeriodLabel     string
	PeriodStart     time.Time
	PeriodEnd       time.Time
	PrintedAt       string
	IntegrityStatus string
	FileSHA256      string
	FileSize        int64
	Records         []RawRecord
	BlankRows       int
	Warnings        []string
}

type ValidationError struct {
	Problems []string
}

func (e *ValidationError) Error() string {
	return "format Excel tidak valid: " + strings.Join(e.Problems, "; ")
}

func ParseXLSX(data []byte) (ParseResult, error) {
	if len(data) < 4 || !bytes.Equal(data[:4], []byte{'P', 'K', 3, 4}) {
		return ParseResult{}, &ValidationError{Problems: []string{"file bukan dokumen XLSX"}}
	}
	container, integrity, err := openContainer(data)
	if err != nil {
		return ParseResult{}, &ValidationError{Problems: []string{err.Error()}}
	}
	sheetName, sheetPath, err := discoverSheet(container)
	if err != nil {
		return ParseResult{}, &ValidationError{Problems: []string{err.Error()}}
	}
	if sheetName != expectedSheetName {
		return ParseResult{}, &ValidationError{Problems: []string{fmt.Sprintf("nama sheet harus %q, ditemukan %q", expectedSheetName, sheetName)}}
	}
	shared, err := readSharedStrings(container)
	if err != nil {
		return ParseResult{}, &ValidationError{Problems: []string{"sharedStrings.xml tidak dapat dibaca: " + err.Error()}}
	}
	sheetReader, err := container.Open(sheetPath)
	if err != nil {
		return ParseResult{}, &ValidationError{Problems: []string{"sheet data tidak ditemukan: " + err.Error()}}
	}
	defer sheetReader.Close()
	parsed, err := parseWorksheet(sheetReader, shared)
	if err != nil {
		return ParseResult{}, err
	}
	hash := sha256.Sum256(data)
	parsed.SheetName = sheetName
	parsed.IntegrityStatus = integrity
	parsed.FileSHA256 = hex.EncodeToString(hash[:])
	parsed.FileSize = int64(len(data))
	if integrity == "recovered_partial_container" {
		parsed.Warnings = append(parsed.Warnings, "Bagian akhir container XLSX tidak lengkap; seluruh sheet utama berhasil dipulihkan dan divalidasi.")
	}
	return parsed, nil
}

type xlsxContainer interface {
	Open(name string) (io.ReadCloser, error)
	Names() []string
}

type standardContainer struct {
	files map[string]*zip.File
}

func (c *standardContainer) Open(name string) (io.ReadCloser, error) {
	file, ok := c.files[strings.TrimPrefix(name, "/")]
	if !ok {
		return nil, fmt.Errorf("entry %q tidak ditemukan", name)
	}
	return file.Open()
}

func (c *standardContainer) Names() []string {
	names := make([]string, 0, len(c.files))
	for name := range c.files {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

type recoveredEntry struct {
	data   []byte
	method uint16
}

type recoveredContainer struct {
	entries map[string]recoveredEntry
}

func (c *recoveredContainer) Open(name string) (io.ReadCloser, error) {
	entry, ok := c.entries[strings.TrimPrefix(name, "/")]
	if !ok {
		return nil, fmt.Errorf("entry %q tidak ditemukan", name)
	}
	switch entry.method {
	case zip.Store:
		return io.NopCloser(bytes.NewReader(entry.data)), nil
	case zip.Deflate:
		return flate.NewReader(bytes.NewReader(entry.data)), nil
	default:
		return nil, fmt.Errorf("metode kompresi %d tidak didukung", entry.method)
	}
}

func (c *recoveredContainer) Names() []string {
	names := make([]string, 0, len(c.entries))
	for name := range c.entries {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func openContainer(data []byte) (xlsxContainer, string, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err == nil {
		files := make(map[string]*zip.File, len(reader.File))
		for _, file := range reader.File {
			files[file.Name] = file
		}
		return &standardContainer{files: files}, "valid", nil
	}
	recovered, recoverErr := recoverLocalEntries(data)
	if recoverErr != nil {
		return nil, "", fmt.Errorf("container XLSX rusak dan tidak dapat dipulihkan: %w", recoverErr)
	}
	return recovered, "recovered_partial_container", nil
}

func recoverLocalEntries(data []byte) (*recoveredContainer, error) {
	const headerSize = 30
	entries := map[string]recoveredEntry{}
	offset := 0
	for offset+headerSize <= len(data) {
		if binary.LittleEndian.Uint32(data[offset:offset+4]) != 0x04034b50 {
			next := bytes.Index(data[offset+1:], []byte{'P', 'K', 3, 4})
			if next < 0 {
				break
			}
			offset += next + 1
			continue
		}
		flags := binary.LittleEndian.Uint16(data[offset+6 : offset+8])
		method := binary.LittleEndian.Uint16(data[offset+8 : offset+10])
		compressedSize := int(binary.LittleEndian.Uint32(data[offset+18 : offset+22]))
		nameLength := int(binary.LittleEndian.Uint16(data[offset+26 : offset+28]))
		extraLength := int(binary.LittleEndian.Uint16(data[offset+28 : offset+30]))
		if flags&0x1 != 0 {
			return nil, errors.New("file terenkripsi tidak didukung")
		}
		if flags&0x8 != 0 {
			return nil, errors.New("pemulihan entry dengan data descriptor tidak didukung")
		}
		nameStart := offset + headerSize
		nameEnd := nameStart + nameLength
		dataStart := nameEnd + extraLength
		dataEnd := dataStart + compressedSize
		if nameEnd > len(data) || dataStart > len(data) {
			break
		}
		if dataEnd > len(data) {
			break
		}
		name := string(data[nameStart:nameEnd])
		entries[name] = recoveredEntry{data: data[dataStart:dataEnd], method: method}
		offset = dataEnd
	}
	for _, required := range []string{"xl/workbook.xml", "xl/_rels/workbook.xml.rels", "xl/worksheets/sheet1.xml"} {
		if _, ok := entries[required]; !ok {
			return nil, fmt.Errorf("entry wajib %q tidak lengkap", required)
		}
	}
	return &recoveredContainer{entries: entries}, nil
}

func discoverSheet(container xlsxContainer) (string, string, error) {
	type sheet struct {
		Name string `xml:"name,attr"`
		RID  string `xml:"http://schemas.openxmlformats.org/officeDocument/2006/relationships id,attr"`
	}
	var workbook struct {
		Sheets []sheet `xml:"sheets>sheet"`
	}
	reader, err := container.Open("xl/workbook.xml")
	if err != nil {
		return "", "", err
	}
	if err := xml.NewDecoder(reader).Decode(&workbook); err != nil {
		reader.Close()
		return "", "", err
	}
	reader.Close()
	if len(workbook.Sheets) != 1 {
		return "", "", fmt.Errorf("workbook harus berisi tepat satu sheet; ditemukan %d", len(workbook.Sheets))
	}

	var relationships struct {
		Items []struct {
			ID     string `xml:"Id,attr"`
			Target string `xml:"Target,attr"`
			Type   string `xml:"Type,attr"`
		} `xml:"Relationship"`
	}
	relsReader, err := container.Open("xl/_rels/workbook.xml.rels")
	if err != nil {
		return "", "", err
	}
	if err := xml.NewDecoder(relsReader).Decode(&relationships); err != nil {
		relsReader.Close()
		return "", "", err
	}
	relsReader.Close()
	for _, relationship := range relationships.Items {
		if relationship.ID == workbook.Sheets[0].RID && strings.Contains(relationship.Type, "/worksheet") {
			target := strings.TrimPrefix(relationship.Target, "/")
			if !strings.HasPrefix(target, "xl/") {
				target = path.Clean("xl/" + target)
			}
			return workbook.Sheets[0].Name, target, nil
		}
	}
	return "", "", errors.New("relasi sheet utama tidak ditemukan")
}

func readSharedStrings(container xlsxContainer) ([]string, error) {
	reader, err := container.Open("xl/sharedStrings.xml")
	if err != nil {
		// Shared strings bersifat opsional bila semua cell memakai inline string.
		return nil, nil
	}
	defer reader.Close()
	decoder := xml.NewDecoder(reader)
	var result []string
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "si" {
			continue
		}
		var item struct {
			Texts []string `xml:"t"`
			Runs  []struct {
				Text string `xml:"t"`
			} `xml:"r"`
		}
		if err := decoder.DecodeElement(&item, &start); err != nil {
			return nil, err
		}
		var builder strings.Builder
		for _, text := range item.Texts {
			builder.WriteString(text)
		}
		for _, run := range item.Runs {
			builder.WriteString(run.Text)
		}
		result = append(result, builder.String())
	}
	return result, nil
}

type worksheetCell struct {
	Reference string `xml:"r,attr"`
	Type      string `xml:"t,attr"`
	Value     string `xml:"v"`
	Inline    string `xml:"is>t"`
}

type worksheetRow struct {
	Number int             `xml:"r,attr"`
	Cells  []worksheetCell `xml:"c"`
}

func parseWorksheet(reader io.Reader, shared []string) (ParseResult, error) {
	decoder := xml.NewDecoder(reader)
	metadata := map[string]string{}
	headers := make([]string, len(expectedHeaders))
	result := ParseResult{}
	seen := map[string]int{}
	problems := []string{}

	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return ParseResult{}, &ValidationError{Problems: []string{"XML sheet tidak dapat dibaca: " + err.Error()}}
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "row" {
			continue
		}
		var row worksheetRow
		if err := decoder.DecodeElement(&row, &start); err != nil {
			return ParseResult{}, &ValidationError{Problems: []string{"baris sheet tidak dapat dibaca: " + err.Error()}}
		}
		values := map[string]string{}
		for _, cell := range row.Cells {
			column := cellColumn(cell.Reference)
			value, valueErr := resolvedCellValue(cell, shared)
			if valueErr != nil {
				appendProblem(&problems, fmt.Sprintf("baris %d kolom %s: %v", row.Number, column, valueErr))
				continue
			}
			values[column] = strings.TrimSpace(value)
		}

		switch {
		case row.Number <= 3:
			for column, value := range values {
				metadata[fmt.Sprintf("%s%d", column, row.Number)] = value
			}
		case row.Number == 4:
			for index := range headers {
				headers[index] = values[columnName(index+1)]
			}
		default:
			if allBlank(values) {
				result.BlankRows++
				continue
			}
			record, recordProblems := parseDataRow(row.Number, values)
			for _, problem := range recordProblems {
				appendProblem(&problems, problem)
			}
			if len(recordProblems) == 0 {
				key := record.NIP + "|" + record.WorkDate.Format("2006-01-02")
				if prior, exists := seen[key]; exists {
					appendProblem(&problems, fmt.Sprintf("NIP %s tanggal %s muncul dua kali (baris %d dan %d)", record.NIP, record.WorkDate.Format("02-01-2006"), prior, row.Number))
				} else {
					seen[key] = row.Number
					result.Records = append(result.Records, record)
				}
			}
		}
	}

	for index, expected := range expectedHeaders {
		if headers[index] != expected {
			appendProblem(&problems, fmt.Sprintf("header %s4 harus %q, ditemukan %q", columnName(index+1), expected, headers[index]))
		}
	}
	periodStart, periodEnd, err := parseIndonesianPeriod(metadata["B2"])
	if err != nil {
		appendProblem(&problems, err.Error())
	} else {
		result.PeriodStart = periodStart
		result.PeriodEnd = periodEnd
		result.PeriodLabel = monthLabel(periodEnd)
		for _, record := range result.Records {
			if record.WorkDate.Before(periodStart) || record.WorkDate.After(periodEnd) {
				appendProblem(&problems, fmt.Sprintf("tanggal baris %d berada di luar periode %s", record.SourceRow, metadata["B2"]))
			}
		}
	}
	result.PrintedAt = metadata["B3"]
	if len(result.Records) == 0 {
		appendProblem(&problems, "tidak ada baris data aktif")
	}
	if len(problems) > 0 {
		return ParseResult{}, &ValidationError{Problems: problems}
	}
	return result, nil
}

func parseDataRow(rowNumber int, values map[string]string) (RawRecord, []string) {
	problems := []string{}
	record := RawRecord{
		SourceRow:          rowNumber,
		Name:               strings.TrimSpace(values["B"]),
		NIP:                strings.TrimSpace(values["C"]),
		SourceDivision:     strings.TrimSpace(values["D"]),
		SourcePlacement:    strings.TrimSpace(values["E"]),
		LateCode:           strings.TrimSpace(values["H"]),
		EarlyLeaveCode:     strings.TrimSpace(values["I"]),
		ShiftCode:          strings.TrimSpace(values["J"]),
		AttendanceStatus:   strings.TrimSpace(values["K"]),
		LeaveType:          strings.TrimSpace(values["L"]),
		AssignmentType:     strings.TrimSpace(values["M"]),
		SourceConfirmation: strings.TrimSpace(values["N"]),
		Notes:              strings.TrimSpace(values["O"]),
	}
	if record.Name == "" {
		problems = append(problems, fmt.Sprintf("baris %d: Nama kosong", rowNumber))
	}
	if !validNIP(record.NIP) {
		problems = append(problems, fmt.Sprintf("baris %d: NIP harus tepat 18 digit", rowNumber))
	}
	date, err := excelSerialDate(values["A"])
	if err != nil {
		problems = append(problems, fmt.Sprintf("baris %d: Tanggal tidak valid (%v)", rowNumber, err))
	} else {
		record.WorkDate = date
	}
	if values["F"] != "" {
		checkIn, err := excelSerialTime(values["F"], date)
		if err != nil {
			problems = append(problems, fmt.Sprintf("baris %d: Jam Masuk tidak valid", rowNumber))
		} else {
			record.CheckIn = &checkIn
		}
	}
	if values["G"] != "" {
		checkOut, err := excelSerialTime(values["G"], date)
		if err != nil {
			problems = append(problems, fmt.Sprintf("baris %d: Jam Pulang tidak valid", rowNumber))
		} else {
			record.CheckOut = &checkOut
		}
	}
	return record, problems
}

func resolvedCellValue(cell worksheetCell, shared []string) (string, error) {
	switch cell.Type {
	case "s":
		index, err := strconv.Atoi(strings.TrimSpace(cell.Value))
		if err != nil || index < 0 || index >= len(shared) {
			return "", errors.New("indeks shared string tidak valid")
		}
		return shared[index], nil
	case "inlineStr":
		return cell.Inline, nil
	default:
		return cell.Value, nil
	}
}

func cellColumn(reference string) string {
	index := 0
	for index < len(reference) && reference[index] >= 'A' && reference[index] <= 'Z' {
		index++
	}
	return reference[:index]
}

func columnName(index int) string {
	var name []byte
	for index > 0 {
		index--
		name = append([]byte{byte('A' + index%26)}, name...)
		index /= 26
	}
	return string(name)
}

func allBlank(values map[string]string) bool {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return false
		}
	}
	return true
}

func appendProblem(problems *[]string, problem string) {
	if len(*problems) < 100 {
		*problems = append(*problems, problem)
	}
}

func excelSerialDate(value string) (time.Time, error) {
	serial, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil || serial < 1 || serial > 2958465 {
		return time.Time{}, errors.New("serial date di luar batas")
	}
	days := int(serial)
	return time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC).AddDate(0, 0, days), nil
}

func excelSerialTime(value string, date time.Time) (time.Time, error) {
	serial, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil || serial < 0 {
		return time.Time{}, errors.New("serial time tidak valid")
	}
	seconds := int64((serial-float64(int64(serial)))*86400 + 0.5)
	seconds %= 86400
	return time.Date(date.Year(), date.Month(), date.Day(), int(seconds/3600), int((seconds%3600)/60), int(seconds%60), 0, time.UTC), nil
}

var periodPattern = regexp.MustCompile(`(?i)^\s*(\d{1,2})\s+([[:alpha:]]+)\s+(\d{4})\s+s\.d\.\s+(\d{1,2})\s+([[:alpha:]]+)\s+(\d{4})\s*$`)

func parseIndonesianPeriod(value string) (time.Time, time.Time, error) {
	match := periodPattern.FindStringSubmatch(value)
	if match == nil {
		return time.Time{}, time.Time{}, fmt.Errorf("cell B2 harus menggunakan format periode seperti '16 Juni 2026 s.d. 15 Juli 2026'")
	}
	months := map[string]time.Month{
		"januari": 1, "februari": 2, "maret": 3, "april": 4, "mei": 5, "juni": 6,
		"juli": 7, "agustus": 8, "september": 9, "oktober": 10, "november": 11, "desember": 12,
	}
	startMonth, startOK := months[strings.ToLower(match[2])]
	endMonth, endOK := months[strings.ToLower(match[5])]
	if !startOK || !endOK {
		return time.Time{}, time.Time{}, errors.New("nama bulan pada cell B2 tidak dikenali")
	}
	startDay, _ := strconv.Atoi(match[1])
	startYear, _ := strconv.Atoi(match[3])
	endDay, _ := strconv.Atoi(match[4])
	endYear, _ := strconv.Atoi(match[6])
	start := time.Date(startYear, startMonth, startDay, 0, 0, 0, 0, time.UTC)
	end := time.Date(endYear, endMonth, endDay, 0, 0, 0, 0, time.UTC)
	if start.Day() != startDay || end.Day() != endDay || end.Before(start) || end.Sub(start) > 62*24*time.Hour {
		return time.Time{}, time.Time{}, errors.New("rentang tanggal pada cell B2 tidak valid")
	}
	return start, end, nil
}

func monthLabel(date time.Time) string {
	months := []string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	return fmt.Sprintf("%s %d", months[int(date.Month())], date.Year())
}

func validNIP(nip string) bool {
	if len(nip) != 18 {
		return false
	}
	for _, char := range nip {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}
