package sbomer

import (
	"github.com/mchmarny/sbomer/pkg/sbomer"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	c "github.com/urfave/cli/v2"
)

var (
	fileFlag = &c.StringFlag{
		Name:    "file",
		Aliases: []string{"f"},
		Usage:   "sbom file location (e.g. ./sbom.json, ./dir, gs://bucket/sbom.json, http://null.io/sbom.json)",
	}

	quietFlag = &c.BoolFlag{
		Name:    "quiet",
		Aliases: []string{"q"},
		Usage:   "suppress output unless error",
	}

	importCmd = &c.Command{
		Name:    "import",
		Aliases: []string{"imp"},
		Usage:   "imports sbom data from file or image",
		Action:  importAction,
		Flags: []c.Flag{
			fileFlag,
			quietFlag,
		},
	}
)

func importAction(c *c.Context) error {
	r := &sbomer.Request{
		Path:  c.String(fileFlag.Name),
		Quiet: c.Bool(quietFlag.Name),
	}

	isQuiet(c)
	log.Info().Msgf(c.App.Version)

	if err := sbomer.Process(c.Context, r); err != nil {
		return errors.Wrap(err, "error executing import")
	}

	return nil
}
