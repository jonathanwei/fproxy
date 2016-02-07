package fileserver

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

// TODO: Use actual stylesheet.
var dirTmpl = template.Must(template.New("dir").Parse(`
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
		<style>
			* {
				box-sizing: border-box;
			}
			body {
				font-family: monospace;
			}
		</style>
	</head>
	<body>
		<div>
			<a href=".?m=zip" download>Download as zip</a>
		</div>
		{{- range .Files -}}
			<div>
				<a href="{{.}}">{{.}}</a>
			</div>
		{{- else -}}
			<div>
				<strong>No files.</strong>
			</div>
		{{- end -}}
	</body>
</html>`))

type htmlRenderer struct {
	w io.Writer

	// Used for HTML template.
	Title string
	Files []string
}

func (h *htmlRenderer) Init(w io.Writer, title string) {
	h.w = w
	h.Title = title
}

func (h *htmlRenderer) Walk(path string, finfo os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	h.Files = append(h.Files, path)

	if finfo.IsDir() {
		return filepath.SkipDir
	}

	return nil
}

func (h *htmlRenderer) Flush() error {
	var b bytes.Buffer
	err := dirTmpl.Execute(&b, h)
	if err != nil {
		return err
	}

	_, err = io.Copy(h.w, &b)
	return err
}
