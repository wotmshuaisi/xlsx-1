package xlsx

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"
	"time"
)

type CellIndexTestCase struct {
	x        uint64
	y        uint64
	expected string
}

func TestCellIndex(t *testing.T) {

	// tests := []CellIndexTestCase{
	// 	{0, 0, "A1"},
	// 	{2, 2, "C3"},
	// 	{26, 45, "AA46"},
	// 	{2600, 100000, "CVA100001"},
	// }

	// for _, c := range tests {
	// 	cellX, cellY := CellIndex(c.x, c.y)
	// 	s := fmt.Sprintf("%s%d", cellX, cellY)
	// 	if s != c.expected {
	// 		t.Errorf("expected %s, got %s", c.expected, s)
	// 	}
	// }
}

type OADateTestCase struct {
	datetime time.Time
	expected string
}

func TestOADate(t *testing.T) {

	tests := []OADateTestCase{
		{time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), "25569"},
		{time.Date(1970, 1, 1, 12, 20, 0, 0, time.UTC), "25569.513889"},
		{time.Date(2014, 12, 20, 0, 0, 0, 0, time.UTC), "41993"},
	}

	for _, d := range tests {
		s := OADate(d.datetime)
		if s != d.expected {
			t.Errorf("expected %s, got %s", d.expected, s)
		}
	}
}

func TestTemplates(t *testing.T) {

	var b bytes.Buffer
	var err error
	var s Sheet

	err = TemplateContentTypes.Execute(&b, nil)
	if err != nil {
		t.Errorf("template TemplateContentTypes failed to Execute returning error %s", err.Error())
	}

	err = TemplateRelationships.Execute(&b, nil)
	if err != nil {
		t.Errorf("template TemplateRelationships failed to Execute returning error %s", err.Error())
	}

	err = TemplateApp.Execute(&b, nil)
	if err != nil {
		t.Errorf("template TemplateApp failed to Execute returning error %s", err.Error())
	}

	err = TemplateCore.Execute(&b, s.DocumentInfo)
	if err != nil {
		t.Errorf("template TemplateCore failed to Execute returning error %s", err.Error())
	}

	err = TemplateWorkbook.Execute(&b, nil)
	if err != nil {
		t.Errorf("template TemplateWorkbook failed to Execute returning error %s", err.Error())
	}

	err = TemplateWorkbookRelationships.Execute(&b, nil)
	if err != nil {
		t.Errorf("template TemplateWorkbookRelationships failed to Execute returning error %s", err.Error())
	}

	err = TemplateStyles.Execute(&b, nil)
	if err != nil {
		t.Errorf("template TemplateStyles failed to Execute returning error %s", err.Error())
	}

	err = TemplateStringLookups.Execute(&b, []string{})
	if err != nil {
		t.Errorf("template TemplateStringLookups failed to Execute returning error %s", err.Error())
	}

	sheet := struct {
		Cols  []Column
		Rows  []string
		Start string
		End   string
	}{
		Cols:  []Column{},
		Rows:  []string{},
		Start: "A1",
		End:   "C3",
	}

	err = TemplateSheetStart.Execute(&b, sheet)
	if err != nil {
		t.Errorf("template TemplateSheetStart failed to Execute returning error %s", err.Error())
	}

	for i := range sheet.Rows {
		rb := &bytes.Buffer{}
		rowString := fmt.Sprintf(`<row r="%d">%s</row>`, uint64(i), rb.String())
		_, err = io.WriteString(&b, rowString)
	}
}

func TestSheetWriter(t *testing.T) {
	outputfile, err := os.Create("test.xlsx")

	w := bufio.NewWriter(outputfile)
	ww := NewWorkbookWriter(w)

	c := []Column{
		{Name: "Col1", Width: 10},
		{Name: "Col2", Width: 10},
	}

	sh := NewSheetWithColumns(c)
	sh.Title = "MySheet"

	sw, err := ww.NewSheetWriter(&sh)

	for i := 0; i < 10; i++ {

		r := sh.NewRow()

		r.Cells = append(r.Cells, Cell{
			Type:  CellTypeNumber,
			Value: strconv.Itoa(i + 1),
		})
		r.Cells = append(r.Cells, Cell{
			Type:  CellTypeNumber,
			Value: strconv.Itoa(i + 1),
		})

		err = sw.WriteRows([]Row{r})
		if err != nil {
			t.Fatal(err)
		}
	}

	err = ww.Close()
	defer w.Flush()
	if err != nil {
		t.Fatal(err)
	}
}
