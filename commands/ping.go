package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lieuweberg/iceberg/util"
)

func init() {
	util.RegisterCommand(util.Command{
		Name: "ping",
		Run: command,
	})
}

func command(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	s.ChannelMessageSend(m.ChannelID, "Pong!")
	return
}