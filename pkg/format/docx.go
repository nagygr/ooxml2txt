// Package format contains format handling types. Each ooxml format that can be
// handled by the library has its own type plus there is a generic XML-handling
// type that forms the basis for the operations on ooxml files.
package format

/*
 * Using code snippets from:
 * https://github.com/nguyenthenguyen/docx/
 * (c) Nguyen The Nguyen
 */

import (
	"github.com/nagygr/ooxml2txt/internal/archive"
	. "github.com/nagygr/ooxml2txt/internal/format"
)

// Docx handles docx documents. Its fields contain the textual information
// corresponding to the different elements of the document: Text contains the
// document text, Links is a list of links that appear in the document (the
// text part contains references to the links), Footnotes contains the list of
// footnotes, and Headers and Footers are also lists and contain the headers
// and footers of the document.
type Docx struct {
	zipReader archive.ZipData
	Text      string
	Links     []string
	Footnotes []string
	Headers   []string
	Footers   []string
}

// MakeDocx creates a Docx that parses the document given by its path. The
// returned instance contains the valid contents of the document if there was
// no error while processing it. If there was an error, it is reported in the
// returned error value).
func MakeDocx(path string) (*Docx, error) {
	reader, err := archive.MakeZipFile(path)

	if err != nil {
		return nil, err
	}

	textXml, err := ReadXml(reader, "word/document.xml")
	if err != nil {
		return nil, err
	}

	text, err := TextFromXml(textXml)
	if err != nil {
		return nil, err
	}

	linksXml, err := ReadXml(reader, "word/_rels/document.xml.rels")
	if err != nil {
		return nil, err
	}

	links, err := LinksFromXml(linksXml)
	if err != nil {
		return nil, err
	}

	headersXmls, err := ReadXmls(reader, "header")
	var headers []string

	if err == nil {
		headers, err = TextListFromXmls(headersXmls)
	}

	if err != nil {
		headers = []string{}
	}

	footersXmls, err := ReadXmls(reader, "footer")
	var footers []string

	if err == nil {
		footers, err = TextListFromXmls(footersXmls)
	}

	if err != nil {
		footers = []string{}
	}

	footnotesXml, err := ReadXml(reader, "word/footnotes.xml")
	var footnotes []string

	if err == nil {
		footnotes, err = TextListFromXml(footnotesXml)
	}

	if err != nil {
		footnotes = []string{}
	}

	return &Docx{
		zipReader: reader,
		Text:      text,
		Links:     links,
		Footnotes: footnotes,
		Headers:   headers,
		Footers:   footers}, nil
}
