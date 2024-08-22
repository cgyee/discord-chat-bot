package responsereader

import (
	"discord-chat-bot/main/pkg/discordmessage"
	"discord-chat-bot/main/pkg/structs"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func Read(ws *websocket.Conn, discordKey string, ch chan int, msg chan structs.Data) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Load", err)
	}
	var sessionId string
	// var seq int
	var resumeGatewayUrl string
	ready := false
	fmt.Println("Read")
	appId := os.Getenv("APP_ID")

	bot := discordmessage.RequestBotInfo(discordKey, appId)
	for {
		_, content, err := ws.ReadMessage()
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
		// seq = int(opCode)
		fmt.Println("Dispacth code:", opCode)

		switch opCode {
		case 10:
			if !ready {
				d := data["d"].(map[string]interface{})
				heartbeat_interval := int(d["heartbeat_interval"].(float64))
				ch <- 10
				go heartbeat(heartbeat_interval, ch)
			}
			ready = true
		case 6:
			ch <- 6
			u := url.URL{Scheme: "wss", Host: resumeGatewayUrl, Path: "/"}
			fmt.Println(u.String())
			ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
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
				message := structs.Data{
					Message:    content,
					Channel:    channel,
					Guild:      guild,
					MessageRef: messageRef,
					SessionId:  sessionId,
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

func heartbeat(heartbeat_interval int, ch chan<- int) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(heartbeat_interval))
	for {
		tm := <-ticker.C
		fmt.Println("💕 Heartbeat 💕 at:", tm)
		ch <- 1
	}
}
