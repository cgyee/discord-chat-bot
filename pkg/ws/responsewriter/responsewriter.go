package responsewriter

import (
	"discord-chat-bot/main/pkg/chatrequest"
	"discord-chat-bot/main/pkg/discordmessage"
	"discord-chat-bot/main/pkg/structs"
	"fmt"

	"github.com/gorilla/websocket"
)

func Write(ws *websocket.Conn, discordKey string, openaiKey string, ch chan int, msg <-chan structs.Data) {
	var sessionId string
	var seq int
	for {
		opCode := <-ch
		switch opCode {
		case 0:
			m := <-msg
			message, channel := m.Message, m.Channel
			sessionId, seq = m.SessionId, m.Seq
			go func() {
				response, err := chatrequest.Request(openaiKey, message)
				if err != nil {
					fmt.Println("chatRequest err", err)

				}
				discordmessage.Request(discordKey, response, channel, "")
			}()
		case 1:
			ws.WriteJSON(map[string]interface{}{
				"op": 1, "d": nil,
			})
		case 6:
			ws.WriteJSON((map[string]interface{}{

				"op": 6,
				"d": map[string]interface{}{
					"token":     discordKey,
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
					"token":   discordKey,
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
