package file

import (
	"context"
	"io"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
)

const (
	timeoutSeconds = 60
	minPartsGCS    = 3
)

func getContentFromGCS(ctx context.Context, path string) ([]byte, error) {
	log.Debug().Msgf("getContentFromGCS(%s)", path)

	if path == "" {
		return nil, errors.New("path is required")
	}

	bucket, name, err := parseGCSPath(path)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing path")
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error creating storage client")
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*timeoutSeconds)
	defer cancel()

	b := client.Bucket(bucket)

	it := b.Objects(ctx, &storage.Query{Prefix: name})
	_, err = it.Next()
	if errors.Is(err, iterator.Done) {
		return nil, errors.New("object not found")
	}

	rc, err := b.Object(name).NewReader(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error creating reader")
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, errors.Wrap(err, "error reading data")
	}

	return data, nil
}

func parseGCSPath(path string) (bucket, name string, err error) {
	if path == "" {
		return "", "", errors.New("path is required")
	}

	parts := strings.Split(path, "/")
	if len(parts) < minPartsGCS {
		return "", "", errors.New("invalid path")
	}

	bucket = parts[2]
	name = strings.Join(parts[3:], "/")
	return bucket, name, nil
}
