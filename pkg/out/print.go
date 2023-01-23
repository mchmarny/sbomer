package out

import (
	"encoding/json"
	"io/fs"
	"os"

	"github.com/pkg/errors"
)

type OutputFormat int64

func Write(path string, data interface{}) error {
	if data == nil {
		return errors.New("nil data")
	}

	b, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "error marshaling")
	}

	if err := os.WriteFile(path, b, fs.ModePerm); err != nil {
		return errors.Wrap(err, "error writing")
	}

	return nil
}
