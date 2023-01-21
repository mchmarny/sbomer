package sbomer

import (
	"fmt"
	"time"

	"github.com/mchmarny/sbomer/pkg/sbomer"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	c "github.com/urfave/cli/v2"
)

const (
	name           = "sbomer"
	metaKeyVersion = "version"
	metaKeyCommit  = "commit"
	metaKeyDate    = "date"
)

var (
	targetFlag = &c.StringFlag{
		Name:    "target",
		Aliases: []string{"t"},
		Usage:   "data store to save results to (e.g. bq://my-project)",
	}

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

	execCmd = &c.Command{
		Name:    "import",
		Aliases: []string{"imp"},
		Usage:   "imports sbom data from file or image",
		Action:  execute,
		Flags: []c.Flag{
			fileFlag,
			targetFlag,
			quietFlag,
		},
	}
)

func Execute(version, commit, date string, args []string) error {
	app, err := newApp(version, commit, date)
	if err != nil {
		return err
	}

	if err := app.Run(args); err != nil {
		return errors.Wrap(err, "error running app")
	}
	return nil
}

func newApp(version, commit, date string) (*c.App, error) {
	if version == "" || commit == "" || date == "" {
		return nil, errors.New("version, commit, and date must be set")
	}

	compileTime, err := time.Parse("2006-01-02T15:04:05Z", date)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse date")
	}
	dateStr := compileTime.UTC().Format("2006-01-02 15:04 UTC")

	app := &c.App{
		EnableBashCompletion: true,
		Suggest:              true,
		Name:                 name,
		Version:              fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, dateStr),
		Usage:                `sbom data processor`,
		Compiled:             compileTime,
		Flags: []c.Flag{
			&c.BoolFlag{
				Name:  "debug",
				Usage: "verbose output",
				Action: func(c *c.Context, debug bool) error {
					if debug {
						zerolog.SetGlobalLevel(zerolog.DebugLevel)
					}
					return nil
				},
			},
		},
		Metadata: map[string]interface{}{
			metaKeyVersion: version,
			metaKeyCommit:  commit,
			metaKeyDate:    date,
		},
		Commands: []*c.Command{
			execCmd,
		},
	}

	return app, nil
}

func isQuiet(c *c.Context) bool {
	_, ok := c.App.Metadata["quiet"]
	if ok {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}

	return ok
}

func execute(c *c.Context) error {
	r := &sbomer.Request{
		Path:   c.String(fileFlag.Name),
		Target: c.String(targetFlag.Name),
		Quiet:  c.Bool(quietFlag.Name),
	}

	isQuiet(c)
	log.Info().Msgf(c.App.Version)

	if err := sbomer.Process(c.Context, r); err != nil {
		return errors.Wrap(err, "error executing")
	}

	return nil
}
