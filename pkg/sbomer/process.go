package sbomer

import (
	"context"

	"github.com/mchmarny/sbomer/pkg/file"
	"github.com/mchmarny/sbomer/pkg/sbom"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func Process(ctx context.Context, req *Request) error {
	log.Debug().Msg("processing request...")
	if req == nil {
		return errors.New("request is required")
	}
	if err := req.Validate(); err != nil {
		return errors.Wrap(err, "invalid request")
	}

	m, err := file.GetContent(ctx, req.Path)
	if err != nil {
		return errors.Wrapf(err, "error getting content: %+v", req)
	}
	log.Debug().Msgf("found %d items", len(m))

	for k, v := range m {
		log.Debug().Msgf("processing %s", k)
		doc, err := sbom.ParseDoc(v)
		if err != nil {
			return errors.Wrapf(err, "error processing item: %s", k)
		}
		log.Debug().Msgf("found %s with %d items", doc.Subject, len(doc.Items))
		// TODO: process doc
	}

	return nil
}
