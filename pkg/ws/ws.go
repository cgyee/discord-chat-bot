package ws

import (
	"discord-chat-bot/main/pkg/discord"
	"discord-chat-bot/main/pkg/structs"
	"discord-chat-bot/main/pkg/ws/responsereader"
	"discord-chat-bot/main/pkg/ws/responsewriter"
	"fmt"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func Run() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Load", err)
	}
	openaiKey := os.Getenv("OPENAI_API_KEY")
	discrodKey := os.Getenv("DISCORD_API_KEY")

	u := discord.GetGateway(discrodKey)
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	ch := make(chan int)
	msg := make(chan structs.Data, 5)
	invalid := make(chan error)
	go responsereader.Read(conn, discrodKey, ch, msg, invalid)
	go responsewriter.Write(conn, discrodKey, openaiKey, ch, msg, invalid)
	for {
		err := <-invalid
		if err != nil {
			break
		}
	}

}
