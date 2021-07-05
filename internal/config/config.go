package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config represents program configuration
type Config struct {
	HTTP     HTTP     `mapstructure:"http"`
	Redis    Redis    `mapstructure:"redis"`
	Logger   Logger   `mapstructure:"logger"`
	Database Database `mapstructure:"database"`
}

// New returns Config object that reads configurations from a file.
func New(configFile string) *Config {

	// Set default configurations
	setDefaults()

	// Set configuration file
	viper.SetConfigFile(configFile)

	// Automatically refresh environment variables
	viper.AutomaticEnv()

	// Read configuration
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("failed to read configuration: ", err.Error())
		os.Exit(1)
	}

	var config *Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println("HTTP server config unmarshal failed", err.Error())
	}

	return config
}

func setDefaults() {
	// Set default database configuration

}

// HTTP is http configuration
type HTTP struct {
	Port                      string        `mapstructure:"port"`
	PortMetrics               int           `mapstructure:"port-metrics"`
	HTTPServerTimeout         time.Duration `mapstructure:"http-server-timeout"`
	HTTPServerShutdownTimeout time.Duration `mapstructure:"http-server-shutdown-timeout"`
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
	Level            string   `mapstructure:"level"`
	OutputPaths      []string `mapstructure:"output-paths"`
	ErrorOutputPaths []string `mapstructure:"error-output-paths"`
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
