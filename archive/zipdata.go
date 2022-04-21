package archive

/*
 * Using code snippets from:
 * https://github.com/nguyenthenguyen/docx/
 * (c) Nguyen The Nguyen
 */

import (
	"archive/zip"
)

type ZipData interface {
	Files() []*zip.File
	close() error

	FileByName(name string) (file *zip.File, err error)
	FilesByName(substring string) (files []*zip.File, err error)
}
