package config

import "time"

// Config is app configuration
type Config struct {
	HTTPClientTimeout         time.Duration `mapstructure:"http-client-timeout"`
	HTTPServerTimeout         time.Duration `mapstructure:"http-server-timeout"`
	HTTPServerShutdownTimeout time.Duration `mapstructure:"http-server-shutdown-timeout"`
	DataPath                  string        `mapstructure:"data-path"`
	Port                      string        `mapstructure:"port"`
	PortMetrics               int           `mapstructure:"port-metrics"`
	HostName                  string        `mapstructure:"hostname"`
	JWTSecret                 string        `mapstructure:"jwt-secret"`
	Redis                     *Redis        `mapstructure:"redis"`
	Logger                    *Logger       `mapstructure:"logger"`
	Database                  *Database     `mapstructure:"database"`
}

// Redis is redis configuration
type Redis struct {
	MaxIdle        int           `mapstructure:"MaxIdle"`
	IdleTimeout    time.Duration `mapstructure:"IdleTimeout"`
	ConnectTimeout time.Duration `mapstructure:"ConnectTimeout"`
	ReadTimeout    time.Duration `mapstructure:"ReadTimeout"`
	WriteTimeout   time.Duration `mapstructure:"WriteTimeout"`
	Host           string        `mapstructure:"Host"`
	Port           string        `mapstructure:"Port"`
	Username       string        `mapstructure:"Username"`
	Password       string        `mapstructure:"Password"`
	DB             int           `mapstructure:"DB"`
}

// Logger is redis configuration
type Logger struct {
	Level string `mapstructure:"level"`
}

// Database is mysql configuration
type Database struct {
	User         string        `mapstructure:"user"`
	Password     string        `mapstructure:"password"`
	DBName       string        `mapstructure:"dbname"`
	Host         string        `mapstructure:"host"`
	Port         string        `mapstructure:"port"`
	Timeout      time.Duration `mapstructure:"timeout"`
	ReadTimeout  time.Duration `mapstructure:"read-timeout"`
	WriteTimeout time.Duration `mapstructure:"write-timeout"`
}
