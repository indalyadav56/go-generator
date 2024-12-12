package logger

type Config struct {
	Debug      bool
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

var ConfigDefault = Config{
	Debug:      false,
	Filename:   "logs/app.log",
	MaxSize:    1,
	MaxBackups: 3,
	MaxAge:     28,
	Compress:   true,
}

func defaultConfig(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	return cfg
}
