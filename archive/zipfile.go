package archive

/*
 * Using code snippets from:
 * https://github.com/nguyenthenguyen/docx/
 * (c) Nguyen The Nguyen
 */

import (
	"archive/zip"
	"errors"
	"fmt"
	"strings"
)

// ZipFile is an implementation of the ZipData interface for actual zip files.
type ZipFile struct {
	data *zip.ReadCloser
}

// MakeZipFile creates a ZipFile for an actual zip file given by its path.
func MakeZipFile(path string) (*ZipFile, error) {
	reader, err := zip.OpenReader(path)

	if err != nil {
		return nil, err
	}

	return &ZipFile{data: reader}, nil
}

// Files returns the files inside the zip archive.
func (z *ZipFile) Files() []*zip.File {
	return z.data.File
}

// close closes the zip reader.
func (z *ZipFile) close() error {
	return z.data.Close()
}

// FileByName finds the file with the given name or returns an error.
func (z *ZipFile) FileByName(name string) (file *zip.File, err error) {
	for _, f := range z.data.File {
		if f.Name == name {
			file = f
			break
		}
	}

	if file == nil {
		err = errors.New(fmt.Sprintf("The file called %s not found", name))
	}

	return

}

// FilesByName finds all the files containing the given substring or
// return an error.
func (z *ZipFile) FilesByName(substring string) (files []*zip.File, err error) {
	for _, f := range z.data.File {
		if strings.Contains(f.Name, substring) {
			files = append(files, f)
		}
	}

	if len(files) == 0 {
		err = errors.New(fmt.Sprintf("No file containing \"%s\" found", substring))
	}

	return
}
