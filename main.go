package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/lieuweberg/iceberg/util"
	
	// Importing the commands package makes them self-register due to their init function being run,
	// but we don't actually need the package itself (due to self-registering)
	_ "github.com/lieuweberg/iceberg/commands"
)

var (
	config util.Configuration
)

func init() {
	configuration, err := util.LoadConfig()
	if err != nil {
		panic(err)
	}
	config = configuration
}

func main() {
	s, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	s.AddHandler(messageCreate)

	err = s.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Iceberg is online!\n")
	// Wait for os terminate events, cleanly close connection when encountered
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT, os.Interrupt, os.Kill)
	<-sc
	fmt.Println("OS termination received, closing WS")
	s.Close()
	fmt.Println("Connection closed, bye bye")
}


func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if c, err := s.State.Channel(m.ChannelID); err != nil || c.Type != discordgo.ChannelTypeGuildText {
		return
	}

	prefix := "go!"
	m.Content = m.Content[len(prefix):]
	message := strings.Split(m.Content, " ")

	util.RunCommand(message[0], s, m)
	return
}