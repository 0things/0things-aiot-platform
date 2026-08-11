package log

import (
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ log.Logger = (*ZapLogger)(nil)

// ZapLogger is a logger impl.
type ZapLogger struct {
	log  *zap.Logger
	Sync func() error
}

// NewZapLogger creates a new zap logger.
func NewZapLogger(opts ...Option) *ZapLogger {
	options := &options{
		level:       zapcore.InfoLevel,
		development: false,
	}
	for _, o := range opts {
		o(options)
	}

	var encoder zapcore.Encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Always use JSON encoder for production
	encoder = zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		options.level,
	)

	zapOpts := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(3), // Skip wrapper calls
	}

	if options.development {
		zapOpts = append(zapOpts, zap.Development())
	}

	zapLogger := zap.New(core, zapOpts...)

	return &ZapLogger{
		log:  zapLogger,
		Sync: zapLogger.Sync,
	}
}

// Log implements log.Logger interface.
func (l *ZapLogger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn(fmt.Sprint("keyvals must appear in pairs: ", keyvals))
		return nil
	}

	var data []zap.Field
	var msg string
	for i := 0; i < len(keyvals); i += 2 {
		key := fmt.Sprint(keyvals[i])
		// Extract msg field as the zap message
		if key == "msg" {
			msg = fmt.Sprint(keyvals[i+1])
			continue
		}
		data = append(data, zap.Any(key, keyvals[i+1]))
	}

	switch level {
	case log.LevelDebug:
		l.log.Debug(msg, data...)
	case log.LevelInfo:
		l.log.Info(msg, data...)
	case log.LevelWarn:
		l.log.Warn(msg, data...)
	case log.LevelError:
		l.log.Error(msg, data...)
	case log.LevelFatal:
		l.log.Fatal(msg, data...)
	default:
		l.log.Debug(msg, data...)
	}
	return nil
}

// options defines the configuration options.
type options struct {
	level       zapcore.Level
	development bool
}

// Option is a function that configures the logger.
type Option func(*options)

// WithLevel sets the log level.
func WithLevel(level string) Option {
	return func(o *options) {
		switch level {
		case "debug":
			o.level = zapcore.DebugLevel
		case "info":
			o.level = zapcore.InfoLevel
		case "warn":
			o.level = zapcore.WarnLevel
		case "error":
			o.level = zapcore.ErrorLevel
		default:
			o.level = zapcore.InfoLevel
		}
	}
}

// WithDevelopment sets development mode.
func WithDevelopment(dev bool) Option {
	return func(o *options) {
		o.development = dev
	}
}
