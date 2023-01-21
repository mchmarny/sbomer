package sbomer

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	JSONFormat OutputFormat = iota
	YAMLFormat

	yamlIndent = 2
)

type OutputFormat int64

func write(w io.Writer, data interface{}, format OutputFormat) error {
	if data == nil {
		return errors.New("nil data")
	}

	switch format {
	case JSONFormat:
		j := json.NewEncoder(w)
		j.SetIndent("", "  ")
		if err := j.Encode(data); err != nil {
			return errors.Wrap(err, "error encoding")
		}
	case YAMLFormat:
		y := yaml.NewEncoder(w)
		y.SetIndent(yamlIndent)
		if err := y.Encode(data); err != nil {
			return errors.Wrap(err, "error encoding")
		}
	default:
		return errors.Errorf("unsupported output format: %d", format)
	}

	return nil
}
