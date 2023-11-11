package config

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Logger *zap.SugaredLogger
	ConfigData *ConfigData
}

type ConfigData struct {
	TwoCaptcha struct {
		Token string `yaml:"token"`
	} `yaml:"twocaptcha"`
	Discord struct {
		Token     string `yaml:"token"`
		ChannelID string `yaml:"channel_id"`
		RoleID    string `yaml:"role_id"`
	} `yaml:"discord"`
	OpenAI struct {
		Token string `yaml:"token"`
	} `yaml:"openai"`
	Proxy struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"proxy"`
}

func (c *Config) LoadConfig() error {
	currentDir, err := os.Getwd()
	if err != nil {
		c.Logger.Errorf("Failed to get current directory: %v", err)
		return err
	}
	configPath := filepath.Join(currentDir, "..", "configs", "config.yml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}

	return nil
}
