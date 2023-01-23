package sbomer

import (
	"github.com/mchmarny/sbomer/pkg/report"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	c "github.com/urfave/cli/v2"
)

var (
	srcFlag = &c.StringFlag{
		Name:    "src",
		Aliases: []string{"s"},
		Usage:   "path to sbom file(s) (e.g. ./sbom.json, ./dir, gs://bucket/sbom.json, http://null.io/sbom.json)",
	}

	outFlag = &c.StringFlag{
		Name:    "out",
		Aliases: []string{"o"},
		Usage:   "path where output will be written",
	}

	importCmd = &c.Command{
		Name:   "report",
		Usage:  "generates SBOM report",
		Action: importAction,
		Flags: []c.Flag{
			srcFlag,
			outFlag,
		},
	}
)

func importAction(c *c.Context) error {
	r := &report.Request{
		Path:   c.String(srcFlag.Name),
		Target: c.String(outFlag.Name),
	}

	isQuiet(c)
	log.Info().Msgf(c.App.Version)

	if err := report.Process(c.Context, r); err != nil {
		return errors.Wrap(err, "error executing import")
	}

	return nil
}
