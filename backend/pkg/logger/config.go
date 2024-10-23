package logger

// Config holds configuration for the logger
type Config struct {
	Debug      bool
	Filename   string
	MaxSize    int // megabytes
	MaxBackups int
	MaxAge     int // days
	Compress   bool
}

var ConfigDefault = Config{
	Debug:      false,
	Filename:   "logs/app.log",
	MaxSize:    1, // megabytes
	MaxBackups: 3,
	MaxAge:     28, // days
	Compress:   true,
}

// defaultConfig returns a default config for http config.
func defaultConfig(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	return cfg
}
