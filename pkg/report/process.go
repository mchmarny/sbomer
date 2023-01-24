package report

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/mchmarny/sbomer/pkg/out"
	"github.com/mchmarny/sbomer/pkg/sbom"
	"github.com/mchmarny/sbomer/pkg/sbom/cdx"
	"github.com/mchmarny/sbomer/pkg/sbom/spdx"
	"github.com/mchmarny/sbomer/pkg/src"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	sbomFormatSPDX = "spdx"
	sbomFormatCDX  = "cdx"
)

type templateData struct {
	GeneratedOn string
	Items       []*sbom.ReportResult
}

func Process(ctx context.Context, req *Request) error {
	log.Debug().Msg("processing request...")
	if req == nil {
		return errors.New("request is required")
	}
	if err := req.Validate(); err != nil {
		return errors.Wrap(err, "invalid request")
	}

	sources, err := src.GetFiles(ctx, req.Path)
	if err != nil {
		return errors.Wrapf(err, "error getting content: %+v", req)
	}
	log.Debug().Msgf("found %d items", len(sources))

	list := make([]*sbom.ReportResult, 0)

	for _, s := range sources {
		log.Debug().Msgf("processing %s", s)

		b, err := os.ReadFile(s)
		if err != nil {
			log.Error().Msgf("error reading file %s: %v", s, err)
			continue
		}

		if len(b) == 0 {
			log.Warn().Msgf("empty file %s", s)
			continue
		}

		sbomType := determineSbomType(b)
		log.Debug().Msgf("type %s", sbomType)

		var r sbom.SbomReport

		switch sbomType {
		case "spdx":
			r = spdx.GetSpdxReport(b)
		case "cdx":
			r = cdx.GetCycloneDXReport(b)
		}

		rep := sbom.GetReport(r)
		list = append(list, &rep)
	}

	log.Debug().Msgf("found %d items in %s", len(list), req.Path)

	d := &templateData{
		GeneratedOn: time.Now().Format(time.RFC1123),
		Items:       list,
	}

	if err := out.WriteTemplate(req.Target, out.ReportTemplate, d); err != nil {
		return errors.Wrapf(err, "error writing output to %s", req.Target)
	}

	return nil
}

func determineSbomType(b []byte) string {
	if strings.Contains(strings.ToLower(string(b)), sbomFormatSPDX) {
		return sbomFormatSPDX
	}
	return sbomFormatCDX
}
