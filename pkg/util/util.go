package util

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
	"go.uber.org/zap"
)

type Util struct {
	Logger *zap.SugaredLogger
	result map[string]interface{}
}

// CheckRecaptcha returns recaptcha parameters and true if recaptcha is present. otherwise false. We consider SEC uses only Cloudflare.
func (u *Util) CheckRecaptcha(url string) (map[string]interface{}, error) {
	pw, err := playwright.Run()
	if err != nil {
		u.Logger.Errorf("could not start playwright: %v", err)
		return make(map[string]interface{}), err
	}

	if err != nil {
		u.Logger.Errorf("Failed to get current directory: %v", err)
		return make(map[string]interface{}), err
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless:          playwright.Bool(false),
		IgnoreDefaultArgs: []string{"--enable-automation"},
		Args:              []string{"--start-maximized"},
	})
	if err != nil {
		u.Logger.Errorf("could not launch browser: %v", err)
		return make(map[string]interface{}), err
	}
	page, err := browser.NewPage()
	if err != nil {
		u.Logger.Errorf("could not create page: %v", err)
		return make(map[string]interface{}), err
	}
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	}
	err = page.SetExtraHTTPHeaders(headers)
	if err != nil {
		u.Logger.Errorf("could not set user agent: %v", err)
		return make(map[string]interface{}), err
	}

	currentDir, err := os.Getwd()
	scriptPath := filepath.Join(currentDir, "..", "pkg", "util", "turnstile.js")
	script := playwright.Script{Path: &scriptPath}
	err = page.AddInitScript(script)
	if err != nil {
		u.Logger.Errorf("could not inject script: %v", err)
		return make(map[string]interface{}), err
	}

	scriptContent := "Object.defineProperty(navigator, 'webdriver', {get: () => undefined})"
	script = playwright.Script{Content: &scriptContent}
	err = page.AddInitScript(script)

	scriptContent = "const originalQuery = window.navigator.permissions.query; window.navigator.permissions.query = (parameters) => (parameters.name === 'notifications' ? Promise.resolve({ state: Notification.permission }) : originalQuery(parameters));"
	script = playwright.Script{Content: &scriptContent}
	err = page.AddInitScript(script)

	scriptContent = "Object.defineProperty(navigator, 'plugins', { get: function () { return [2]; }, });"
	script = playwright.Script{Content: &scriptContent}
	err = page.AddInitScript(script)

	scriptContent = "Object.defineProperty(navigator, 'languages', { get: function () { return ['en-US', 'en']; }, });"
	script = playwright.Script{Content: &scriptContent}
	err = page.AddInitScript(script)

	var wg sync.WaitGroup
	wg.Add(1)

	page.OnConsole(func(msg playwright.ConsoleMessage) {
		if strings.Contains(msg.Text(), "sitekey") {
			var result map[string]interface{}
			err := json.Unmarshal([]byte(msg.Text()), &result)
			if err != nil {
				u.Logger.Errorf("Error unmarshalling JSON: %v", err)
				return
			}

			u.result = result
			wg.Done()
		}
	})

	_, err = page.Goto(url)
	if err != nil {
		u.Logger.Errorf("could not navigate to website: %v", err)
		return make(map[string]interface{}), err
	}

	timeout := time.After(10 * time.Second)
	done := make(chan bool)

	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-timeout:
		u.Logger.Info("Timeout reached while waiting for console message")
		return make(map[string]interface{}), err
	case <-done:
	}

	if err = browser.Close(); err != nil {
		u.Logger.Infof("could not close browser: %v", err)
		return make(map[string]interface{}), err
	}
	pw.Stop()

	return u.result, err
}
