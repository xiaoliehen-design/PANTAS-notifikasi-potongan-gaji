package httpapi

import (
	"archive/zip"
	"bytes"
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"
)

func buildTreasuryRecapXLSX(data treasuryRecapData, generatedAt time.Time) ([]byte, error) {
	var output bytes.Buffer
	archive := zip.NewWriter(&output)
	parts := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>
  <Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>
  <Override PartName="/xl/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml"/>
</Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>
</Relationships>`,
		"xl/workbook.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <sheets><sheet name="Rekap Potongan" sheetId="1" r:id="rId1"/></sheets>
</workbook>`,
		"xl/_rels/workbook.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/>
  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>
</Relationships>`,
		"xl/styles.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">
  <fonts count="3">
    <font><sz val="11"/><name val="Aptos"/><family val="2"/></font>
    <font><b/><sz val="16"/><color rgb="FF0B3B66"/><name val="Aptos Display"/><family val="2"/></font>
    <font><b/><sz val="11"/><color rgb="FFFFFFFF"/><name val="Aptos"/><family val="2"/></font>
  </fonts>
  <fills count="3">
    <fill><patternFill patternType="none"/></fill>
    <fill><patternFill patternType="gray125"/></fill>
    <fill><patternFill patternType="solid"><fgColor rgb="FF1479BF"/><bgColor indexed="64"/></patternFill></fill>
  </fills>
  <borders count="2">
    <border><left/><right/><top/><bottom/><diagonal/></border>
    <border><left style="thin"><color rgb="FFDCE5EC"/></left><right style="thin"><color rgb="FFDCE5EC"/></right><top style="thin"><color rgb="FFDCE5EC"/></top><bottom style="thin"><color rgb="FFDCE5EC"/></bottom><diagonal/></border>
  </borders>
  <cellStyleXfs count="1"><xf numFmtId="0" fontId="0" fillId="0" borderId="0"/></cellStyleXfs>
  <cellXfs count="5">
    <xf numFmtId="0" fontId="0" fillId="0" borderId="0" xfId="0"/>
    <xf numFmtId="0" fontId="1" fillId="0" borderId="0" xfId="0" applyFont="1"/>
    <xf numFmtId="0" fontId="2" fillId="2" borderId="1" xfId="0" applyFont="1" applyFill="1" applyBorder="1"><alignment horizontal="center" vertical="center"/></xf>
    <xf numFmtId="0" fontId="0" fillId="0" borderId="1" xfId="0" applyBorder="1"><alignment vertical="center"/></xf>
    <xf numFmtId="10" fontId="0" fillId="0" borderId="1" xfId="0" applyNumberFormat="1" applyBorder="1"><alignment horizontal="right" vertical="center"/></xf>
  </cellXfs>
  <cellStyles count="1"><cellStyle name="Normal" xfId="0" builtinId="0"/></cellStyles>
</styleSheet>`,
	}
	for _, name := range []string{"[Content_Types].xml", "_rels/.rels", "xl/workbook.xml", "xl/_rels/workbook.xml.rels", "xl/styles.xml"} {
		if err := writeXLSXPart(archive, name, parts[name]); err != nil {
			return nil, err
		}
	}
	if err := writeXLSXPart(archive, "xl/worksheets/sheet1.xml", treasuryWorksheetXML(data, generatedAt)); err != nil {
		return nil, err
	}
	if err := archive.Close(); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func treasuryWorksheetXML(data treasuryRecapData, generatedAt time.Time) string {
	lastRow := len(data.Rows) + 5
	if lastRow < 5 {
		lastRow = 5
	}
	var xml strings.Builder
	xml.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	xml.WriteString(`<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">`)
	fmt.Fprintf(&xml, `<dimension ref="A1:E%d"/>`, lastRow)
	xml.WriteString(`<sheetViews><sheetView workbookViewId="0"><pane ySplit="5" topLeftCell="A6" activePane="bottomLeft" state="frozen"/></sheetView></sheetViews>`)
	xml.WriteString(`<sheetFormatPr defaultRowHeight="18"/>`)
	xml.WriteString(`<cols><col min="1" max="1" width="7" customWidth="1"/><col min="2" max="2" width="34" customWidth="1"/><col min="3" max="3" width="23" customWidth="1"/><col min="4" max="4" width="38" customWidth="1"/><col min="5" max="5" width="22" customWidth="1"/></cols>`)
	xml.WriteString(`<sheetData>`)
	xml.WriteString(`<row r="1" ht="25" customHeight="1">` + inlineStringCell("A1", "REKAPITULASI POTONGAN EFEKTIF", 1) + `</row>`)
	xml.WriteString(`<row r="2">` + inlineStringCell("A2", "Periode: "+data.Period.Label+" ("+data.Period.Start+" s.d. "+data.Period.End+")", 0) + `</row>`)
	xml.WriteString(`<row r="3">` + inlineStringCell("A3", "Dibuat: "+generatedAt.Format("02-01-2006 15:04 UTC"), 0) + `</row>`)
	xml.WriteString(`<row r="5" ht="24" customHeight="1">`)
	for index, label := range []string{"No.", "Nama", "NIP", "Unit", "Potongan Efektif"} {
		xml.WriteString(inlineStringCell(cellReference(index+1, 5), label, 2))
	}
	xml.WriteString(`</row>`)
	for index, item := range data.Rows {
		row := index + 6
		fmt.Fprintf(&xml, `<row r="%d">`, row)
		xml.WriteString(numberCell(cellReference(1, row), float64(item.Number), 3))
		xml.WriteString(inlineStringCell(cellReference(2, row), item.Name, 3))
		xml.WriteString(inlineStringCell(cellReference(3, row), item.NIP, 3))
		xml.WriteString(inlineStringCell(cellReference(4, row), item.Unit, 3))
		xml.WriteString(numberCell(cellReference(5, row), item.EffectiveRate, 4))
		xml.WriteString(`</row>`)
	}
	xml.WriteString(`</sheetData>`)
	xml.WriteString(`<mergeCells count="1"><mergeCell ref="A1:E1"/></mergeCells>`)
	fmt.Fprintf(&xml, `<autoFilter ref="A5:E%d"/>`, lastRow)
	xml.WriteString(`<pageMargins left="0.4" right="0.4" top="0.6" bottom="0.6" header="0.3" footer="0.3"/>`)
	xml.WriteString(`<pageSetup orientation="landscape" fitToWidth="1" fitToHeight="0"/>`)
	xml.WriteString(`</worksheet>`)
	return xml.String()
}

func writeXLSXPart(archive *zip.Writer, name, content string) error {
	header := &zip.FileHeader{Name: name, Method: zip.Deflate}
	header.SetModTime(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
	part, err := archive.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = part.Write([]byte(content))
	return err
}

func inlineStringCell(reference, value string, style int) string {
	return fmt.Sprintf(`<c r="%s" s="%d" t="inlineStr"><is><t xml:space="preserve">%s</t></is></c>`, reference, style, html.EscapeString(value))
}

func numberCell(reference string, value float64, style int) string {
	return fmt.Sprintf(`<c r="%s" s="%d"><v>%s</v></c>`, reference, style, strconv.FormatFloat(value, 'f', -1, 64))
}

func cellReference(column, row int) string {
	name := ""
	for column > 0 {
		column--
		name = string(rune('A'+column%26)) + name
		column /= 26
	}
	return name + strconv.Itoa(row)
}
