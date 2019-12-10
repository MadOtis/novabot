package bot

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/madotis/novabot/config"
	"strings"
)

var botID string
var goBot *discordgo.Session

var DB *sql.DB

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	goBot.AddHandler(messageHandler)
	botID = u.ID

	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Connecting to database...")
	DB, err = sql.Open("mysql", config.MysqlUser + ":" + config.MysqlPass + "@tcp("+config.MysqlHost+":"+config.MysqlPort+")/"+config.MysqlDatabase+"?charset=utf8&parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	defer DB.Close()

	fmt.Println("Bot is running!")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Content, config.BotPrefix) {
		command := strings.TrimPrefix(m.Content, config.BotPrefix)
		if m.Author.ID == botID {
			return
		}

		if strings.HasPrefix(command, "help") {
			resultString := "HELP TOPICS...\n"
			resultString = resultString + "============================\n"
			resultString = resultString + "!help - this content, duh!\n"
			resultString = resultString + "  NOTE: All commands accept an optional [handle] argument - if specified, I will return the requested\n"
			resultString = resultString + "  info for that user, if I can find that user, so make sure his or her handle is correct\n\n"
			resultString = resultString + "!ships [handle]\n"
			resultString = resultString + "!bio [handle] - not implemented, yet... stay tuned!\n"
			_, _ = s.ChannelMessageSend(m.ChannelID, resultString)
		}

		if strings.HasPrefix(command, "bio") {
			resultString := "Fucker, I told you this isn't implemented yet!"
			_, _ = s.ChannelMessageSend(m.ChannelID, resultString)
		}

		if strings.HasPrefix(command, "ships") {
			userSpecified := strings.TrimSpace(strings.TrimPrefix(command, "ships"))
			var dbrows *sql.Rows
			var err error
			if len(userSpecified) > 0 {
				dbrows, err = DB.Query("select s.name, s.nickname, s.crewsize from ships s, ownedships os, users u where s.id = os.shipid and os.userid = u.userid and u.handle = ?", userSpecified)
				if err != nil {
					panic(err.Error())
				}
				defer dbrows.Close()
			} else {
				dbrows, err = DB.Query("select s.name, s.nickname, s.crewsize from ships s, ownedships os, users u where s.id = os.shipid and os.userid = u.userid and u.handle = ?", m.Author.Username)
				if err != nil {
					panic(err.Error())
				}
				defer dbrows.Close()
			}

			var resultMessage string
			resultMessage = "Ship | Nickname | Crew size\n"

			for dbrows.Next() {
				var shipname, shipnickname, crewsize string
				err := dbrows.Scan(&shipname, &shipnickname, &crewsize)
				if err != nil {
					panic(err.Error())
				}
				resultMessage = resultMessage + shipname + "\t" + shipnickname + "\t" + crewsize + "\n"
			}
			_, _ = s.ChannelMessageSend(m.ChannelID, resultMessage)
		}
	}
}

