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
		log.Fatal("Couldn't open docx file: %s", err)
	}

	fmt.Printf("Text: %s\n", docxfile.Text())

	fmt.Printf("Links: %s\n", docxfile.Links())
}
