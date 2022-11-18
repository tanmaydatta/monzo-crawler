package config

type Config struct {
	WorkersToFetchURL int `mapstructure:"workers"`
	MaxDepth          int `mapstructure:"max_depth"`
	Env               string
	LogFile           string `mapstructure:"log_file"`
	OutFile           string `mapstructure:"out_file"`
	RequestTimeout    int    `mapstructure:"request_timeout"`
	MaxDelay          int    `mapstructure:"max_delay"`
	MaxAttempts       int    `mapstructure:"max_attempts"`
}

var Conf Config
