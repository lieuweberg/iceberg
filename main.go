package main

import (
	"log"
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
	// Same thing but for the database
	_ "github.com/lieuweberg/iceberg/util/db"
)

func main() {
	s, err := discordgo.New("Bot " + util.Config.Token)
	if err != nil {
		log.Fatalf("Couldn't create session: %s", err)
		return
	}

	s.AddHandler(ready)
	s.AddHandler(messageCreate)
	s.AddHandler(voiceStateUpdate)
	s.AddHandler(voiceServerUpdate)

	err = s.Open()
	if err != nil {
		log.Fatalf("Couldn't create WS: %s", err)
	}

	// Wait for os terminate events, cleanly close connection when encountered
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT, os.Interrupt, os.Kill)
	<-sc
	log.Print("OS termination received, closing WS")
	s.Close()
	log.Print("Connection closed, bye bye")
}

func ready(s *discordgo.Session, e *discordgo.Ready) {
	u, err := s.User("@me")
	if err != nil {
		log.Fatal(err)
	}

	err = util.LavalinkSetup(u.ID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Iceberg is online!")
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
		s.ChannelMessageSend(m.ChannelID, "An error occured during execution.\n"+err.Error())
	}
}

func voiceStateUpdate(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
	ms, ok := util.Music[e.GuildID]
	if !ok {
		return
	}

	guild, err := s.State.Guild(e.GuildID)
	if err != nil {
		log.Printf("No guild found in State for %s: %s", e.GuildID, err)
		return
	}
	if util.GetUsersInVoice(guild) == 0 {
		ms.SongEnd <- "end"
		if returnedMessage := util.LeaveAndDestroy(s, e.GuildID); returnedMessage != "" {
			log.Print(returnedMessage)
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
		if err := ms.Player.Forward(s.State.SessionID, vsu); err != nil {
			log.Printf("Player not forwarded for %s: %s", e.GuildID, err)
		}
		return
	}

	node, err := util.Lavalink.BestNode()
	if err != nil {
		log.Printf("Couldn't find best node: %s", err)
		return
	}

	handler := new(util.EventHandler)
	player, err := node.CreatePlayer(e.GuildID, s.State.SessionID, vsu, handler)
	if err != nil {
		log.Printf("Couldn't create player for %s: %s", e.GuildID, err)
		return
	}

	util.Music[e.GuildID].Player = player
	ms.PlayerCreated <- true
}
