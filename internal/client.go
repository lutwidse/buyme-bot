package internal

import (
	"buyme-bot/pkg/captcha"
	"buyme-bot/pkg/config"
	"buyme-bot/pkg/proxy"
	"buyme-bot/pkg/util"

	"go.uber.org/zap"
)

type ClientFactory struct {
	Logger  *zap.SugaredLogger
	Config  *config.Config
	Proxy   *proxy.Proxy
	Util    *util.Util
	Captcha *captcha.Client
}

func NewClientFactory(logger *zap.SugaredLogger) *ClientFactory {
	return &ClientFactory{
		Logger: logger,
	}
}

func (cf *ClientFactory) Init() {
	cf.Config = cf.newConfig()
	cf.Proxy = cf.newProxy()
	cf.Util = cf.newUtil()
	cf.Captcha = cf.newCaptcha()
}

func (cf *ClientFactory) newConfig() *config.Config {
	Config := &config.Config{Logger: cf.Logger, ConfigData: &config.ConfigData{}}
	Config.LoadConfig()
	return Config
}

func (cf *ClientFactory) newProxy() *proxy.Proxy {
	return &proxy.Proxy{Logger: cf.Logger, Config: cf.Config}
}

func (cf *ClientFactory) newUtil() *util.Util {
	return &util.Util{Logger: cf.Logger}
}

func (cf *ClientFactory) newCaptcha() *captcha.Client {
	return captcha.NewClient(cf.Config.ConfigData.TwoCaptcha.Token)
}
