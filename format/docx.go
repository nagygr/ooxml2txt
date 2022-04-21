package format

/*
 * Using code snippets from:
 * https://github.com/nguyenthenguyen/docx/
 * (c) Nguyen The Nguyen
 */

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"github.com/nagygr/ooxml2txt/archive"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

type Docx struct {
	zipReader archive.ZipData
	content   string
	links     string
	footnotes string
	headers   map[string]string
	footers   map[string]string
}

func MakeDocx(path string) (*Docx, error) {
	reader, err := archive.MakeZipFile(path)

	if err != nil {
		return nil, err
	}

	content, err := readText(reader.Files())
	if err != nil {
		return nil, err
	}

	links, err := readLinks(reader.Files())
	if err != nil {
		return nil, err
	}

	/*
	 * Error not handled: it is ok to not find headers/footers
	 */
	headers, footers, _ := readHeaderFooter(reader.Files())

	/*
	 * Error not handled: it is ok to not find footnotes
	 */
	footnotes, _ := readFootnotes(reader.Files())

	return &Docx {
		zipReader: reader,
		content: content,
		links: links,
		footnotes: footnotes,
		headers: headers,
		footers: footers}, nil
}

func (d *Docx) Text() string {
	var (
		contents = strings.NewReader(d.content)
		decoder = xml.NewDecoder(contents)
		text strings.Builder
		inText bool = false
	)

	for {
		token, err := decoder.Token()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error while parsing xml file: %s", err.Error())
		}

		switch t := token.(type) {
			case xml.CharData:
				if inText {
					text.WriteString(string(t))
				}
			case xml.StartElement:
				if t.Name.Local == "t" {
					inText = true
				}
			case xml.EndElement:
				if t.Name.Local == "t" {
					inText = false
				}
			default:
		}
	}

	return text.String()
}

func (d *Docx) Links() (links []string) {
	const (
		tagName = "Relationship"
		typeName = "Type"
		targetName = "Target"
		urlType = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink"
	)

	var (
		contents = strings.NewReader(d.links)
		decoder = xml.NewDecoder(contents)
		urlFound bool
	)

	for {
		token, err := decoder.Token()
		urlFound = false

		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error while parsing xml file: %s", err.Error())
		}

		switch t := token.(type) {
			case xml.StartElement:
				if t.Name.Local == tagName {
					var url string

					for _, a := range t.Attr {
						if a.Name.Local == typeName && a.Value == urlType {
							urlFound = true
						} else if a.Name.Local == targetName {
							url = a.Value
						}
					}

					if urlFound {
						links = append(links, url)
					}
				}
			default:
		}
	}

	return
}

func (d *Docx) Footnotes() (footnotes []string) {
	var (
		reader = strings.NewReader(d.footnotes)
		decoder = xml.NewDecoder(reader)
		inText bool = false
	)

	for {
		token, err := decoder.Token()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error while parsing xml file: %s", err.Error())
		}

		switch t := token.(type) {
			case xml.CharData:
				if inText {
					footnotes = append(footnotes, string(t))
				}
			case xml.StartElement:
				if t.Name.Local == "t" {
					inText = true
				}
			case xml.EndElement:
				if t.Name.Local == "t" {
					inText = false
				}
			default:
		}
	}

	return
}

func readFootnotes(files []*zip.File) (string, error) {
	f, err := retrieveFootnoteDoc(files)

	if err != nil {
		return "", err
	}

	footnotes, err := buildFootnotes(f)
	if err != nil {
		return "", err
	}

	return footnotes, nil
}

func buildFootnotes(footnotesFile *zip.File) (string, error) {
	documentReader, err := footnotesFile.Open()
	if err != nil {
		return "", err
	}

	footnotesText, err := wordDocToString(documentReader)
	if err != nil {
		return "", err
	}

	return footnotesText, nil
}

func readHeaderFooter(files []*zip.File) (headerText map[string]string, footerText map[string]string, err error) {
	h, f, err := retrieveHeaderFooterDoc(files)

	if err != nil {
		return map[string]string{}, map[string]string{}, err
	}

	headerText, err = buildHeaderFooter(h)
	if err != nil {
		return map[string]string{}, map[string]string{}, err
	}

	footerText, err = buildHeaderFooter(f)
	if err != nil {
		return map[string]string{}, map[string]string{}, err
	}

	return headerText, footerText, err
}

func buildHeaderFooter(headerFooter []*zip.File) (map[string]string, error) {
	headerFooterText := make(map[string]string)

	for _, element := range headerFooter {
		documentReader, err := element.Open()
		if err != nil {
			return map[string]string{}, err
		}

		text, err := wordDocToString(documentReader)
		if err != nil {
			return map[string]string{}, err
		}

		headerFooterText[element.Name] = text
	}

	return headerFooterText, nil
}

func readText(files []*zip.File) (text string, err error) {
	var documentFile *zip.File
	documentFile, err = retrieveWordDoc(files)
	if err != nil {
		return text, err
	}

	var documentReader io.ReadCloser
	documentReader, err = documentFile.Open()
	if err != nil {
		return text, err
	}

	text, err = wordDocToString(documentReader)
	return
}

func readLinks(files []*zip.File) (text string, err error) {
	var documentFile *zip.File
	documentFile, err = retrieveLinkDoc(files)
	if err != nil {
		return text, err
	}

	var documentReader io.ReadCloser
	documentReader, err = documentFile.Open()
	if err != nil {
		return text, err
	}

	text, err = wordDocToString(documentReader)
	return
}

func wordDocToString(reader io.Reader) (string, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func retrieveWordDoc(files []*zip.File) (file *zip.File, err error) {
	for _, f := range files {
		if f.Name == "word/document.xml" {
			file = f
		}
	}
	if file == nil {
		err = errors.New("document.xml file not found")
	}
	return
}

func retrieveLinkDoc(files []*zip.File) (file *zip.File, err error) {
	for _, f := range files {
		if f.Name == "word/_rels/document.xml.rels" {
			file = f
		}
	}
	if file == nil {
		err = errors.New("document.xml.rels file not found")
	}
	return
}

func retrieveHeaderFooterDoc(files []*zip.File) (headers []*zip.File, footers []*zip.File, err error) {
	for _, f := range files {

		if strings.Contains(f.Name, "header") {
			headers = append(headers, f)
		}
		if strings.Contains(f.Name, "footer") {
			footers = append(footers, f)
		}
	}
	if len(headers) == 0 && len(footers) == 0 {
		err = errors.New("headers[1-3].xml file not found and footers[1-3].xml file not found")
	}
	return
}

func retrieveFootnoteDoc(files []*zip.File) (footnotes *zip.File, err error) {
	for _, f := range files {
		if strings.Contains(f.Name, "footnotes") {
			footnotes =  f
			break
		}
	}
	if footnotes == nil {
		err = errors.New("footnotes.xml file not found")
	}
	return
}
