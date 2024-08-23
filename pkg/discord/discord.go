package discord

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func GetGateway(discordKey string) string {
	client := &http.Client{}
	u := url.URL{Scheme: "http", Host: "discord.com", Path: "/api/v10/gateway/bot", RawPath: "/api/v10/gateway/bot"}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		fmt.Println("gate get error", err)
	}

	req.Header.Set("Authorization", "Bot "+discordKey)
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	fmt.Println(res.Status)
	bBytes, _ := io.ReadAll(res.Body)
	jsonBody := map[string]interface{}{}
	if err = json.Unmarshal(bBytes, &jsonBody); err != nil {
		fmt.Println(err)
	}

	gatewayUrl := jsonBody["url"]
	if gatewayUrl == nil {
		return ""
	}
	fmt.Println(gatewayUrl)
	return gatewayUrl.(string)

}

func SendMessage(token string, message string, channel string, messageRef string) {
	if message != "" || channel != "" {
		client := &http.Client{}
		u := url.URL{
			Scheme:  "https",
			Host:    "discord.com",
			Path:    "/api/v10/channels/" + channel + "/messages",
			RawPath: "/api/v10/channels/" + channel + "/messages?",
		}
		data := url.Values{
			"content": {message},
		}
		req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(data.Encode()))
		if err != nil {
			fmt.Println(err)
		}

		req.Header = http.Header{}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bot "+token)
		req.Header.Set("User-Agent", "DiscordBot (https://github.com/cgyee/discord-chat-bot, 0.1.0)")
		req.Header.Set("Accept", "application/json")
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()

		bBytes, _ := io.ReadAll(res.Body)
		jsonBody := map[string]interface{}{}
		if err = json.Unmarshal(bBytes, &jsonBody); err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(bBytes))

	}
}

type Bot struct {
	Name string `json:name,omitempty`
}

func RequestBotInfo(token string, appId string) string {
	client := &http.Client{}
	u := url.URL{
		Scheme:  "https",
		Host:    "discord.com",
		Path:    "/api/v10/applications/" + appId,
		RawPath: "/api/v10/applications/" + appId,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header = http.Header{}
	req.Header.Set("Authorization", "Bot "+token)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	fmt.Println(res.Status)
	bBytes, _ := io.ReadAll(res.Body)
	botInfo := map[string]interface{}{}
	fmt.Println(string(bBytes))
	json.Unmarshal(bBytes, &botInfo)

	return botInfo["name"].(string)

}
