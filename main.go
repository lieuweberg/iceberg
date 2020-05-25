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
	s.AddHandler(voiceServerUpdate)

	err = s.Open()
	if err != nil {
		panic(err)
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
		panic(err)
	}

	err = util.LavalinkSetup(u.ID)
	if err != nil {
		panic(err)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if c, err := s.State.Channel(m.ChannelID); err != nil || c.Type != discordgo.ChannelTypeGuildText {
		return
	}

	prefix := "go!"
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}
	m.Content = m.Content[len(prefix):]
	message := strings.Split(m.Content, " ")

	err := util.RunCommand(message[0], s, m)
	if err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, "An error occured during execution.\n"+err.Error())
		if err != nil {
			fmt.Println(err)
		}
	}
}

func voiceServerUpdate(s *discordgo.Session, e *discordgo.VoiceServerUpdate) {
	vsu := gavalink.VoiceServerUpdate{
		GuildID:  e.GuildID,
		Endpoint: e.Endpoint,
		Token:    e.Token,
	}

	if p, err := util.Lavalink.GetPlayer(e.GuildID); err == nil {
		err = p.Forward(s.State.SessionID, vsu)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	node, err := util.Lavalink.BestNode()
	if err != nil {
		fmt.Println(err)
		return
	}

	handler := new(util.EventHandler)
	player, err := node.CreatePlayer(e.GuildID, s.State.SessionID, vsu, handler)
	if err != nil {
		fmt.Println(err)
		return
	}

	util.LavaPlayers[e.GuildID] = player
}
