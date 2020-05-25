package commands

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/bwmarrin/discordgo"
	"github.com/lieuweberg/iceberg/util"
)

func init() {
	util.RegisterCommand(util.Command{
		Name: "server",
		Run: server,
	})
}

type serverQuery struct {
	Status      string   `json:"status"`
	Online      bool     `json:"online"`
	Error       string   `json:"error"`
	MOTD        string   `json:"motd"`
	Version     string   `json:"version"`
	GameType    string   `json:"game_type"`
	GameID      string   `json:"game_id"`
	ServerMod   string   `json:"server_mod"`
	Map         string   `json:"map"`
	Players     players  `json:"players"`
	Plugins     []string `json:"plugins"`
	LastOnline  string   `json:"last_online"`
	LastUpdated string   `json:"last_updated"`
	Duration    int      `json:"duration"`
}

type players struct {
	Max  int      `json:"max"`
	Now  int      `json:"now"`
	List []string `json:"list"`
}

type queryCache struct {
	lastQuery *serverQuery
	lastUpdated int64
}
var cache queryCache

func server(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	query := &serverQuery{}
	var footer string
	if cache.lastQuery == nil || time.Now().Unix() - cache.lastUpdated > 300 {
		resp, err := http.Get("https://mcapi.us/server/query?ip=server.blockhermit.com")
		if err != nil {
			return err
		}
	
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
	
		err = json.Unmarshal(data, query)
		if err != nil {
			return err
		}

		cache.lastQuery = query
		t, err := strconv.Atoi(query.LastUpdated)
		if err != nil {
			return err
		}
		cache.lastUpdated = int64(t)

		footer = "queried"
	} else {
		query = cache.lastQuery
		footer = "from cache"
	}

	if query.Status != "success" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Something went wrong with fetching the data from the server. Try again later.")
		return err
	}

	statusEmoji := "<:Online:417788119824203776>"
	color := 0x43B581
	if !query.Online {
		statusEmoji = "<:Offline:417788109300695060>"
		color = 0xF04747
	}

	playerList := "There are no players online"
	if query.Players.Now > 0 {
		playerList = strings.Join(query.Players.List, ", ")
	}

	lastUpdated, err := strconv.Atoi(query.LastUpdated)
	if err != nil {
		return
	}
	lastUpdated = int(time.Now().Unix()) - lastUpdated
	var lastUpdatedFormatted string
	if lastUpdated >= 60 {
		if lastUpdated < 1 {
			lastUpdatedFormatted = "Just now"
		} else if lastUpdated % 60 == 0 {
			lastUpdatedFormatted = strconv.Itoa(lastUpdated/60) + " minutes ago"
		} else {
			lastUpdatedFormatted = strconv.Itoa(lastUpdated/60) + " minutes and " + strconv.Itoa(lastUpdated%60) + " seconds ago"
		}
	} else {
		lastUpdatedFormatted = strconv.Itoa(lastUpdated) + " seconds ago"
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Color: color,
		Title: "BlockHermit",
		Description: statusEmoji + " server.blockhermit.com",
		Footer: &discordgo.MessageEmbedFooter{
			Text: "data " + footer,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: strconv.Itoa(query.Players.Now) + "/" + strconv.Itoa(query.Players.Max) + " players online",
				Value: "```" + playerList + "```",
			},
			{
				Name: "Last updated",
				Value: lastUpdatedFormatted,
			},
			{
				Name: "Running",
				Value: query.ServerMod,
			},
		},
	})

	return nil
}