package main

import (
	client "buyme-bot/internal"
	"buyme-bot/pkg/captcha"
	"log"

	"go.uber.org/zap"
)

// TODO: Separate the code to internal.
// TODO: Change the temporary code.

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	sugar := logger.Sugar()

	buymeClient := client.NewClientFactory(sugar)

	result, err := buymeClient.UtilClient.CheckRecaptcha("https://nopecha.com/demo/cloudflare")
	if err != nil {
		sugar.Errorf("Recaptcha check failed: %v", err)
		return
	}

	if len(result) > 3 {
		sugar.Infof("Captcha Type: Challenge Page")

		sitekey, _ := result["sitekey"].(string)
		pageurl, _ := result["pageurl"].(string)
		data, _ := result["data"].(string)
		pagedata, _ := result["pagedata"].(string)
		action, _ := result["action"].(string)
		userAgent, _ := result["userAgent"].(string)
		sugar.Infof("Site Key: %s, Page URL: %s, Data: %s, Page Data: %s, Action: %s, User Agent: %s",
			sitekey, pageurl, data, pagedata, action, userAgent)

		cap := captcha.CloudflareTurnstile{
			SiteKey:   sitekey,
			Url:       pageurl,
			Data:      data,
			PageData:  pagedata,
			Action:    action,
			UserAgent: userAgent,
		}

		sessionProxy := buymeClient.ProxyClient.GetSessionProxyURL("Japan")

		req := cap.ToRequest()
		req.SetProxy("HTTP", sessionProxy)

		code, err := buymeClient.CaptchaClient.Send(req)
		if err != nil {
			sugar.Errorf("Failed to send captcha: %v", err)
			return
		}

		captchaResult, err := buymeClient.CaptchaClient.WaitForResult(code, 60, 15)
		if err != nil {
			sugar.Errorf("Failed to wait for captcha result: %v", err)
			return
		}
		sugar.Infof("Captcha Result: %s", captchaResult)
	} else if len(result) < 3 {
		sugar.Infof("Captcha Type: Standalone")

		sitekey, _ := result["sitekey"].(string)
		pageurl, _ := result["pageurl"].(string)
		userAgent, _ := result["userAgent"].(string)
		sugar.Infof("Site Key: %s, Page URL: %s, User Agent: %s",
			sitekey, pageurl, userAgent)

		cap := captcha.CloudflareTurnstile{
			SiteKey: sitekey,
			Url:     pageurl,
		}

		sessionProxy := buymeClient.ProxyClient.GetSessionProxyURL("Japan")

		req := cap.ToRequest()
		req.SetProxy("HTTP", sessionProxy)

		code, err := buymeClient.CaptchaClient.Send(req)
		if err != nil {
			sugar.Errorf("Failed to send captcha: %v", err)
			return
		}

		captchaResult, err := buymeClient.CaptchaClient.WaitForResult(code, 60, 15)
		if err != nil {
			sugar.Errorf("Failed to wait for captcha result: %v", err)
			return
		}
		sugar.Infof("Captcha Result: %s", captchaResult)
	}
}
