package format

import (
	"github.com/nagygr/ooxml2txt/internal/archive"
	. "github.com/nagygr/ooxml2txt/internal/format"
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

	return makePptxFromReader(reader)
}

// MakePptxFromUrl creates a Pptx from an URL to a presentation. The returned
// instance contains the valid contents of the document if there was no error
// while processing it (which is then reported in the returned error value).
func MakePptxFromUrl(url string) (*Pptx, error) {
	reader, err := archive.MakeZipFileFromUrl(url)

	if err != nil {
		return nil, err
	}

	return makePptxFromReader(reader)
}

func makePptxFromReader(reader archive.ZipData) (*Pptx, error) {
	slideXmls, err := ReadXmls(reader, "ppt/slides/slide")
	var slideTexts []string

	if err == nil {
		slideTexts, err = TextListFromXmls(slideXmls)
	}

	if err != nil {
		slideTexts = []string{}
	}

	return &Pptx{
		zipReader: reader,
		Text:      slideTexts}, nil
}
