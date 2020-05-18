package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/lieuweberg/iceberg/util"
)

var (
	config util.Configuration
)

func init() {
	configuration, err := util.Config()
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

	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Println(err.Error())
		_, err = s.ChannelMessageSend(m.ChannelID, "Error getting this channel, please try again")
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}

	if channel.Type != discordgo.ChannelTypeGuildText {
		return
	}

	if m.Content == "go!ping" {
		_, err = s.ChannelMessageSend(m.ChannelID, "Pong! :D " + strconv.Itoa(int(s.HeartbeatLatency().Milliseconds())) + "ms")
		if err != nil {
			fmt.Println(err.Error())
		}

		return
	}
}