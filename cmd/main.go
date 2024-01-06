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
	DisplayNames    []string `json:"display_names"`
	RootForm        string   `json:"root_form"`
	Form            string   `json:"form"`
	FileDate        string   `json:"file_date"`
	Ciks            []string `json:"ciks"`
	BizLocations    []string `json:"biz_locations"`
	FileNum         []string `json:"file_num"`
	FilmNum         []string `json:"film_num"`
	FileType        string   `json:"file_type"`
	FileDescription string   `json:"file_description"`
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

	validRootForms := map[string]bool{
		"S-1":    true,
		"S-3":    true,
		"8-A12B": true,
	}

	searchTerms := []string{"BTC", "Bitcoin"}

	loop := func(startdt string, enddt string) {
		queryParams := []string{
			// API accepts sql query
			"q=" + url.QueryEscape(strings.Join(searchTerms, " OR ")),
			"dateRange=custom",
			// Registration statements and prospectuses
			"category=form-cat5",
			"startdt=" + url.QueryEscape(startdt),
			"enddt=" + url.QueryEscape(enddt),
			// form-cat5
			"forms=" + url.QueryEscape("10-12B,10-12G,18-12B,20FR12B,20FR12G,40-24B2,40FR12B,40FR12G,424A,424B1,424B2,424B3,424B4,424B5,424B7,424B8,424H,425,485APOS,485BPOS,485BXT,487,497,497J,497K,8-A12B,8-A12G,AW,AW WD,DEL AM,DRS,F-1,F-10,F-10EF,F-10POS,F-3,F-3ASR,F-3D,F-3DPOS,F-3MEF,F-4,F-4 POS,F-4MEF,F-6,F-6 POS,F-6EF,F-7,F-7 POS,F-8,F-8 POS,F-80,F-80POS,F-9,F-9 POS,F-N,F-X,FWP,N-2,POS AM,POS EX,POS462B,POS462C,POSASR,RW,RW WD,S-1,S-11,S-11MEF,S-1MEF,S-20,S-3,S-3ASR,S-3D,S-3DPOS,S-3MEF,S-4,S-4 POS,S-4EF,S-4MEF,S-6,S-8,S-8 POS,S-B,S-BMEF,SF-1,SF-3,SUPPL,UNDER"),
		}

		baseURL := "https://efts.sec.gov/LATEST/search-index"
		urlStr := baseURL + "?" + strings.Join(queryParams, "&")
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

		// API returns error most of the time
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
			This is called only once, NEVER called next time.

			API returns the previous items that are not in the date range we passed on the last day of the month.
			I'm not sure if this is a form update or not. thus, we will check if the file date is after the start date.
			We don't want to miss any updates.
		*/
		if len(processedItems) == 0 {
			for _, item := range result.Hits.Hits {
				processedItems[item.ID] = true
				client.Logger.Debugf("Added to processed (init): %s", item.ID)
			}
		}

		for _, item := range result.Hits.Hits {
			_source := item.Source
			_id := item.ID

			for j := range _source.DisplayNames {
				client.Logger.Debugf("ID: %s", _id)

				if processedItems[_id] {
					client.Logger.Debugf("Already processed: %s", _id)
					continue
				}

				// Limiting the scope of rootForm in the request parameter might be better
				if !validRootForms[_source.RootForm] {
					client.Logger.Debugf("Unexpected root form : %v", _source.RootForm)
					processedItems[_id] = true
					client.Logger.Debugf("Added to processed (unexpected root form): %s", _id)
					continue
				}

				_ids := strings.Split(_id, ":")

				rootForm := _source.RootForm
				form := _source.Form
				fileDate := _source.FileDate
				fileType := _source.FileType
				fileDescription := _source.FileDescription
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

				if !fileDateTime.After(startdtTime) {
					client.Logger.Debugf("File date is not after start date: %v", fileDateTime)
					continue
				}

				for i := 0; i < len(_source.FileNum); i++ {
					fileNum := _source.FileNum[i]
					filmNum := _source.FilmNum[i]
					fileNumURL := fmt.Sprintf("https://www.sec.gov/cgi-bin/browse-edgar/?filenum=%s&action=getcompany", fileNum)

					client.Logger.Debugf("rootForm: %s, form: %s, fileDate: %s, displayName: %s, cik: %s, bizLocation: %s, fileNum: %s, filmNum: %s, fileType: %s, rootFormURL: %s, fileNumURL: %s", rootForm, form, fileDate, displayName, cik, bizLocation, fileNum, filmNum, fileType, fileDescription, rootFormURL, fileNumURL)

					embed := &discordgo.MessageEmbed{
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  "Form & File",
								Value: fmt.Sprintf("[%s - %s (%s)](%s)", form, fileType, fileDescription, rootFormURL),
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

					// We use ChannelMessageSendEmbed instead of Embeds because both functions handle Embed one by one
					_, err = client.Discord.ChannelMessageSendEmbed(client.Config.ConfigData.Discord.ChannelID, embed)
					if err != nil {
						client.Logger.Errorf("Error sending embed: ", err)
					}
					processedItems[_id] = true
					client.Logger.Debugf("Added to processed: %s", _id)
				}
			}
		}
	}

	for {
		now := time.Now().UTC()
		startdt := now.AddDate(0, 0, -1).Format("2006-01-02")
		enddt := now.AddDate(0, 0, 1).Format("2006-01-02")
		loop(startdt, enddt)
		time.Sleep(1 * time.Second)
	}
}
