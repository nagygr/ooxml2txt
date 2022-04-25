package format

import (
	"github.com/nagygr/ooxml2txt/internal/archive"
	. "github.com/nagygr/ooxml2txt/internal/format"
)

// Xlsx handles xlsx documents. The Text member is a list of strings where each
// element corresponds to a string value in the document. Only the strings are
// collected, numbers, formulas and binary data is ignored. Only unique strings
// are collected, i.e. if a piece of text appears multiple times in the
// document, it will only show up once in the list.
type Xlsx struct {
	zipReader archive.ZipData
	Text      []string
}

// MakeXlsx creates a Xlsx from the path to a spreadsheet document. The
// returned instance contains the valid contents of the document if there was
// no error while processing it (which is then reported in the returned error
// value).
func MakeXlsx(path string) (*Xlsx, error) {
	reader, err := archive.MakeZipFile(path)

	if err != nil {
		return nil, err
	}

	sharedStringsXml, err := ReadXml(reader, "xl/sharedStrings.xml")

	if err != nil {
		return nil, err
	}

	var sharedStrings []string

	sharedStrings, err = XlsxSharedStringsFromXml(sharedStringsXml)

	if err != nil {
		return nil, err
	}

	return &Xlsx{
		zipReader: reader,
		Text: sharedStrings}, nil

}
