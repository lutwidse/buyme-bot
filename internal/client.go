package internal

import (
	"buyme-bot/pkg/captcha"
	"buyme-bot/pkg/config"
	"buyme-bot/pkg/proxy"
	"buyme-bot/pkg/util"

	"go.uber.org/zap"
)

type ClientFactory struct {
	Logger        *zap.SugaredLogger
	ConfigClient  *config.ConfigClient
	ProxyClient   *proxy.Proxy
	UtilClient    *util.Util
	CaptchaClient *captcha.Client
}

func NewClientFactory(logger *zap.SugaredLogger) *ClientFactory {
	return &ClientFactory{
		Logger: logger,
	}
}

func (cf *ClientFactory) Init() {
	cf.ConfigClient = cf.newConfigClient()
	cf.ProxyClient = cf.newProxyClient()
	cf.UtilClient = cf.newUtilClient()
	cf.CaptchaClient = cf.newCaptchaClient()
}

func (cf *ClientFactory) newConfigClient() *config.ConfigClient {
	configClient := &config.ConfigClient{Logger: cf.Logger, Config: &config.Config{}}
	configClient.LoadConfig()
	return configClient
}

func (cf *ClientFactory) newProxyClient() *proxy.Proxy {
	return &proxy.Proxy{Logger: cf.Logger, ConfigClient: cf.ConfigClient}
}

func (cf *ClientFactory) newUtilClient() *util.Util {
	return &util.Util{Logger: cf.Logger}
}

func (cf *ClientFactory) newCaptchaClient() *captcha.Client {
	return captcha.NewClient(cf.ConfigClient.Config.TwoCaptcha.Token)
}
