package fileserver

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

type zipRenderer struct {
	w *zip.Writer

	root string
}

func (z *zipRenderer) Init(w io.Writer, title string) error {
	z.w = zip.NewWriter(w)
	return nil
}

func (z *zipRenderer) Walk(path string, finfo os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if finfo.IsDir() {
		return nil
	}

	fh, err := zip.FileInfoHeader(finfo)
	if err != nil {
		return err
	}

	fh.Name = "download/" + path
	w, err := z.w.CreateHeader(fh)
	if err != nil {
		return err
	}

	f, err := os.Open(filepath.Join(z.root, path))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(w, f)
	return err
}

func (z *zipRenderer) Flush() error {
	return z.w.Close()
}
