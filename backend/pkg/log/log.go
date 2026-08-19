package log

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"sync"
	"time"
)

const ctxLoggerKey = "zapLogger"

type Logger struct {
	*zap.Logger
}

// DateWriteSyncer handles daily log rotation by appending date to filename.
type DateWriteSyncer struct {
	basePath string // e.g. "./storage/logs/server.log"
	mu       sync.Mutex
	hook     *lumberjack.Logger
	date     string // current date "2006-01-02"
	level    zapcore.Level
	maxSize  int
	maxBackups int
	maxAge    int
	compress  bool
}

func newDateWriteSyncer(basePath string, maxSize, maxBackups, maxAge int, compress bool) *DateWriteSyncer {
	return &DateWriteSyncer{
		basePath:   basePath,
		date:       time.Now().Format("2006-01-02"),
		maxSize:    maxSize,
		maxBackups: maxBackups,
		maxAge:     maxAge,
		compress:   compress,
	}
}

func (d *DateWriteSyncer) Write(p []byte) (n int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	today := time.Now().Format("2006-01-02")
	if today != d.date || d.hook == nil {
		d.rotate(today)
	}
	return d.hook.Write(p)
}

func (d *DateWriteSyncer) rotate(today string) {
	if d.hook != nil {
		d.hook.Close()
	}
	// Insert date before extension: server.log -> server-2025-01-15.log
	path := d.basePath
	if idx := strings.LastIndex(path, "."); idx != -1 {
		path = path[:idx] + "-" + today + path[idx:]
	} else {
		path = path + "-" + today
	}
	d.hook = &lumberjack.Logger{
		Filename:   path,
		MaxSize:    d.maxSize,
		MaxBackups: d.maxBackups,
		MaxAge:     d.maxAge,
		Compress:   d.compress,
	}
	d.date = today
}

func (d *DateWriteSyncer) Sync() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.hook != nil {
		return d.hook.Close()
	}
	return nil
}

func NewLog(conf *viper.Viper) *Logger {
	// log address "out.log" User-defined
	lp := conf.GetString("log.log_file_name")
	lv := conf.GetString("log.log_level")
	var level zapcore.Level
	//debug<info<warn<error<fatal<panic
	switch lv {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	hook := &DateWriteSyncer{
		basePath:   lp,
		date:       time.Now().Format("2006-01-02"),
		maxSize:    conf.GetInt("log.max_size"),
		maxBackups: conf.GetInt("log.max_backups"),
		maxAge:     conf.GetInt("log.max_age"),
		compress:   conf.GetBool("log.compress"),
	}

	var encoder zapcore.Encoder
	if conf.GetString("log.encoding") == "console" {
		encoder = zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "Logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
			EncodeTime:     timeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		})
	} else {
		encoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})
	}
	// default(both) log to console and file
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook)), // Print to console and file
		level,
	)
	mode := conf.GetString("log.mode")
	switch mode {
	case "console":
		core = zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
	case "file":
		core = zapcore.NewCore(
			encoder,
			zapcore.AddSync(hook),
			level,
		)
	}
	if conf.GetString("env") != "prod" {
		return &Logger{zap.New(core, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
	}
	return &Logger{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//enc.AppendString(t.Format("2006-01-02 15:04:05"))
	enc.AppendString(t.Format("2006-01-02 15:04:05.000000000"))
}

// WithValue Adds a field to the specified context
func (l *Logger) WithValue(ctx context.Context, fields ...zapcore.Field) context.Context {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
		c.Request = c.Request.WithContext(context.WithValue(ctx, ctxLoggerKey, l.WithContext(ctx).With(fields...)))
		return c
	}
	return context.WithValue(ctx, ctxLoggerKey, l.WithContext(ctx).With(fields...))
}

// WithContext Returns a zap instance from the specified context
func (l *Logger) WithContext(ctx context.Context) *Logger {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
	}
	zl := ctx.Value(ctxLoggerKey)
	ctxLogger, ok := zl.(*zap.Logger)
	if ok {
		return &Logger{ctxLogger}
	}
	return l
}
