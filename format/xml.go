package format

import (
	"archive/zip"
	"github.com/nagygr/ooxml2txt/archive"
	"io"
	"io/ioutil"
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
