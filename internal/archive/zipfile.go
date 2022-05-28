package archive

/*
 * Using code snippets from:
 * https://github.com/nguyenthenguyen/docx/
 * (c) Nguyen The Nguyen
 */

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type zipReader interface {
	Files() []*zip.File
	Close() error
}

type zipFileReader struct {
	reader *zip.ReadCloser
}

var _ zipReader = (*zipFileReader)(nil)

func (r *zipFileReader) Files() []*zip.File {
	return r.reader.File
}

func (r *zipFileReader) Close() error {
	return r.reader.Close()
}

type zipUrlReader struct {
	reader *zip.Reader
}

var _ zipReader = (*zipUrlReader)(nil)

func (r *zipUrlReader) Files() []*zip.File {
	return r.reader.File
}

func (r *zipUrlReader) Close() error {
	return nil
}

// ZipFile is an implementation of the ZipData interface for actual zip files.
type ZipFile struct {
	data zipReader
}

// MakeZipFile creates a ZipFile for an actual zip file given by its path.
func MakeZipFile(path string) (*ZipFile, error) {
	reader, err := zip.OpenReader(path)

	if err != nil {
		return nil, err
	}

	return &ZipFile{data: &zipFileReader{reader}}, nil
}

// MakeZipFileFromUrl creates a ZipFile from a URL
func MakeZipFileFromUrl(url string) (*ZipFile, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	reader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, err
	}

	return &ZipFile{data: &zipUrlReader{reader}}, nil
}

// Files returns the files inside the zip archive.
func (z *ZipFile) Files() []*zip.File {
	return z.data.Files()
}

// close closes the zip reader.
func (z *ZipFile) close() error {
	return z.data.Close()
}

// FileByName finds the file with the given name or returns an error.
func (z *ZipFile) FileByName(name string) (file *zip.File, err error) {
	for _, f := range z.data.Files() {
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
	for _, f := range z.data.Files() {
		if strings.Contains(f.Name, substring) {
			files = append(files, f)
		}
	}

	if len(files) == 0 {
		err = errors.New(fmt.Sprintf("No file containing \"%s\" found", substring))
	}

	return
}
