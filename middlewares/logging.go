package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Logging -
type Logging struct {
	logger  *zap.Logger
	sugar   *zap.SugaredLogger
	logging func(code int, r *http.Request, elapsed time.Duration)
}

// LoggingConfig -
type LoggingConfig struct {
	Logger  *zap.Logger
	Console bool
}

// MustLogging -
func MustLogging(cfg *LoggingConfig) *Logging {
	l, err := NewLogging(cfg)
	if err != nil {
		panic(err)
	}
	return l
}

// NewLogging -
func NewLogging(cfg *LoggingConfig) (*Logging, error) {
	l := &Logging{logger: cfg.Logger}
	l.logging = l.stdLogging

	if cfg.Console {
		// addCaller optionを適用したくない
		l.sugar = zap.New(l.logger.Core()).Sugar()
		l.logging = l.console

	}
	return l, nil
}

// ServeHTTP -
func (m *Logging) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	lw := &loggingWriter{ResponseWriter: w}
	start := time.Now()
	next.ServeHTTP(lw, r)
	elapsed := time.Since(start)

	// この時点で書き込まれていない場合がある。
	if lw.statusCode == 0 {
		lw.statusCode = 200
	}
	m.logging(lw.statusCode, r, elapsed)
}

func (m *Logging) stdLogging(code int, r *http.Request, elapsed time.Duration) {
	if code >= 400 {
		// Errorで出すと、logic側のloggingとstack traceが重複する. optionの指定でなんとかできるかもしれない.
		m.logger.Warn("req",
			zap.Int("code", code),
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.Float64("et", elapsed.Seconds()))
	} else {
		m.logger.Info("req",
			zap.Int("code", code),
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.Float64("et", elapsed.Seconds()))
	}
}

func (m *Logging) console(code int, r *http.Request, elapsed time.Duration) {
	if strings.HasPrefix(r.URL.Path, "/static") {
		return
	}
	msg := fmt.Sprintf("|%3d| %-4s %-40s %.3f", code, r.Method, r.URL.String(), elapsed.Seconds())
	if code >= 400 {
		m.sugar.Warn(msg)
	} else {
		m.sugar.Info(msg)
	}
}

type loggingWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lw *loggingWriter) WriteHeader(statusCode int) {
	lw.statusCode = statusCode
	lw.ResponseWriter.WriteHeader(statusCode)
}
