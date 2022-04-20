package main

import (
	"fmt"
	"github.com/nagygr/ooxml2txt/archive"
	"log"
)

func main() {
	zipfile, err := archive.MakeZipFile("test_data/example.docx")

	if err != nil {
		log.Fatal("Couldn't open zip file: %s", err)
	}

	fmt.Printf("Number of files: %d\n", len(zipfile.Files()))
}
