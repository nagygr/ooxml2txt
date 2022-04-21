package format

import (
	"github.com/nagygr/ooxml2txt/archive"
)

type Pptx struct {
	zipReader archive.ZipData
	Text      []string
}

func MakePptx(path string) (*Pptx, error) {
	reader, err := archive.MakeZipFile(path)

	if err != nil {
		return nil, err
	}

	slideXmls, err := readXmls(reader, "ppt/slides/slide")
	var slideTexts []string

	if err == nil {
		slideTexts, err = textListFromXmls(slideXmls)
	}

	if err != nil {
		slideTexts = []string{}
	}

	return &Pptx{
		zipReader: reader,
		Text: slideTexts}, nil

}
