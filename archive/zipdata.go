package archive

import (
	"archive/zip"
)

type ZipData interface {
	Files() []*zip.File
	close() error
}
