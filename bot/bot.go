package bot

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	strip "github.com/grokify/html-strip-tags-go"

	// needed for the database/sql package
	_ "github.com/go-sql-driver/mysql"
)

const userMarker = "<@!"

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
			userSpecified := strings.TrimSpace(strings.TrimPrefix(command, "bio"))
			if len(userSpecified) > 0 {
				if strings.HasPrefix(userSpecified, userMarker) {
					userSpecified = strings.TrimSuffix(strings.TrimPrefix(userSpecified, userMarker), ">")
					discordUser, err := s.User(userSpecified)
					if err != nil {
						panic("User unknown")
					}
					userSpecified = discordUser.Username
				}
			} else {
				userSpecified = m.Author.Username
			}
			sendUserBio(userSpecified, m, s)
		}

		if strings.HasPrefix(command, "ships") {
			userSpecified := strings.TrimSpace(strings.TrimPrefix(command, "ships"))
			if len(userSpecified) > 0 {
				if strings.HasPrefix(userSpecified, userMarker) {
					userSpecified = strings.TrimSuffix(strings.TrimPrefix(userSpecified, userMarker), ">")
					discordUser, err := s.User(userSpecified)
					if err != nil {
						panic("User unknown")
					}
					userSpecified = discordUser.Username
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
	dbrows, err := DB.Query("select s.name, s.nickname, s.crewsize from ships s, ownedShips os, users u where s.id = os.shipid and os.userid = u.userid and u.handle = ?", userSpecified)
	if err != nil {
		panic(err.Error())
	}
	defer dbrows.Close()

	var resultMessage string
	var bFoundShips = false

	resultMessage = "Ship | Nickname | Crew size\n"

	for dbrows.Next() {
		var shipname, shipnickname, crewsize string
		err := dbrows.Scan(&shipname, &shipnickname, &crewsize)
		if err != nil {
			panic(err.Error())
		}
		resultMessage = resultMessage + shipname + "\t" + shipnickname + "\t" + crewsize + "\n"
		bFoundShips = true
	}
	if !bFoundShips {
		resultMessage = "No ships for you!"
	}
	return resultMessage
}

func sendUserBio(userSpecified string, m *discordgo.MessageCreate, s *discordgo.Session) {
	dbrows, err := DB.Query("select u.handle, u.shortBio, u.img, r.name as rank, p.name as position from users u, rank r, positions p where u.rank = r.rankid and u.position = p.positionid and u.handle = ?", userSpecified)
	if err != nil {
		panic(err.Error)
	}
	defer dbrows.Close()

	for dbrows.Next() {
		var handle, shortBio, img, rank, position string
		err := dbrows.Scan(&handle, &shortBio, &img, &rank, &position)
		if err != nil {
			panic(err.Error())
		}

		var imgURL string
		if len(img) == 0 {
			imgURL = "https://i.imgur.com/GhsS0cq.jpg"
		} else {
			imgURL = "http://www.novabl4ck.org" + img
		}
		resultMessage := NewEmbed().SetTitle(handle).SetDescription(strip.StripTags(shortBio)).SetColor(0xBA55D3).SetAuthor(m.Author.Username).SetImage(imgURL).AddField("Rank", rank).AddField("Position", position).MessageEmbed
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, resultMessage)
	}
}
