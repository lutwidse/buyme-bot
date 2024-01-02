package internal

import (
	"buyme-bot/pkg/captcha"
	"buyme-bot/pkg/config"
	"buyme-bot/pkg/proxy"
	"buyme-bot/pkg/util"
	"log"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type ClientFactory struct {
	Logger  *zap.SugaredLogger
	Config  *config.Config
	Proxy   *proxy.Proxy
	Util    *util.Util
	Captcha *captcha.Client
	Discord *discordgo.Session
}

func NewClientFactory(debug bool) *ClientFactory {
	var logger *zap.Logger
	var err error

	if debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	sugar := logger.Sugar()

	cf := &ClientFactory{Logger: sugar}
	cf.Init()
	return cf
}

func (cf *ClientFactory) Init() {
	cf.Config = cf.newConfig()
	cf.Proxy = cf.newProxy()
	cf.Util = cf.newUtil()
	cf.Captcha = cf.newCaptcha()
	cf.Discord = cf.newDiscord()
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

func (cf *ClientFactory) newDiscord() *discordgo.Session {
	discord, _ := discordgo.New("Bot " + cf.Config.ConfigData.Discord.Token)

	return discord
}
