package logging

import (
	"io"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	encodeConsole = "console"
	encodeJSON    = "json"
)

// Config -
type Config struct {
	Out   io.Writer
	Level string
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
	var encode string
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
		encode = encodeConsole
		encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
		encode = encodeJSON
	}

	var encoder zapcore.Encoder
	if encode == encodeConsole {
		encoder = zapcore.NewConsoleEncoder(encCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encCfg)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(cfg.Out), level)
	z := zap.New(core, zap.AddCaller())

	return z, nil
}
