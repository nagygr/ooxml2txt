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
	"fmt"
	"github.com/nagygr/ooxml2txt/archive"
	"io"
	"io/ioutil"
	"strings"
)

func readXmls(zipReader archive.ZipData, nameFragment string) (xmlTexts map[string]string, err error) {
	xmlFiles, err := zipReader.FilesByName(nameFragment);

	if err != nil {
		return map[string]string{}, err
	}

	xmlTexts, err = readTextFromXmls(xmlFiles)
	if err != nil {
		return map[string]string{}, err
	}

	return
}

func readTextFromXmls(xmlFiles []*zip.File) (map[string]string, error) {
	xmlText := make(map[string]string)

	for _, element := range xmlFiles {
		documentReader, err := element.Open()
		if err != nil {
			return map[string]string{}, err
		}

		text, err := xmlFileToString(documentReader)
		if err != nil {
			return map[string]string{}, err
		}

		xmlText[element.Name] = text
	}

	return xmlText, nil
}

func readXml(zipReader archive.ZipData, path string) (text string, err error) {
	var documentFile *zip.File
	documentFile, err = zipReader.FileByName(path)
	if err != nil {
		return text, err
	}

	var documentReader io.ReadCloser
	documentReader, err = documentFile.Open()
	if err != nil {
		return text, err
	}

	text, err = xmlFileToString(documentReader)
	return
}


func xmlFileToString(reader io.Reader) (string, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func textFromXml(xmlText string) (string, error) {
	var (
		contents = strings.NewReader(xmlText)
		decoder = xml.NewDecoder(contents)
		text strings.Builder
		inText bool = false
	)

	for {
		token, err := decoder.Token()

		if err == io.EOF {
			break
		} else if err != nil {
			return "",
				errors.New(fmt.Sprintf("Error while parsing xml file: %s", err.Error()))
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

	return text.String(), nil
}

func linksFromXml(xmlLinks string) (links []string, err error) {
	const (
		tagName = "Relationship"
		typeName = "Type"
		targetName = "Target"
		urlType = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink"
	)

	var (
		contents = strings.NewReader(xmlLinks)
		decoder = xml.NewDecoder(contents)
		urlFound bool
	)

	for {
		token, decErr := decoder.Token()
		urlFound = false

		if decErr == io.EOF {
			break
		} else if decErr != nil {
			err = errors.New(fmt.Sprintf("Error while parsing xml file: %s", decErr.Error()))
			return
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

func textListFromXml(textXml string) (textList []string, err error) {
	var (
		reader = strings.NewReader(textXml)
		decoder = xml.NewDecoder(reader)
		inText bool = false
	)

	for {
		token, decErr := decoder.Token()

		if decErr == io.EOF {
			break
		} else if decErr != nil {
			err = errors.New(fmt.Sprintf("Error while parsing xml file: %s", decErr.Error()))
		}

		switch t := token.(type) {
			case xml.CharData:
				if inText {
					textList = append(textList, string(t))
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

func textListFromXmls(textXmls map[string]string) (textList []string, err error) {
	for _, textXml := range textXmls {
		var (
			reader = strings.NewReader(textXml)
			decoder = xml.NewDecoder(reader)
			inText bool = false
		)

		for {
			token, decErr := decoder.Token()

			if decErr == io.EOF {
				break
			} else if decErr != nil {
				err = errors.New(fmt.Sprintf("Error while parsing xml file: %s", decErr.Error()))
			}

			switch t := token.(type) {
				case xml.CharData:
					if inText {
						textList = append(textList, string(t))
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
