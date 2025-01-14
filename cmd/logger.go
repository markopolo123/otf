package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
)

const DefaultLogLevel = "info"

type LoggerConfig struct {
	Level string

	// Toggle log colors. Must be one of auto, true, or false.
	Color string
}

// NewLoggerConfigFromFlags adds flags to the given flagset, and, after the
// flagset is parsed by the caller, the flags populate the returned logger
// config.
func NewLoggerConfigFromFlags(flags *pflag.FlagSet) *LoggerConfig {
	cfg := LoggerConfig{}

	flags.StringVarP(&cfg.Level, "log-level", "l", DefaultLogLevel, "Logging level")
	flags.StringVar(&cfg.Color, "log-color", "auto", "Toggle log colors: auto, true or false. Auto enables colors if using a TTY.")

	return &cfg
}

func NewLogger(cfg *LoggerConfig) (logr.Logger, error) {
	// disable "v=<log-level>" field on every log line
	zerologr.VerbosityFieldName = ""
	// don't display floats for time durations
	zerolog.DurationFieldInteger = true

	zlvl, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return logr.Logger{}, err
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	// insert a delimiter, '|', between the message and the logfmt key values,
	// to aid log parsing
	consoleWriter.FormatMessage = func(msg interface{}) string {
		return fmt.Sprintf(`%s |`, msg)
	}

	switch cfg.Color {
	case "auto":
		// Disable color if stdout is not a tty
		if !isatty.IsTerminal(os.Stdout.Fd()) {
			consoleWriter.NoColor = true
		}
	case "true":
		consoleWriter.NoColor = false
	case "false":
		consoleWriter.NoColor = true
	default:
		return logr.Logger{}, fmt.Errorf("invalid choice for log color: %s. Must be one of auto, true, or false", cfg.Color)
	}

	zlogger := zerolog.New(consoleWriter).Level(zlvl).With().Timestamp().Logger()

	// wrap within logr wrapper
	return zerologr.New(&zlogger), nil
}
