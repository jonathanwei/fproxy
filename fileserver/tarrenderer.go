package fileserver

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

type tarRenderer struct {
	w *tar.Writer

	root string
}

func (t *tarRenderer) Init(w io.Writer, title string) error {
	t.w = tar.NewWriter(w)
	return t.w.WriteHeader(&tar.Header{
		Name:     "download/",
		Typeflag: tar.TypeDir,
		Mode:     0755 | 040000, // rwx + rx and c_ISDIR from archive/tar.
	})
}

func (t *tarRenderer) Walk(path string, finfo os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	fh, err := tar.FileInfoHeader(finfo, "")
	if err != nil {
		return err
	}

	fh.Name = "download/" + path
	err = t.w.WriteHeader(fh)
	if err != nil {
		return err
	}

	if finfo.IsDir() {
		return nil
	}

	f, err := os.Open(filepath.Join(t.root, path))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(t.w, f)
	return err
}

func (t *tarRenderer) Flush() error {
	return t.w.Close()
}
