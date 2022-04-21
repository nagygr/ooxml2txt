package main

import (
	"fmt"
	"github.com/nagygr/ooxml2txt/archive"
	"github.com/nagygr/ooxml2txt/format"
	"log"
)

func main() {
	zipfile, err := archive.MakeZipFile("test_data/example.docx")

	if err != nil {
		log.Fatal("Couldn't open zip file: %s", err)
	}

	fmt.Printf("Number of files: %d\n", len(zipfile.Files()))

	docxfile, err := format.MakeDocx("test_data/example.docx")

	if err != nil {
		log.Fatal(fmt.Sprintf("Couldn't open docx file: %s", err))
	}

	fmt.Printf("Text: %s\n", docxfile.Text)

	fmt.Printf("Links: %s\n", docxfile.Links)

	fmt.Printf("Footnotes: %s\n", docxfile.Footnotes)

	fmt.Printf("Headers: %s\n", docxfile.Headers)
	fmt.Printf("Footers: %s\n", docxfile.Footers)

	pptxfile, err := format.MakePptx("test_data/example.pptx")

	if err != nil {
		log.Fatal(fmt.Sprintf("Couldn't open pptx file: %s", err))
	}

	fmt.Printf("Text: %s\n", pptxfile.Text)

	ppt, _ := format.MakePptx("test_data/example.pptx")

	for n, slide := range ppt.Text {
		fmt.Printf("Slide %d: %s\n", n, slide)
	}
}
