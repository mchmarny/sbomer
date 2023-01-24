package out

import (
	"bytes"
	"embed"
	"os"
	"text/template"

	"github.com/pkg/errors"
)

const (
	ReportTemplate = "report.md"
)

var (
	//go:embed templates
	f embed.FS
)

func WriteTemplate(path, name string, data interface{}) error {
	b, err := ExecTemplate(name, data)
	if err != nil {
		return errors.Wrap(err, "error executing template")
	}

	if err := os.WriteFile(path, b, os.ModePerm); err != nil {
		return errors.Wrap(err, "error writing template")
	}

	return nil
}

func ExecTemplate(name string, data interface{}) ([]byte, error) {
	t := template.Must(template.New("").ParseFS(f, "templates/*"))

	buf := new(bytes.Buffer)
	if err := t.Lookup(name).Execute(buf, data); err != nil {
		return nil, errors.Wrap(err, "error executing template")
	}

	b := buf.Bytes()
	if len(b) == 0 {
		return nil, errors.New("empty content after template execution")
	}

	return buf.Bytes(), nil
}
