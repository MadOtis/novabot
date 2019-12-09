package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
)

const token string = "NjUzMzk1OTY4ODEyNzc3NDcz.Xe2Z2w.4kFMzinZ8C3iBeB-D2qKByRhQrY"

var botID string
var DB *sqlx.DB

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := dg.User("@me")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	dg.AddHandler(messageHandler)
	botID = u.ID

	err = dg.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Connecting to database...")
	DB = sqlx.MustConnect("mysql", "novabot:novabot@tcp(localhost:3306)/novablack?charset=utf8&parseTime=true")
	defer DB.Close()

	fmt.Println("Bot is running!")

	<-make(chan struct{})
	return
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	var users []User

	if m.Author.ID == botID {
		return
	}

	if m.Content == "!ships" {
		fmt.Println(m.Author.Username)

		err := DB.Select(&users, "select * from users")
		if err != nil {
			return
		}
		fmt.Println(len(users))
		_, _ = s.ChannelMessageSend(m.ChannelID, users[0].firstName)
	}
	fmt.Println(m.Content)
}
