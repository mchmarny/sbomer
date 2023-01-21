package file

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	defaultHTTPRequestTimeout = 30 * time.Second
)

func getContentFromURL(ctx context.Context, url string) ([]byte, error) {
	log.Debug().Msgf("getContentFromURL(%s)", url)
	if url == "" {
		return nil, errors.New("url is required")
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request")
	}

	req = req.WithContext(ctx)

	c := &http.Client{
		Timeout: defaultHTTPRequestTimeout,
	}
	r, err := c.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error executing request")
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(err, "invalid response code: %s", r.Status)
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading body")
	}

	return b, nil
}
