package logging

import (
	"io"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	EncodeConsole = "console"
	EncodeJSON    = "json"
)

// Config -
type Config struct {
	Out    io.Writer
	Level  string
	Encode string
	Color  bool
}

// Must -
func Must(cfg *Config) *zap.Logger {
	z, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return z
}

// New -
func New(cfg *Config) (*zap.Logger, error) {

	encCfg := zapcore.EncoderConfig{
		TimeKey:        "t",
		LevelKey:       "l",
		CallerKey:      "c",
		MessageKey:     "m",
		StacktraceKey:  "s",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var level zap.AtomicLevel
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	default:
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	if cfg.Color {
		encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var encoder zapcore.Encoder
	if cfg.Encode == EncodeConsole {
		encoder = zapcore.NewConsoleEncoder(encCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encCfg)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(cfg.Out), level)
	z := zap.New(core, zap.AddCaller())

	return z, nil
}
