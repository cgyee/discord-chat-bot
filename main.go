package main

import (
	"discord-chat-bot/main/pkg/chatrequest"
	"discord-chat-bot/main/pkg/discordmessage"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type DiscordData struct {
	message       string
	channel       string
	guild         string
	hasMessageRef bool
	messageRef    string
}

type Bot struct {
	name string
}

var ready bool = false
var ch = make(chan int, 5)
var msg = make(chan DiscordData, 5)
var resumeGatewayUrl string
var sessionId string
var seq int
var appId string
var openaiKey string
var discrodApiKey string

func heartbeat(heartbeat_interval int) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(heartbeat_interval))
	for {
		tm := <-ticker.C
		fmt.Println("💕 Heartbeat 💕 at:", tm)
		ch <- 1
	}
}
func responseWriter(ws *websocket.Conn, token string) {
	for {
		opCode := <-ch
		switch opCode {
		case 0:
			m := <-msg
			message, channel := m.message, m.channel
			go func() {
				response, err := chatrequest.Request(openaiKey, message)
				if err != nil {
					fmt.Println("chatRequest err", err)

				}
				discordmessage.Request(discrodApiKey, response, channel, "")
			}()
		case 1:
			ws.WriteJSON(map[string]interface{}{
				"op": 1, "d": nil,
			})
		case 6:
			ws.WriteJSON((map[string]interface{}{
				"op": 6,
				"d": map[string]interface{}{
					"token":     token,
					"sessionId": sessionId,
					"seq":       seq,
				},
			}))
		case 7:
			ws.WriteJSON(map[string]interface{}{
				"op": 7,
				"d":  nil,
			})
		case 10:
			ws.WriteJSON(map[string]interface{}{
				"op": 2,
				"d": map[string]interface{}{
					"token":   token,
					"intents": 512,
					"properties": map[string]interface{}{
						"os":      "macos",
						"browser": "chrome",
						"device":  "macbook air",
					},
				},
			})
		}

	}

}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Load", err)
	}
	openaiKey = os.Getenv("OPENAI_API_KEY")
	discrodApiKey = os.Getenv("DISCORD_API_KEY")
	appId := os.Getenv("APP_ID")
	bot := discordmessage.RequestBotInfo(discrodApiKey, appId)
	fmt.Println(discrodApiKey)
	u := url.URL{Scheme: "wss", Host: "gateway.discord.gg", Path: "/"}
	fmt.Println(u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	if err != nil {
		fmt.Println("Socket", err)
	}
	go responseWriter(conn, discrodApiKey)
	for {
		_, content, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read fail ", err)
			return
		}

		data := map[string]interface{}{}
		_ = json.Unmarshal(content, &data)
		var opCode float64
		if data["op"] != nil {
			opCode = data["op"].(float64)
		}
		seq = int(opCode)
		fmt.Println("Dispacth code:", opCode)

		switch opCode {
		case 10:
			if !ready {
				d := data["d"].(map[string]interface{})
				heartbeat_interval := int(d["heartbeat_interval"].(float64))
				ch <- 10
				go heartbeat(heartbeat_interval)
			}
			ready = true
		case 6:
			ch <- 6
			u := url.URL{Scheme: "wss", Host: resumeGatewayUrl, Path: "/"}
			fmt.Println(u.String())
			conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				fmt.Println(err)
			}
		case 7:
			ch <- 7
		}
		status := data["t"]
		switch status {

		case "READY":
			fmt.Println("READY")
			d := data["d"].(map[string]interface{})
			sessionId = d["session_id"].(string)
			resumeGatewayUrl = d["resume_gateway_url"].(string)

		case "MESSAGE_CREATE":
			d := data["d"].(map[string]interface{})
			bytes, _ := json.Marshal(d)
			fmt.Println(string(bytes))
			author := d["author"].(map[string]interface{})
			username := author["username"].(string)
			if bot != username {
				content := d["content"].(string)
				guild := d["guild_id"].(string)
				channel := d["channel_id"].(string)
				messageRef := d["id"].(string)
				fmt.Println(username, "said", content)
				ch <- 0
				message := DiscordData{
					message:    content,
					channel:    channel,
					guild:      guild,
					messageRef: messageRef,
				}
				msg <- message
			}
		default:
			if status != nil {
				fmt.Println("Status:", status)
			}
		}

	}
}
