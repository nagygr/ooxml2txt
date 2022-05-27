package format

import (
	"strings"
	"testing"
)

func TestMakingDocxGoodPath(t *testing.T) {
	path := "../../test_data/example.docx"
	_, err := MakeDocx(path)

	if err != nil {
		t.Errorf("Expected to open %s successfully", path)
	}
}

func TestMakingDocxBadPath(t *testing.T) {
	path := "../../test_data/wrong_example.docx"
	_, err := MakeDocx(path)

	if err == nil {
		t.Errorf("Expected to fail to open %s successfully", path)
	}
}

func TestReadingDocxText(t *testing.T) {
	path := "../../test_data/example.docx"
	doc, err := MakeDocx(path)

	if err != nil {
		t.Errorf("Expected to open %s successfully", path)
	}

	text := "The largest city"
	if !strings.HasPrefix(doc.Text, text) {
		t.Errorf("Expected the text to start with: \"%s\"", text)
	}

	text = "tourism in Jutland."
	if !strings.HasSuffix(doc.Text, text) {
		t.Errorf("Expected the text to end with: \"%s\"", text)
	}

	if len(doc.Links) != 7 {
		t.Errorf("Expected to have 7 links, has: %d", len(doc.Links))
	} else {
		linkFragment := "Denmark_Region"
		if !strings.Contains(doc.Links[0], linkFragment) {
			t.Errorf("First link (%s) expected to contain: %s", doc.Links[0], linkFragment)
		}
	}

	if len(doc.Footnotes) != 1 {
		t.Errorf("Expected to have one footnote, has %d", len(doc.Footnotes))
	} else {
		footnote := "This is a footnote."
		if doc.Footnotes[0] != footnote {
			t.Errorf("Expected the first footnote to be: %s, was: %s", footnote, doc.Footnotes[0])
		}
	}

	if len(doc.Headers) != 1 {
		t.Errorf("Expected to have one header, has %d", len(doc.Headers))
	} else {
		header := "Header text"
		if doc.Headers[0] != header {
			t.Errorf("Expected the first header to be: %s, was: %s", header, doc.Headers[0])
		}
	}

	if len(doc.Footers) != 1 {
		t.Errorf("Expected to have one footer, has %d", len(doc.Footers))
	} else {
		footer := "Footer text"
		if doc.Footers[0] != footer {
			t.Errorf("Expected the first footer to be: %s, was: %s", footer, doc.Footers[0])
		}
	}
}

func TestMissingDocumentXml(t *testing.T) {
	path := "../../test_data/broken_missing_document_xml.docx"
	_, err := MakeDocx(path)

	t.Logf("%s", err.Error())
	if err == nil {
		t.Errorf("Expected to get an error due to the missing document.xml.")
	}
}
