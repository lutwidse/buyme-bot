package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	Logger     *zap.SugaredLogger
	ConfigData *ConfigData
}

type ConfigData struct {
	TwoCaptcha struct {
		Token string
	}
	Discord struct {
		Token     string
		ChannelID string
		RoleID    string
	}
	OpenAI struct {
		Token string
	}
	Proxy struct {
		Host     string
		Port     string
		User     string
		Password string
	}
	Setting struct {
		UserAgent string
	}
}

func (c *Config) LoadConfig() error {
	err := godotenv.Load(".env")
	if err != nil {
		c.Logger.Errorf("Error loading .env file")
		return err
	}

	c.ConfigData = &ConfigData{
		TwoCaptcha: struct {
			Token string
		}{
			Token: os.Getenv("TWO_CAPTCHA_TOKEN"),
		},
		Discord: struct {
			Token     string
			ChannelID string
			RoleID    string
		}{
			Token:     os.Getenv("DISCORD_TOKEN"),
			ChannelID: os.Getenv("DISCORD_CHANNEL_ID"),
			RoleID:    os.Getenv("DISCORD_ROLE_ID"),
		},
		OpenAI: struct {
			Token string
		}{
			Token: os.Getenv("OPENAI_TOKEN"),
		},
		Proxy: struct {
			Host     string
			Port     string
			User     string
			Password string
		}{
			Host:     os.Getenv("PROXY_HOST"),
			Port:     os.Getenv("PROXY_PORT"),
			User:     os.Getenv("PROXY_USER"),
			Password: os.Getenv("PROXY_PASSWORD"),
		},
		Setting: struct {
			UserAgent string
		}{
			UserAgent: os.Getenv("USER_AGENT"),
		},
	}

	return nil
}
