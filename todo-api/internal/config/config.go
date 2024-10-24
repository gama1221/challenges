package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

// AppConfig holds application-specific configurations.
type AppConfig struct {
	Env string `mapstructure:"env"`
}

// DBConfig holds database configurations.
type DBConfig struct {
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	Keyspace       string `mapstructure:"keyspace"`
	ContactPoints  string `mapstructure:"contact_points"`
	Host           string `mapstructure:"host"`
	Port           string `mapstructure:"port"`
	Name           string `mapstructure:"name"`
	Driver         string `mapstructure:"driver"`
	MigrationsPath string `mapstructure:"migrations_path"`
}

// ServerConfig holds server configurations.
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// Config holds the configuration for the application.
type Config struct {
	App    AppConfig    `mapstructure:"app"`
	DB     DBConfig     `mapstructure:"db"`
	Server ServerConfig `mapstructure:"server"`
}

var cfg Config

// Get returns the global configuration.
func Get() Config {
	return cfg
}

// ReadConfig reads and parses the configuration file.
func ReadConfig(filename string) (Config, error) {
	if len(filename) < 1 {
		filename = "config.yml"
	}

	viper.SetConfigFile(filename)
	viper.SetConfigType("yaml")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	for _, k := range viper.AllKeys() {
		v := viper.GetString(k)
		viper.Set(k, os.ExpandEnv(v))
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	cfg = config
	return config, nil
}
