package ws

import (
	"discord-chat-bot/main/pkg/structs"
	"discord-chat-bot/main/pkg/ws/responsereader"
	"discord-chat-bot/main/pkg/ws/responsewriter"
	"fmt"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

func Run() {
	openaiKey := os.Getenv("OPENAI_API_KEY")
	discrodKey := os.Getenv("DISCORD_API_KEY")

	u := url.URL{Scheme: "wss", Host: "gateway.discord.gg", Path: "/"}
	fmt.Println(u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
	}
	ch := make(chan int)
	msg := make(chan structs.Data, 5)
	go responsereader.Read(conn, discrodKey, ch, msg)
	go responsewriter.Write(conn, discrodKey, openaiKey, ch, msg)
}
