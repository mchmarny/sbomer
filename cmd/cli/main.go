package main

import (
	"os"

	sbomer "github.com/mchmarny/sbomer/cmd/cli/sbomer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	logLevelEnvVar = "debug"
)

var (
	version = "v0.0.1-default"
	commit  = "none"
	date    = "unknown"
)

func main() {
	initLogging()
	if err := sbomer.Execute(version, commit, date, os.Args); err != nil {
		log.Error().Msg(err.Error())
	}
}

func initLogging() {
	level := zerolog.InfoLevel
	levStr := os.Getenv(logLevelEnvVar)
	if levStr == "true" {
		level = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(level)

	out := zerolog.ConsoleWriter{
		Out: os.Stderr,
		PartsExclude: []string{
			zerolog.TimestampFieldName,
		},
	}

	log.Logger = zerolog.New(out)
}
