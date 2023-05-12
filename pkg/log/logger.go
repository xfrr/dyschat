package log

import (
	"os"

	"github.com/rs/zerolog"
)

func NewZeroLogger(lvl zerolog.Level) *zerolog.Logger {
	output := zerolog.ConsoleWriter{
		Out: os.Stdout,
		// FormatTimestamp: func(i interface{}) string {
		// 	parse, _ := time.Parse(time.RFC3339, i.(string))
		// 	return parse.Format
		// },
	}

	l := zerolog.New(output).With().
		Timestamp().
		// CallerWithSkipFrameCount(2).
		Logger()

	l.Level(lvl)
	return &l
}

func ParseLogLevel(lvl string) zerolog.Level {
	switch lvl {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.DebugLevel
	}
}
