package format

import (
	"github.com/nagygr/ooxml2txt/archive"
)

// Pptx handles pptx documents. The Text member is a list of strings where each
// element corresponds to a slide in the presentation.
type Pptx struct {
	zipReader archive.ZipData
	Text      []string
}

// MakePptx creates a Pptx from the path to a presentation. The
// returned instance contains the valid contents of the document if there was
// no error while processing it (which is then reported in the returned error
// value).
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
