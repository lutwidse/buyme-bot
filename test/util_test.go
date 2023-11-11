package util_test

import (
	"buyme-bot/pkg/util"
	"testing"

	"go.uber.org/zap"
)

// TODO: Subdivision of test code. need to learn more.

func TestCloudFlareFailCase(t *testing.T) {
	sugar := zap.NewExample().Sugar()

	utilClient := &util.Util{
		Logger: sugar,
	}
	result, err := utilClient.CheckCloudFlareRecaptcha("https://nopecha.com/demo/turnstile")
	if err != nil {
		sugar.Errorf("Recaptcha check failed: %v", err)
	} else {
		if len(result) > 3 {
			sugar.Infof("Captcha Type: Challenge Page")

			sitekey, _ := result["sitekey"].(string)
			pageurl, _ := result["pageurl"].(string)
			data, _ := result["data"].(string)
			pagedata, _ := result["pagedata"].(string)
			action, _ := result["action"].(string)
			useragent, _ := result["userAgent"].(string)
			sugar.Infof("Site Key: %s, Page URL: %s, Data: %s, Page Data: %s, Action: %s, User Agent: %s",
				sitekey, pageurl, data, pagedata, action, useragent)
		} else {
			sugar.Infof("Captcha Type: Standalone")

			sitekey, _ := result["sitekey"].(string)
			pageurl, _ := result["pageurl"].(string)
			useragent, _ := result["userAgent"].(string)
			sugar.Infof("Site Key: %s, Page URL: %s, User Agent: %s",
				sitekey, pageurl, useragent)
		}
	}
}

func TestStandalone(t *testing.T) {
	sugar := zap.NewExample().Sugar()

	utilClient := &util.Util{
		Logger: sugar,
	}
	result, err := utilClient.CheckCloudFlareRecaptcha("https://nopecha.com/demo/turnstile")
	if err != nil {
		sugar.Errorf("Recaptcha check failed: %v", err)
	} else {
		if len(result) > 3 {
			sugar.Infof("Captcha Type: Challenge Page")

			sitekey, _ := result["sitekey"].(string)
			pageurl, _ := result["pageurl"].(string)
			data, _ := result["data"].(string)
			pagedata, _ := result["pagedata"].(string)
			action, _ := result["action"].(string)
			useragent, _ := result["userAgent"].(string)
			sugar.Infof("Site Key: %s, Page URL: %s, Data: %s, Page Data: %s, Action: %s, User Agent: %s",
				sitekey, pageurl, data, pagedata, action, useragent)
		} else {
			sugar.Infof("Captcha Type: Standalone")

			sitekey, _ := result["sitekey"].(string)
			pageurl, _ := result["pageurl"].(string)
			useragent, _ := result["userAgent"].(string)
			sugar.Infof("Site Key: %s, Page URL: %s, User Agent: %s",
				sitekey, pageurl, useragent)
		}
	}
}

func TestChallengePage(t *testing.T) {
	sugar := zap.NewExample().Sugar()

	utilClient := &util.Util{
		Logger: sugar,
	}
	result, err := utilClient.CheckCloudFlareRecaptcha("https://nopecha.com/demo/cloudflare")
	if err != nil {
		sugar.Errorf("Recaptcha check failed: %v", err)
	} else {
		if len(result) > 3 {
			sugar.Infof("Captcha Type: Challenge Page")

			sitekey, _ := result["sitekey"].(string)
			pageurl, _ := result["pageurl"].(string)
			data, _ := result["data"].(string)
			pagedata, _ := result["pagedata"].(string)
			action, _ := result["action"].(string)
			useragent, _ := result["userAgent"].(string)
			sugar.Infof("Site Key: %s, Page URL: %s, Data: %s, Page Data: %s, Action: %s, User Agent: %s",
				sitekey, pageurl, data, pagedata, action, useragent)
		} else {
			sugar.Infof("Captcha Type: Standalone")

			sitekey, _ := result["sitekey"].(string)
			pageurl, _ := result["pageurl"].(string)
			useragent, _ := result["userAgent"].(string)
			sugar.Infof("Site Key: %s, Page URL: %s, User Agent: %s",
				sitekey, pageurl, useragent)
		}
	}
}
