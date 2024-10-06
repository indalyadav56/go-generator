package logger

import (
	"os"
	"reflect"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	With(fields ...interface{}) Logger
	Close()
}

type zapLogger struct {
	log *zap.Logger
}

func New(config ...Config) (Logger, error) {
	cfg := defaultConfig(config...)

	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// Define our console output.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	// Configure console output.
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// Configure file output.
	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// encoderConfig.CallerKey = "caller"
	// encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// Create lumberjack syncer for file output
	fileSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize, // megabytes
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge, // days
		// Compress:   cfg.Compress,
	})

	// Create tee for file core
	fileCore := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, fileSyncer, zap.InfoLevel),
	)

	// Create tee for console core
	consoleCore := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	// Combine cores
	core := zapcore.NewTee(consoleCore, fileCore)

	// Create logger
	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))

	if cfg.Debug {
		log = log.WithOptions(zap.Development())
	}

	return &zapLogger{log: log}, nil
}

// structToFields converts a struct to a slice of zap.Field
func structToFields(obj interface{}) []zap.Field {
	fields := []zap.Field{}
	v := reflect.ValueOf(obj)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fields
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name

		switch field.Kind() {
		case reflect.String:
			fields = append(fields, zap.String(fieldName, field.String()))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fields = append(fields, zap.Int64(fieldName, field.Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fields = append(fields, zap.Uint64(fieldName, field.Uint()))
		case reflect.Float32, reflect.Float64:
			fields = append(fields, zap.Float64(fieldName, field.Float()))
		case reflect.Bool:
			fields = append(fields, zap.Bool(fieldName, field.Bool()))
		default:
			fields = append(fields, zap.Any(fieldName, field.Interface()))
		}
	}

	return fields
}

// processFields converts interface{} slice to zap.Field slice
func (l *zapLogger) processFields(fields ...interface{}) []zap.Field {
	zapFields := []zap.Field{}
	for _, field := range fields {
		switch f := field.(type) {
		case zap.Field:
			zapFields = append(zapFields, f)
		case string:
			zapFields = append(zapFields, zap.String("additionalInfo", f))
		default:
			zapFields = append(zapFields, structToFields(f)...)
		}
	}
	return zapFields
}

func (l *zapLogger) Debug(msg string, fields ...interface{}) {
	l.log.Debug(msg, l.processFields(fields...)...)
}

func (l *zapLogger) Info(msg string, fields ...interface{}) {
	l.log.Info(msg, l.processFields(fields...)...)
}

func (l *zapLogger) Warn(msg string, fields ...interface{}) {
	l.log.Warn(msg, l.processFields(fields...)...)
}

func (l *zapLogger) Error(msg string, fields ...interface{}) {
	l.log.Error(msg, l.processFields(fields...)...)
}

func (l *zapLogger) Fatal(msg string, fields ...interface{}) {
	l.log.Fatal(msg, l.processFields(fields...)...)
}

func (l *zapLogger) With(fields ...interface{}) Logger {
	return &zapLogger{log: l.log.With(l.processFields(fields...)...)}
}

func (l *zapLogger) Close() {
	_ = l.log.Sync() // flush buffered logs
}
