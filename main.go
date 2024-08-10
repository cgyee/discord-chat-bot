package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var ready bool = false
var ch = make(chan int, 5)
var msg = make(chan string, 5)
var resumeGatewayUrl string
var sessionId string
var seq int
var encoder *form.Encoder

func chatRequest(host string, token string) {
	client := &http.Client{}
	u := url.URL{Scheme: "https", Host: host, Path: "/v1/chat/completions"}
	fmt.Println(u.String())

	payload := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]interface{}{
			{"role": "system",
				"content": "You are a helpful assistant."},
			{"role": "user",
				"content": "Hello!"},
		},
	}
	pBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(pBytes))
	fmt.Println(string(pBytes))

	if err != nil {
		fmt.Println("Request failed", err)
	}
	req.Header = http.Header{}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	fmt.Println(req.Header)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Do req failed", err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	fmt.Println(res.Status)
	fmt.Println(res.StatusCode)
	fmt.Println(string(body))

}

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

	openaiKey := os.Getenv("OPENAI_API_KEY")
	openaiApi := os.Getenv("OPENAI_API")
	chatRequest(openaiApi, openaiKey)
	// tkn := os.Getenv("TOKEN")
	// u := url.URL{Scheme: "wss", Host: "gateway.discord.gg", Path: "/"}
	// fmt.Println(u.String())
	// conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer conn.Close()
	// if err != nil {
	// 	fmt.Println("Socket", err)
	// }
	// go responseWriter(conn, tkn)
	// for {
	// 	_, content, err := conn.ReadMessage()
	// 	if err != nil {
	// 		fmt.Println("Read fail ", err)
	// 		return
	// 	}

	// 	data := map[string]interface{}{}
	// 	_ = json.Unmarshal(content, &data)
	// 	var opCode float64
	// 	if data["op"] != nil {
	// 		opCode = data["op"].(float64)
	// 	}
	// 	seq = int(opCode)
	// 	fmt.Println("Dispacth code:", opCode)
	// 	switch opCode {
	// 	case 10:
	// 		if !ready {
	// 			d := data["d"].(map[string]interface{})
	// 			heartbeat_interval := int(d["heartbeat_interval"].(float64))
	// 			ch <- 10
	// 			go heartbeat(heartbeat_interval)
	// 		}
	// 		ready = true
	// 	case 6:
	// 		ch <- 6
	// 		u := url.URL{Scheme: "wss", Host: resumeGatewayUrl, Path: "/"}
	// 		fmt.Println(u.String())
	// 		conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 	case 7:
	// 		ch <- 7

	// 	}
	// 	status := data["t"]
	// 	switch status {

	// 	case "READY":
	// 		fmt.Println("READY")
	// 		d := data["d"].(map[string]interface{})
	// 		sessionId = d["session_id"].(string)
	// 		resumeGatewayUrl = d["resume_gateway_url"].(string)

	// 	case "MESSAGE_CREATE":
	// 		d := data["d"].(map[string]interface{})
	// 		content := d["content"].(string)
	// 		author := d["author"].(map[string]interface{})
	// 		username := author["username"].(string)
	// 		fmt.Println(username, "said", content)
	// 		fmt.Println(d)
	// 	default:
	// 		if status != nil {
	// 			fmt.Println("Status:", status)
	// 		}
	// 	}

	// }
}
