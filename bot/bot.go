package bot

import (
	"database/sql"
	"fmt"
	"strconv"
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

// Rank structure for rank array
type Rank struct {
	RankID   int
	Name     string
	Sequence int
}

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
			return
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
			return
		}

		if strings.HasPrefix(command, "ship") {
			shipSpecified := strings.TrimSpace(strings.TrimPrefix(command, "ship"))
			parsedFields := strings.Fields(shipSpecified)
			if len(parsedFields) == 0 {
				sendShipManufacturers(m, s)
			} else if len(parsedFields) == 1 {
				if _, err := strconv.Atoi(parsedFields[0]); err == nil {
					sendShipInfoByID(m, s, parsedFields[0])
				} else {
					sendShipsForManufacturer(m, s, parsedFields[0])
				}
			} else if len(parsedFields) >= 2 {
				sendShipInfo(m, s, parsedFields)
			}
			return
		}

		if strings.HasPrefix(command, "shitlist") {
			userSpecified := strings.TrimSpace(strings.TrimPrefix(command, "shitlist"))
			shitlist(userSpecified, m, s)
		}
	}
}

func buildHelpMessage() string {
	resultString := "HELP TOPICS...\n"
	resultString = resultString + "============================\n"
	resultString = resultString + "!help - this content, duh!\n"
	resultString = resultString + "  NOTE: All commands accept an optional [handle] argument - if specified, I will return the requested\n"
	resultString = resultString + "  info for that user, if I can find that user, so make sure his or her handle is correct\n\n"
	resultString = resultString + "!ships [handle] - displays a list of ships you or the specified player owns\n"
	resultString = resultString + "!ship [manufacturer] [name] - displays information about a specified ship\n"
	resultString = resultString + "!bio [handle] - displays a player's BIO in the organization\n"
	return resultString
}

func buildShipList(userSpecified string) string {
	dbrows, err := DB.Query("select s.manufacturer, s.name, s.crewsize from ships s, ownedShips os, users u where os.status = 1 and s.id = os.shipid and os.userid = u.userid and u.handle = ?", userSpecified)
	if err != nil {
		panic(err.Error())
	}
	defer dbrows.Close()

	var resultMessage string
	var bFoundShips = false

	resultMessage = "Manufacturer | Ship Name | Crew size\n"

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

func sendShipManufacturers(m *discordgo.MessageCreate, s *discordgo.Session) {
	queryString := `select distinct(manufacturer) from ships where active = 1;`
	dbrows, err := DB.Query(queryString)
	if err != nil {
		panic(err.Error())
	}
	defer dbrows.Close()

	var resultMessage string
	resultMessage = "Ship Manufacturers\n==================\n"
	for dbrows.Next() {
		var manuName string
		err := dbrows.Scan(&manuName)
		if err != nil {
			panic(err.Error())
		}
		resultMessage = resultMessage + manuName + "\n"
	}
	_, _ = s.ChannelMessageSend(m.ChannelID, resultMessage)
	return
}

func sendShipsForManufacturer(m *discordgo.MessageCreate, s *discordgo.Session, manufacturer string) {
	queryString := `select id, name from ships where manufacturer = "` + manufacturer + `" and active = 1;`
	dbrows, err := DB.Query(queryString)
	if err != nil {
		panic(err.Error())
	}
	defer dbrows.Close()

	var resultMessage = "Ships for " + manufacturer + "\n"
	resultMessage = resultMessage + strings.Repeat("=", len(resultMessage)) + "\n"
	for dbrows.Next() {
		var id, shipName string
		err := dbrows.Scan(&id, &shipName)
		if err != nil {
			panic(err.Error())
		}
		resultMessage = resultMessage + "(#" + id + " )" + shipName + "\n"
	}
	_, _ = s.ChannelMessageSend(m.ChannelID, resultMessage)
	return
}

func sendShipInfoByID(m *discordgo.MessageCreate, s *discordgo.Session, shipStr string) {
	shipID, err := strconv.Atoi(shipStr)
	if err != nil {
		panic(err.Error())
	}
	dbrows, err := DB.Query(`select id, img, name, manufacturer, nickname, crewsize, count(t2.key) qtyInOrg from ships left join (ownedShips t2) on (id = shipid) where t2.status = 1 and active = 1 and shipId = ?`, shipID)
	if err != nil {
		panic(err.Error())
	}
	defer dbrows.Close()

	for dbrows.Next() {
		var id, img, name, manufacturer, nickname, crewsize, qtyInOrg string
		err := dbrows.Scan(&id, &img, &name, &manufacturer, &nickname, &crewsize, &qtyInOrg)
		if err != nil {
			panic(err.Error())
		}

		var imgURL string

		if len(img) == 0 {
			imgURL = "https://i.imgur.com/GhsS0cq.jpg"
		} else {
			imgURL = "http://www.novabl4ck.org" + img
		}
		if len(nickname) == 0 {
			nickname = name
		}
		q, _ := strconv.Atoi(qtyInOrg)
		if q > 0 {

			usrhandles, err := DB.Query(`select u.handle from ownedShips o inner join (users u) on (u.userid = o.userid) where u.status = 1 and o.status = 1 and o.shipid = ?`, shipID)
			if err != nil {
				panic(err.Error())
			}
			defer usrhandles.Close()

			var users []string
			for usrhandles.Next() {
				var userhandle string
				err := usrhandles.Scan(&userhandle)
				if err != nil {
					panic(err.Error())
				}
				users = append(users, userhandle)
			}

			members := strings.Join(users, "\n")
			resultMessage := NewEmbed().SetTitle(name).SetDescription(manufacturer+" "+name).SetColor(0xBA55D3).SetAuthor(m.Author.Username).SetImage(imgURL).AddField("Crew Size", crewsize).AddField("Nickname", nickname).AddField("Qty in the Org", qtyInOrg).AddField("Members who own one:", members).MessageEmbed
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, resultMessage)
			if err != nil {
				panic(err.Error())
			}
		} else {
			resultMessage := NewEmbed().SetTitle(name).SetDescription(manufacturer+" "+name).SetColor(0xBA55D3).SetAuthor(m.Author.Username).SetImage(imgURL).AddField("Crew Size", crewsize).AddField("Nickname", nickname).AddField("Qty in the Org", qtyInOrg).MessageEmbed
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, resultMessage)
			if err != nil {
				panic(err.Error())
			}
		}
	}
}

func sendShipInfo(m *discordgo.MessageCreate, s *discordgo.Session, fields []string) {
	var shipName, manufacturer string
	manufacturer = fields[0]
	for x := 1; x < len(fields); x++ {
		shipName = shipName + " " + fields[x]
	}
	shipName = strings.TrimSpace(shipName)
	queryString := `select id, img, name, manufacturer, nickname, crewsize, count(t2.key) qtyInOrg from ships left join (ownedShips t2) on (id = shipid) where active = 1 and manufacturer = "` + manufacturer + `" and name = "` + shipName + `"`
	dbrows, err := DB.Query(queryString)
	if err != nil {
		panic(err.Error())
	}
	defer dbrows.Close()

	for dbrows.Next() {
		var id, img, name, manufacturer, nickname, crewsize, qtyInOrg string
		err := dbrows.Scan(&id, &img, &name, &manufacturer, &nickname, &crewsize, &qtyInOrg)
		if err != nil {
			panic(err.Error())
		}

		var imgURL string
		if len(img) == 0 {
			imgURL = "https://i.imgur.com/GhsS0cq.jpg"
		} else {
			imgURL = "http://www.novabl4ck.org" + img
		}
		if len(nickname) == 0 {
			nickname = name
		}
		q, _ := strconv.Atoi(qtyInOrg)
		if q > 0 {
			shipID, _ := strconv.Atoi(id)
			usrhandles, err := DB.Query(`select u.handle from ownedShips o inner join (users u) on (u.userid = o.userid) where o.shipid = ?`, shipID)
			if err != nil {
				panic(err.Error())
			}
			defer usrhandles.Close()

			var users []string
			for usrhandles.Next() {
				var userhandle string
				err := usrhandles.Scan(&userhandle)
				if err != nil {
					panic(err.Error())
				}
				users = append(users, userhandle)
			}

			members := strings.Join(users, "\n")
			resultMessage := NewEmbed().SetTitle(name).SetDescription(manufacturer+" "+name).SetColor(0xBA55D3).SetAuthor(m.Author.Username).SetImage(imgURL).AddField("Crew Size", crewsize).AddField("Nickname", nickname).AddField("Qty in the Org", qtyInOrg).AddField("Members who own one:", members).MessageEmbed
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, resultMessage)
			if err != nil {
				panic(err.Error())
			}
		} else {
			resultMessage := NewEmbed().SetTitle(name).SetDescription(manufacturer+" "+name).SetColor(0xBA55D3).SetAuthor(m.Author.Username).SetImage(imgURL).AddField("Crew Size", crewsize).AddField("Nickname", nickname).AddField("Qty in the Org", qtyInOrg).MessageEmbed
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, resultMessage)
			if err != nil {
				panic(err.Error())
			}
		}
	}
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

func shitlist(shitlistedUser string, m *discordgo.MessageCreate, s *discordgo.Session) {

	//myRolename := getUserDiscordRole(m, s)
	var myRolename string
	myRolename = "Captain"

	var maxLevel int

	ranks := getOrgRanks()
	for _, oRank := range ranks {
		if oRank.Name == myRolename {
			maxLevel = oRank.Sequence
		}
	}

	dbrows, err := DB.Query("select rank, handle from users where status = 1 and handle = ?", shitlistedUser)
	if err != nil {
		panic(err.Error())
	}
	defer dbrows.Close()

	for dbrows.Next() {
		var hisranks, hishandle string
		err := dbrows.Scan(&hisranks, &hishandle)
		if err != nil {
			panic(err.Error())
		}
		hisrankID, _ := strconv.Atoi(hisranks)
		hisSequence := getRankSequence(hisrankID)
		if maxLevel < hisSequence {
			// he can be shitlisted
			fmt.Printf("%s can be shitlisted!", hishandle)
		} else {
			fmt.Printf("%s's rank is too high for you to shitlist him; max level you can shitlist is %d", hishandle, maxLevel)
		}
	}
}

func getUserDiscordRole(m *discordgo.MessageCreate, s *discordgo.Session) string {
	var myRoleName string

	guild, err := s.State.Guild(m.GuildID)
	guildRoles := guild.Roles
	if err != nil {
		panic(err.Error())
	}
	myroles := m.Member.Roles
	for _, gRole := range guildRoles {
		if gRole.Name != "@everyone" {
			for _, myRole := range myroles {
				if myRole == gRole.ID {
					myRoleName = gRole.Name
				}
			}
		}
	}
	return myRoleName
}

func getRankSequence(rankID int) int {
	ranks := getOrgRanks()

	for _, rank := range ranks {
		if rank.RankID == rankID {
			return rank.Sequence
		}
	}
	return 1000
}

func getOrgRanks() []Rank {

	var ranks []Rank

	dbrows, err := DB.Query("Select rankid, name, sequence from rank where status = 1 order by sequence asc;")
	if err != nil {
		panic(err.Error())
	}
	defer dbrows.Close()

	for dbrows.Next() {
		var srankid, name, ssequence string
		err := dbrows.Scan(&srankid, &name, &ssequence)
		if err != nil {
			panic(err.Error())
		}

		rankID, _ := strconv.Atoi(srankid)
		sequence, _ := strconv.Atoi(ssequence)
		rank := Rank{RankID: rankID, Name: name, Sequence: sequence}
		ranks = append(ranks, rank)
	}

	return ranks
}
