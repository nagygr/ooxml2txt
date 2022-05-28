package format

import (
	"strings"
	"testing"
)

func TestMakingPptxGoodPath(t *testing.T) {
	path := "../../test_data/example.pptx"
	_, err := MakePptx(path)

	if err != nil {
		t.Errorf("Expected to open %s successfully", path)
	}
}

func TestMakingPptxBadPath(t *testing.T) {
	path := "../../test_data/wrong_example.pptx"
	_, err := MakePptx(path)

	if err == nil {
		t.Errorf("Expected to fail to open %s successfully", path)
	}
}

func TestReadingPptxText(t *testing.T) {
	path := "../../test_data/example.pptx"
	ppt, err := MakePptx(path)

	if err != nil {
		t.Errorf("Expected to open %s successfully", path)
	}

	if len(ppt.Text) != 2 {
		t.Errorf("Expected to have 2 slides, has: %d", len(ppt.Text))
	}

	text := "Aalborg"
	if !strings.HasPrefix(ppt.Text[0], text) {
		t.Errorf("Expected the text to start with: \"%s\" (%s)", text, ppt.Text[0])
	}

	text = "Limfjordsbroen"
	if !strings.HasSuffix(strings.TrimSpace(ppt.Text[1]), text) {
		t.Errorf("Expected the text to end with: \"%s\" (%s)", text, ppt.Text[1])
	}
}

func TestReadingPptxTextFromUrl(t *testing.T) {
	url := "https://github.com/nagygr/ooxml2txt/raw/main/test_data/example.pptx"
	ppt, err := MakePptxFromUrl(url)

	if err != nil {
		t.Fatalf("Expected to open %s successfully", url)
	}

	if len(ppt.Text) != 2 {
		t.Errorf("Expected to have 2 slides, has: %d", len(ppt.Text))
	}

	text := "Aalborg"
	if !strings.HasPrefix(ppt.Text[0], text) {
		t.Errorf("Expected the text to start with: \"%s\" (%s)", text, ppt.Text[0])
	}

	text = "Limfjordsbroen"
	if !strings.HasSuffix(strings.TrimSpace(ppt.Text[1]), text) {
		t.Errorf("Expected the text to end with: \"%s\" (%s)", text, ppt.Text[1])
	}
}
