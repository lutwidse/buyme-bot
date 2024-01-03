package main

import (
	client "buyme-bot/internal"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// TODO: Separate code to internal.
// TODO: Change temporary code.

type Source struct {
	DisplayNames []string `json:"display_names"`
	RootForm     string   `json:"root_form"`
	Form         string   `json:"form"`
	FileDate     string   `json:"file_date"`
	Ciks         []string `json:"ciks"`
	BizLocations []string `json:"biz_locations"`
	FileNum      []string `json:"file_num"`
	FilmNum      []string `json:"film_num"`
}

type Item struct {
	ID     string `json:"_id"`
	Source Source `json:"_source"`
}

type Hits struct {
	Hits []Item `json:"hits"`
}

type Response struct {
	Hits Hits `json:"hits"`
}

func main() {
	debug := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()
	buymeClient := client.NewClientFactory(*debug)
	monitorEdgar(buymeClient)
}

func monitorEdgar(client *client.ClientFactory) {
	processedItems := make(map[string]bool)

	loop := func(startdt string, enddt string) {
		urlStr := fmt.Sprintf("https://efts.sec.gov/LATEST/search-index?q=BTC&dateRange=custom&startdt=%s&enddt=%s", startdt, enddt)
		req, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			client.Logger.Errorf("Failed to create request: %v", err)
			return
		}
		req.Header.Set("User-Agent", client.Config.ConfigData.Setting.UserAgent)

		proxyURL, err := url.Parse("http://" + client.Config.ConfigData.Proxy.Host + ":" + client.Config.ConfigData.Proxy.Port)
		if err != nil {
			client.Logger.Errorf("Failed to parse proxy URL: %v", err)
			return
		}

		proxyURL.User = url.UserPassword(client.Config.ConfigData.Proxy.User, client.Proxy.GetSessionProxy("Japan"))

		httpClient := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			client.Logger.Errorf("Failed to send request: %v", err)
			return
		}

		// Their API returns error most of the time.
		if resp.StatusCode != 200 {
			client.Logger.Debugf("Failed to get response: %v", resp.StatusCode)
			return
		}

		var result Response

		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			client.Logger.Errorf("Failed to decode response body: %v", err)
			return
		}

		defer resp.Body.Close()

		/*
			Initialize processedItems only once to prevent notifications on the first run.
			This is called only once. even if all the loops within an hour complete, it will NEVER be called next time.
			Their API weirdly returns the previous items that are not in the date range we passed.
			But it's okay, we're checking if start date is before file date.
			We don't want to miss any update.
		*/
		if len(processedItems) == 0 {
			for _, item := range result.Hits.Hits {
				processedItems[item.ID] = true
				client.Logger.Debugf("Added to processed (init): %s", item.ID)
			}
		}

		for _, item := range result.Hits.Hits {
			_source := item.Source
			for j := range _source.DisplayNames {
				client.Logger.Debugf("ID: %s", item.ID)

				if processedItems[item.ID] {
					client.Logger.Debugf("Already processed: %s", item.ID)
					continue
				}

				_ids := strings.Split(item.ID, ":")

				rootForm := _source.RootForm
				form := _source.Form
				fileDate := _source.FileDate
				displayName := _source.DisplayNames[j]
				cik := _source.Ciks[j]
				bizLocation := _source.BizLocations[j]

				idParamWhat := strings.ReplaceAll(_ids[0], "-", "")
				idParamName := _ids[1]
				rootFormURL := fmt.Sprintf("https://www.sec.gov/Archives/edgar/data/%s/%s/%s", cik[3:], idParamWhat, idParamName)

				fileDateTime, err := time.Parse("2006-01-02", fileDate)
				if err != nil {
					client.Logger.Errorf("Failed to parse date: %v", err)
					continue
				}
				startdtTime, err := time.Parse("2006-01-02", startdt)
				if err != nil {
					client.Logger.Errorf("Failed to parse date: %v", err)
					continue
				}

				if !startdtTime.After(fileDateTime) {
					client.Logger.Debugf("File date is not after start date: %v", fileDateTime)
					continue
				}

				if _source.RootForm != "S-1" || _source.RootForm != "S-3" {
					client.Logger.Debugf("Unexpected root form : %v", _source.RootForm)
					continue
				}

				for i := 0; i < len(_source.FileNum); i++ {
					fileNum := _source.FileNum[i]
					filmNum := _source.FilmNum[i]
					fileNumURL := fmt.Sprintf("https://www.sec.gov/cgi-bin/browse-edgar/?filenum=%s&action=getcompany", fileNum)

					client.Logger.Debugf("rootForm: %s, form: %s, fileDate: %s, displayName: %s, cik: %s, bizLocation: %s, fileNum: %s, filmNum: %s, rootFormURL: %s, fileNumURL: %s", rootForm, form, fileDate, displayName, cik, bizLocation, fileNum, filmNum, rootFormURL, fileNumURL)

					embed := &discordgo.MessageEmbed{
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  "Form & File",
								Value: fmt.Sprintf("[%s](%s)", form, rootFormURL),
							},
							{
								Name:  "Filed",
								Value: fileDate,
							},
							{
								Name:  "Filing entity/person",
								Value: displayName,
							},
							{
								Name:  "CIKS",
								Value: cik,
							},
							{
								Name:  "Located",
								Value: bizLocation,
							},
							{
								Name:  "File number",
								Value: fmt.Sprintf("[%s](%s)", fileNum, fileNumURL),
							},
							{
								Name:  "Film number",
								Value: filmNum,
							},
						},
					}

					mentionEmbed := &discordgo.MessageEmbed{
						Title: "EDGAR FILE LISTED 📔 ✅",
						Color: 0x5865F2,
					}

					_, err := client.Discord.ChannelMessageSendEmbed(client.Config.ConfigData.Discord.ChannelID, mentionEmbed)
					if err != nil {
						client.Logger.Errorf("Error sending mentionEmbed message: ", err)
					}
					_, err = client.Discord.ChannelMessageSend(client.Config.ConfigData.Discord.ChannelID, fmt.Sprintf("<@&%s>", client.Config.ConfigData.Discord.RoleID))
					if err != nil {
						client.Logger.Errorf("Error sending mention message: ", err)
					}

					// We use ChannelMessageSendEmbed instead of Embeds because both functions handle Embeds one by one.
					_, err = client.Discord.ChannelMessageSendEmbed(client.Config.ConfigData.Discord.ChannelID, embed)
					if err != nil {
						client.Logger.Errorf("Error sending embed: ", err)
					}
					processedItems[item.ID] = true
					client.Logger.Debugf("Added to processed: %s", item.ID)
				}
			}
		}
	}

	for {
		now := time.Now().UTC()
		startdt := now.AddDate(0, 0, -1).Format("2006-01-02")
		enddt := now.AddDate(0, 0, 1).Format("2006-01-02")

		for i := 0; i < 60*60; i++ {
			loop(startdt, enddt)
			time.Sleep(1 * time.Second)
		}
	}
}
