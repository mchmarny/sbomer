package sbomer

import (
	"fmt"
	"time"

	"github.com/mchmarny/sbomer/pkg/config"
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

func Execute(version, commit, date string, args []string) error {
	homeDir, created, err := config.GetOrCreateHomeDir(name)
	if err != nil {
		return errors.Wrap(err, "failed to get home dir")
	}
	log.Debug().Msgf("home dir (created: %v): %s", created, homeDir)

	cfg, err := config.ReadOrCreate(homeDir)
	if err != nil {
		return errors.Wrap(err, "failed to read config")
	}

	app, err := newApp(version, commit, date)
	if err != nil {
		return err
	}

	if err := app.Run(args); err != nil {
		return errors.Wrap(err, "error running app")
	}

	cfg.LastExec = time.Now()
	cfg.LastVersion = version
	if err := config.Save(homeDir, cfg); err != nil {
		return errors.Wrap(err, "failed to save config")
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
			importCmd,
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
