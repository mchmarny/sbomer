package src

import (
	"context"
	"os"
	"path"

	"github.com/pkg/errors"
)

// GetFiles returns content from the given path.
func GetFiles(ctx context.Context, src string) ([]string, error) {
	if src == "" {
		return nil, errors.New("path is required")
	}

	// FILE
	fi, err := os.Stat(src)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to find path: %s", src)
	}

	paths := make([]string, 0)
	if fi.IsDir() {
		files, err := os.ReadDir(src)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to read dir: %s", src)
		}
		for _, f := range files {
			paths = append(paths, path.Join(src, f.Name()))
		}
	} else {
		paths = append(paths, src)
	}

	return paths, nil
}
