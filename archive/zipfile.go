package archive

import (
	"archive/zip"
)

type ZipFile struct {
	data *zip.ReadCloser
}

func MakeZipFile(path string) (*ZipFile, error) {
	reader, err := zip.OpenReader(path)

	if err != nil {
		return nil, err
	}

	return &ZipFile{data: reader}, nil
}

func (z *ZipFile) Files() []*zip.File {
	return z.data.File
}

func (z *ZipFile) close() error {
	return z.data.Close()
}
