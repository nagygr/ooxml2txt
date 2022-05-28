package format

import (
	"strings"
	"testing"
)

func TestMakingXlsxGoodPath(t *testing.T) {
	path := "../../test_data/example.xlsx"
	_, err := MakeXlsx(path)

	if err != nil {
		t.Errorf("Expected to open %s successfully", path)
	}
}

func TestMakingXlsxBadPath(t *testing.T) {
	path := "../../test_data/wrong_example.xlsx"
	_, err := MakeXlsx(path)

	if err == nil {
		t.Errorf("Expected to fail to open %s successfully", path)
	}
}

func TestReadingXlsxText(t *testing.T) {
	path := "../../test_data/example.xlsx"
	xls, err := MakeXlsx(path)

	if err != nil {
		t.Errorf("Expected to open %s successfully", path)
	}
	if len(xls.Text) != 3 {
		t.Errorf("Expected to have 3 strings, has: %d", len(xls.Text))
	}

	text := "Odense"
	if !strings.HasPrefix(xls.Text[0], text) {
		t.Errorf("Expected the text to start with: \"%s\"", text)
	}

	text = "Copenhagen"
	if !strings.HasSuffix(strings.TrimSpace(xls.Text[2]), text) {
		t.Errorf("Expected the text to end with: \"%s\"", text)
	}
}

func TestReadingXlsxTextFromUrl(t *testing.T) {
	url := "https://github.com/nagygr/ooxml2txt/raw/main/test_data/example.xlsx"
	xls, err := MakeXlsxFromUrl(url)

	if err != nil {
		t.Errorf("Expected to open %s successfully", url)
	}
	if len(xls.Text) != 3 {
		t.Errorf("Expected to have 3 strings, has: %d", len(xls.Text))
	}

	text := "Odense"
	if !strings.HasPrefix(xls.Text[0], text) {
		t.Errorf("Expected the text to start with: \"%s\"", text)
	}

	text = "Copenhagen"
	if !strings.HasSuffix(strings.TrimSpace(xls.Text[2]), text) {
		t.Errorf("Expected the text to end with: \"%s\"", text)
	}
}
