// Package proxy is designed to use Proxy-Cheap. proxy must be rotatable, but be static.
// This is because Cloudflare's Turnstile requires the same IP and User-Agent.
// Proxy-Cheap doesn't provide Session IP in API but we reverse engineer it.
// Tips: reCAPTCHA is not implemented SEC however they ban IP regularly every few days.
package proxy

import (
	"buyme-bot/pkg/config"
	"math/rand"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Proxy struct {
	Logger *zap.SugaredLogger
	Config *config.Config
}

func (p *Proxy) generateRandomString() string {
	characters := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, 8)
	for i := range result {
		result[i] = characters[rand.Intn(len(characters))]
	}
	return string(result)
}

func (p *Proxy) generatePassword(e map[string]string) string {
	password := e["password"]
	if e["country"] != "random" {
		password += "_country-" + strings.ReplaceAll(e["country"], " ", "")
	}
	if e["session"] == "sticky" {
		password += "_session-" + p.generateRandomString()
	}
	return password
}

func (p *Proxy) GetSessionProxy(country string) string {
	e := map[string]string{
		"connection": "not_ssl",
		"country":    country,
		"password":   p.Config.ConfigData.Proxy.Password,
		"session":    "sticky",
	}

	sessionProxy := p.generatePassword(e)

	return sessionProxy
}

func (p *Proxy) GetSessionProxyURL(country string) string {
	return p.Config.ConfigData.Proxy.User + ":" + p.GetSessionProxy(country) + "@" + p.Config.ConfigData.Proxy.Host + ":" + p.Config.ConfigData.Proxy.Port
}
