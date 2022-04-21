package format

/*
 * Using code snippets from:
 * https://github.com/nguyenthenguyen/docx/
 * (c) Nguyen The Nguyen
 */

import (
	"encoding/xml"
	"github.com/nagygr/ooxml2txt/archive"
	"io"
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

	content, err := readXml(reader, "word/document.xml")
	if err != nil {
		return nil, err
	}

	links, err := readXml(reader, "word/_rels/document.xml.rels")
	if err != nil {
		return nil, err
	}

	/*
	 * Error not handled: it is ok to not find headers/footers
	 */
	headers, _ := readXmls(reader, "header")
	footers, _ := readXmls(reader, "footer")

	/*
	 * Error not handled: it is ok to not find footnotes
	 */
	footnotes, _ := readXml(reader, "word/footnotes.xml")

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

func (d *Docx) Headers() (headers []string) {
	for _, headerXml := range d.headers {
		var (
			reader = strings.NewReader(headerXml)
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
						headers = append(headers, string(t))
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
	}

	return
}

func (d *Docx) Footers() (footers []string) {
	for _, footerXml := range d.footers {
		var (
			reader = strings.NewReader(footerXml)
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
						footers = append(footers, string(t))
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
	}

	return
}
