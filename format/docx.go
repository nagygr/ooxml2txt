package format

/*
 * Using code snippets from:
 * https://github.com/nguyenthenguyen/docx/
 * (c) Nguyen The Nguyen
 */

import (
	"github.com/nagygr/ooxml2txt/archive"
)

type Docx struct {
	zipReader archive.ZipData
	Text      string
	Links     []string
	Footnotes []string
	Headers   []string
	Footers   []string
}

func MakeDocx(path string) (*Docx, error) {
	reader, err := archive.MakeZipFile(path)

	if err != nil {
		return nil, err
	}

	textXml, err := readXml(reader, "word/document.xml")
	if err != nil {
		return nil, err
	}

	text, err := textFromXml(textXml)
	if err != nil {
		return nil, err
	}

	linksXml, err := readXml(reader, "word/_rels/document.xml.rels")
	if err != nil {
		return nil, err
	}

	links, err := linksFromXml(linksXml)
	if err != nil {
		return nil, err
	}

	headersXmls, err := readXmls(reader, "header")
	var headers []string

	if err == nil {
		headers, err = textListFromXmls(headersXmls)
	}

	if err != nil {
		headers = []string{}
	}

	footersXmls, err := readXmls(reader, "footer")
	var footers []string

	if err == nil {
		footers, err = textListFromXmls(footersXmls)
	}

	if err != nil {
		footers = []string{}
	}

	footnotesXml, err := readXml(reader, "word/footnotes.xml")
	var footnotes []string

	if err == nil {
		footnotes, err = textListFromXml(footnotesXml)
	}

	if err != nil {
		footnotes = []string{}
	}

	return &Docx {
		zipReader: reader,
		Text: text,
		Links: links,
		Footnotes: footnotes,
		Headers: headers,
		Footers: footers}, nil
}
