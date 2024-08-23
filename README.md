# Discord Gateway Connector

Minimal code implementation with the purpose of connecting to Discord's gateway API, monitoring the text channel and having the bot respond with ChatGPT's chat bot.

# Installation

### Dependancies

Go Version 1.22.5

[Gorilla Websocket](https://pkg.go.dev/github.com/gorilla/websocket@v1.5.3)

## Environment Variables

To run this project, you will need to create and add the following environment variables to your .env file

`TOKEN`

Where `TOKEN` is your bot token provided to you by discord

## Run Locally

Clone the project

```bash
  git clone https://github.com/cgyee/discord-chat-bot
```

Go to the project directory

```bash
  cd discord-chat-bot
```

Install dependencies

```bash
  go get github.com/gorilla/websocket
```

Start the application

```bash
  go run main.go
```
