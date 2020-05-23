package commands

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/lieuweberg/iceberg/util"
)

func init() {
	util.RegisterCommand(util.Command{
		Name: "ping",
		Run: ping,
	})
}

func ping(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	s.ChannelMessageSend(m.ChannelID, "Pong! " + strconv.Itoa(int(s.HeartbeatLatency().Milliseconds())) + "ms :stopwatch:")
	return
}