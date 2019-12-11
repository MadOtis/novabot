package main

import (
	"os"

	"github.com/madotis/novabot/bot"
)

func main() {

	botPrefix := os.Getenv("BOT_PREFIX")
	botToken := os.Getenv("BOT_TOKEN")
	sqlUser := os.Getenv("SQL_USER")
	sqlPass := os.Getenv("SQL_PASS")
	sqlHost := os.Getenv("SQL_HOST")
	sqlPort := os.Getenv("SQL_PORT")
	sqlDatabase := os.Getenv("SQL_DATABASE")

	bot.Start(botPrefix, botToken, sqlUser, sqlPass, sqlHost, sqlPort, sqlDatabase)

	<-make(chan struct{})
	return
}
