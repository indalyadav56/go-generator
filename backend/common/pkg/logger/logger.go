package logger

import (
	"os"
	"reflect"
	"strings"
	"time"

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

func New(config ...Config) Logger {
	cfg := defaultConfig(config...)

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	fileSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		// Compress:   cfg.Compress,
	})

	fileCore := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, fileSyncer, zap.InfoLevel),
	)

	consoleCore := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	core := zapcore.NewTee(consoleCore, fileCore)

	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))

	if cfg.Debug {
		log = log.WithOptions(zap.Development())
	}

	return &zapLogger{log: log}
}

func structToFields(obj interface{}) []zap.Field {
	if obj == nil {
		return nil
	}

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	t := v.Type()
	// Pre-allocate slice with capacity of struct fields
	fields := make([]zap.Field, 0, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanInterface() {
			continue
		}

		fieldType := t.Field(i)
		fieldName := parseFieldName(fieldType)
		if fieldName == "-" {
			continue
		}

		// Skip empty fields if omitempty is set
		if shouldOmitEmpty(fieldType) && isEmptyValue(field) {
			continue
		}

		if zapField := valueToField(fieldName, field); zapField != nil {
			fields = append(fields, *zapField)
		}
	}

	return fields
}

// parseFieldName extracts the field name from struct tags or falls back to field name
func parseFieldName(field reflect.StructField) string {
	if tag := field.Tag.Get("json"); tag != "" {
		if idx := strings.Index(tag, ","); idx != -1 {
			return tag[:idx]
		}
		return tag
	}
	return field.Name
}

// shouldOmitEmpty checks if the field should be omitted when empty
func shouldOmitEmpty(field reflect.StructField) bool {
	if tag := field.Tag.Get("json"); tag != "" {
		return strings.Contains(tag, "omitempty")
	}
	return false
}

// isEmptyValue determines if a field's value should be considered empty
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Struct:
		if t, ok := v.Interface().(time.Time); ok {
			return t.IsZero()
		}
	}
	return false
}

// valueToField converts a reflect.Value to a *zap.Field based on its kind
func valueToField(name string, v reflect.Value) *zap.Field {
	switch v.Kind() {
	case reflect.String:
		val := v.String()
		if val != "" {
			field := zap.String(name, val)
			return &field
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := v.Int()
		if val != 0 {
			field := zap.Int64(name, val)
			return &field
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val := v.Uint()
		if val != 0 {
			field := zap.Uint64(name, val)
			return &field
		}

	case reflect.Float32, reflect.Float64:
		val := v.Float()
		if val != 0 {
			field := zap.Float64(name, val)
			return &field
		}

	case reflect.Bool:
		val := v.Bool()
		field := zap.Bool(name, val)
		return &field

	case reflect.Struct:
		if t, ok := v.Interface().(time.Time); ok {
			if !t.IsZero() {
				field := zap.Time(name, t)
				return &field
			}
			return nil
		}
		fallthrough

	default:
		if v.IsValid() && !isEmptyValue(v) {
			field := zap.Any(name, v.Interface())
			return &field
		}
	}

	return nil
}

// processFields converts interface{} slice to zap.Field slice with proper key handling
func (l *zapLogger) processFields(fields ...interface{}) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))

	for i := 0; i < len(fields); i++ {
		field := fields[i]

		switch f := field.(type) {
		case zap.Field:
			zapFields = append(zapFields, f)

		case string:
			// Handle string key-value pairs
			if i+1 < len(fields) {
				// If next item exists, treat current as key and next as value
				if key := sanitizeKey(f); key != "" {
					zapFields = append(zapFields, processKeyValue(key, fields[i+1]))
					i++ // Skip next item since we used it as value
				}
			}

		case map[string]interface{}:
			for key, value := range f {
				if key = sanitizeKey(key); key != "" {
					zapFields = append(zapFields, processKeyValue(key, value))
				}
			}

		default:
			zapFields = append(zapFields, structToFields(f)...)
		}
	}

	return zapFields
}

// sanitizeKey ensures consistent key naming
func sanitizeKey(key string) string {
	// Remove spaces and convert to snake_case
	key = strings.TrimSpace(strings.ToLower(key))
	key = strings.ReplaceAll(key, " ", "_")
	return key
}

// processKeyValue converts a key-value pair to zap.Field
func processKeyValue(key string, value interface{}) zap.Field {
	switch v := value.(type) {
	case string:
		return zap.String(key, v)
	case int:
		return zap.Int(key, v)
	case int64:
		return zap.Int64(key, v)
	case float64:
		return zap.Float64(key, v)
	case bool:
		return zap.Bool(key, v)
	case error:
		return zap.Error(v)
	default:
		return zap.Any(key, v)
	}
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
	_ = l.log.Sync()
}
