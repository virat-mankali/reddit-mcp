package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	configDirName  = ".rdcli"
	configFileName = "config"
	configFileType = "json"
)

type Config struct {
	ClientID     string `mapstructure:"client_id" json:"client_id"`
	ClientSecret string `mapstructure:"client_secret" json:"client_secret"`
	UserAgent    string `mapstructure:"user_agent" json:"user_agent"`
}

func Load() (*Config, error) {
	dir, err := configDir()
	if err != nil {
		return nil, err
	}

	v := viper.New()
	v.SetConfigName(configFileName)
	v.SetConfigType(configFileType)
	v.AddConfigPath(dir)

	v.SetEnvPrefix("RD")
	v.BindEnv("client_id")
	v.BindEnv("client_secret")
	v.BindEnv("user_agent")
	v.BindEnv("client_id", "REDDIT_CLIENT_ID")
	v.BindEnv("client_secret", "REDDIT_CLIENT_SECRET")
	v.BindEnv("user_agent", "REDDIT_USER_AGENT")

	v.SetDefault("user_agent", "reddit-mcp/0.1.0")

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return nil, err
		}
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Save() error {
	dir, err := configDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	v := viper.New()
	v.SetConfigName(configFileName)
	v.SetConfigType(configFileType)
	v.AddConfigPath(dir)

	v.Set("client_id", c.ClientID)
	v.Set("client_secret", c.ClientSecret)
	v.Set("user_agent", c.UserAgent)

	path := filepath.Join(dir, configFileName+"."+configFileType)
	v.SetConfigFile(path)

	if _, err := os.Stat(path); err == nil {
		return v.WriteConfig()
	}
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return v.SafeWriteConfigAs(path)
}

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, configDirName), nil
}
