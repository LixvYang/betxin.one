package logger

import (
	"io"
	"os"
	"path"
	"runtime/debug"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

// Configuration for logging
type LogConfig struct {
	ConsoleLoggingEnabled bool   `mapstructure:"console_logging_enabled"`
	EncodeLogsAsJson      bool   `mapstructure:"encode_logs_as_json"`
	FileLoggingEnabled    bool   `mapstructure:"file_logging_enabled"`
	Directory             string `mapstructure:"directory"`
	Filename              string `mapstructure:"filename"`
	MaxSize               int    `mapstructure:"max_size"`
	MaxBackups            int    `mapstructure:"max_backups"`
	MaxAge                int    `mapstructure:"max_age"`
	Level                 int    `mapstructure:"level"`
	LocalTime             bool   `mapstructure:"local_time"`
	Compress              bool   `mapstructure:"compress"`
}

// Configure sets up the logging framework
//
// In production, the container logs will be collected and file logging should be disabled. However,
// during development it's nicer to see logs as text and optionally write to a file when debugging
// problems in the containerized pipeline
//
// The output log file will be located at /var/log/service-xyz/service-xyz.log and
// will be rolled according to configuration set.
func New(config *LogConfig) *zerolog.Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	var logLevel = zerolog.Level(config.Level)
	if config.Level < -1 || config.Level > 7 {
		logLevel = zerolog.InfoLevel // default to INFO
	}

	var writers []io.Writer

	if config.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
			FieldsExclude: []string{
				"user_agent",
				"git_revision",
				"go_version",
			}})
	}

	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}
	mw := io.MultiWriter(writers...)

	var gitRevision string

	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		for _, v := range buildInfo.Settings {
			if v.Key == "vcs.revision" {
				gitRevision = v.Value
				break
			}
		}
	}

	Lg := zerolog.New(mw).
		Level(zerolog.Level(logLevel)).
		With().
		Str("git_revision", gitRevision).
		Str("go_version", buildInfo.GoVersion).
		Timestamp().
		CallerWithSkipFrameCount(2).
		Logger()

	Lg.Info().
		Bool("fileLogging", config.FileLoggingEnabled).
		Bool("jsonLogOutput", config.EncodeLogsAsJson).
		Str("logDirectory", config.Directory).
		Str("fileName", config.Filename).
		Int("maxSizeMB", config.MaxSize).
		Int("maxBackups", config.MaxBackups).
		Int("maxAgeInDays", config.MaxAge).
		Bool("localtime", config.LocalTime).
		Bool("compress", config.Compress).
		Msg("logging configured")

	return &Lg
}

func newRollingFile(config *LogConfig) io.Writer {
	if err := os.MkdirAll(config.Directory, 0744); err != nil {
		log.Error().Err(err).Str("path", config.Directory).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
		LocalTime:  config.LocalTime,
		Compress:   config.Compress,
	}
}
