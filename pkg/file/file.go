package file

import (
	"context"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// GetContent returns content from the given URI.
// Supported URI types: http://, https://, gs://.
// If not prefixed with protocol, assumed to be a file path.
// Relative, absolute, or directory paths are supported.
func GetContent(ctx context.Context, src string) (map[string][]byte, error) {
	log.Debug().Msgf("GetReader(%s)", src)
	if src == "" {
		return nil, errors.New("path is required")
	}

	m := make(map[string][]byte)
	test := strings.TrimSpace(strings.ToLower(src))

	// URL
	if strings.HasPrefix(test, "http://") || strings.HasPrefix(test, "https://") {
		b, err := getContentFromURL(ctx, src)
		if err != nil {
			return nil, errors.Wrapf(err, "error getting content from url: %s", src)
		}
		m[src] = b
		return m, nil
	}

	// OBJECT
	if strings.HasPrefix(test, "gs://") {
		b, err := getContentFromGCS(ctx, src)
		if err != nil {
			return nil, errors.Wrapf(err, "error getting content from gcs: %s", src)
		}
		m[src] = b
		return m, nil
	}

	// FILE
	fi, err := os.Stat(src)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to find file: %s", src)
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

	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to read file: %s", p)
		}
		m[p] = b
	}

	return m, nil
}
