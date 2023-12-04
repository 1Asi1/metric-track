package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger() zerolog.Logger {
	out := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006-01-02 15:04:05 -0700",
		NoColor:    true,
	}
	return zerolog.New(out)
}
