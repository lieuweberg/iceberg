package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/foxbot/gavalink"
	"github.com/lieuweberg/iceberg/util"

	// Importing the commands package makes them self-register due to their init function being run,
	// but we don't actually need the package itself (due to self-registering)
	_ "github.com/lieuweberg/iceberg/commands"
)

func main() {
	s, err := discordgo.New("Bot " + util.Config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	s.AddHandler(ready)
	s.AddHandler(messageCreate)
	s.AddHandler(voiceStateUpdate)
	s.AddHandler(voiceServerUpdate)

	err = s.Open()
	if err != nil {
		panic(err.Error())
	}

	// Wait for os terminate events, cleanly close connection when encountered
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT, os.Interrupt, os.Kill)
	<-sc
	fmt.Println("OS termination received, closing WS")
	s.Close()
	fmt.Println("Connection closed, bye bye")
}

func ready(s *discordgo.Session, e *discordgo.Ready) {
	fmt.Printf("Iceberg is online!\n")

	u, err := s.User("@me")
	if err != nil {
		panic(err.Error())
	}

	err = util.LavalinkSetup(u.ID)
	if err != nil {
		panic(err.Error())
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	c, err := s.State.Channel(m.ChannelID)
	if err != nil || c.Type != discordgo.ChannelTypeGuildText {
		return
	}

	if !strings.HasPrefix(m.Content, util.Config.Prefix) {
		return
	}
	m.Content = m.Content[len(util.Config.Prefix):]
	message := strings.Split(m.Content, " ")

	err = util.RunCommand(message[0], s, m)
	if err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, "An error occured during execution.\n"+err.Error())
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func voiceStateUpdate(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
	ms, ok := util.Music[e.GuildID]
	if !ok {
		return
	}

	guild, err := s.State.Guild(e.GuildID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if util.GetUsersInVoice(guild) == 0 {
		fmt.Println("removing player")
		ms.SongEnd <- "end"
		player := ms.Player;
		delete(util.Music, e.GuildID)
		if err = player.Destroy(); err != nil {
			fmt.Println(err)
		}
		if err = s.ChannelVoiceJoinManual(e.GuildID, "", false, false); err != nil {
			fmt.Println(err)
		}
		return
	}
}

func voiceServerUpdate(s *discordgo.Session, e *discordgo.VoiceServerUpdate) {
	vsu := gavalink.VoiceServerUpdate{
		GuildID:  e.GuildID,
		Endpoint: e.Endpoint,
		Token:    e.Token,
	}

	ms, ok := util.Music[e.GuildID]
	if ok && ms.Player != nil {
		err := ms.Player.Forward(s.State.SessionID, vsu)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}

	node, err := util.Lavalink.BestNode()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	handler := new(util.EventHandler)
	player, err := node.CreatePlayer(e.GuildID, s.State.SessionID, vsu, handler)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	util.Music[e.GuildID].Player = player
	fmt.Println("adding player")
}
