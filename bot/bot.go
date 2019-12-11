package bot

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	// needed for the database/sql package
	_ "github.com/go-sql-driver/mysql"
)

const userMarker = "@"

var botID string
var goBot *discordgo.Session

// DB is a global database pool object
var DB *sql.DB

var botPrefix string

// Start the bot running
func Start(prefix string, botToken string, sqlUser string, sqlPass string, sqlHost string, sqlPort string, sqlDatabase string) {
	botPrefix = prefix
	goBot, err := discordgo.New("Bot " + botToken)
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
	sqlConnectString := sqlUser + ":" + sqlPass + "@tcp(" + sqlHost + ":" + sqlPort + ")/" + sqlDatabase + "?charset=utf8&parseTime=true"
	DB, err = sql.Open("mysql", sqlConnectString)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Bot is running!")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Content, botPrefix) {
		command := strings.TrimPrefix(m.Content, botPrefix)
		if m.Author.ID == botID {
			return
		}

		if strings.HasPrefix(command, "help") {
			resultString := buildHelpMessage()
			_, _ = s.ChannelMessageSend(m.ChannelID, resultString)
		}

		if strings.HasPrefix(command, "bio") {
			resultString := "Fucker, I told you this isn't implemented yet!"
			_, _ = s.ChannelMessageSend(m.ChannelID, resultString)
		}

		if strings.HasPrefix(command, "ships") {
			userSpecified := strings.TrimSpace(strings.TrimPrefix(command, "ships"))
			if len(userSpecified) > 0 {
				if strings.HasPrefix(userSpecified, userMarker) {
					userSpecified = strings.TrimPrefix(userSpecified, userMarker)
				}
			} else {
				userSpecified = m.Author.Username
			}

			resultMessage := buildShipList(userSpecified)
			_, _ = s.ChannelMessageSend(m.ChannelID, resultMessage)
		}
	}
}

func buildHelpMessage() string {
	resultString := "HELP TOPICS...\n"
	resultString = resultString + "============================\n"
	resultString = resultString + "!help - this content, duh!\n"
	resultString = resultString + "  NOTE: All commands accept an optional [handle] argument - if specified, I will return the requested\n"
	resultString = resultString + "  info for that user, if I can find that user, so make sure his or her handle is correct\n\n"
	resultString = resultString + "!ships [handle]\n"
	resultString = resultString + "!bio [handle] - not implemented, yet... stay tuned!\n"
	return resultString
}

func buildShipList(userSpecified string) string {
	dbrows, err := DB.Query("select s.name, s.nickname, s.crewsize from ships s, ownedships os, users u where s.id = os.shipid and os.userid = u.userid and u.handle = ?", userSpecified)
	if err != nil {
		panic(err.Error())
	}
	defer dbrows.Close()

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
	return resultMessage
}
