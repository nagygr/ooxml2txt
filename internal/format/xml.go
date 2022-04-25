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
	"github.com/nagygr/ooxml2txt/internal/archive"
	"io"
	"io/ioutil"
	"strings"
)

func ReadXmls(zipReader archive.ZipData, nameFragment string) (xmlTexts []string, err error) {
	xmlFiles, err := zipReader.FilesByName(nameFragment);

	if err != nil {
		return []string{}, err
	}

	xmlTexts, err = ReadTextFromXmls(xmlFiles)
	if err != nil {
		return []string{}, err
	}

	return
}

func ReadTextFromXmls(xmlFiles []*zip.File) ([]string, error) {
	xmlText := []string{}

	for _, element := range xmlFiles {
		documentReader, err := element.Open()
		if err != nil {
			return []string{}, err
		}

		text, err := XmlFileToString(documentReader)
		if err != nil {
			return []string{}, err
		}

		xmlText = append(xmlText, text)
	}

	return xmlText, nil
}

func ReadXml(zipReader archive.ZipData, path string) (text string, err error) {
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

	text, err = XmlFileToString(documentReader)
	return
}


func XmlFileToString(reader io.Reader) (string, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func TextFromXml(xmlText string) (string, error) {
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

func LinksFromXml(xmlLinks string) (links []string, err error) {
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

func TextListFromXml(textXml string) (textList []string, err error) {
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

func TextListFromXmls(textXmls []string) (textList []string, err error) {
	for _, textXml := range textXmls {
		tmpList, errXml := TextListFromXml(textXml)

		if errXml != nil {
			err = errXml
			return
		}

		slideText := strings.Join(tmpList, " ")
		textList = append(textList, slideText)
	}

	return
}

func XlsxSharedStringsFromXml(sharedStringsXml string) (sharedStrings []string, err error) {
	var (
		reader = strings.NewReader(sharedStringsXml)
		decoder = xml.NewDecoder(reader)
		inSi bool = false
		inT bool = false
		currentString strings.Builder
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
				if inT {
					currentString.WriteString(string(t))
				}
			case xml.StartElement:
				if t.Name.Local == "si" {
					inSi = true
				} else if t.Name.Local == "t" && inSi {
					inT = true
				}

			case xml.EndElement:
				if t.Name.Local == "si" {
					inSi = false
					sharedStrings = append(sharedStrings, currentString.String())
					currentString.Reset()
				} else if t.Name.Local == "t" {
					inT = false
				}
			default:
		}
	}

	return
}

