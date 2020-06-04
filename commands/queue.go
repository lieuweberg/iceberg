package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/lieuweberg/iceberg/util"
)

func init() {
	util.RegisterCommand(util.Command{
		Name: "queue",
		Run: queue,
	})
}

func queue(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	ms, ok := util.Music[m.GuildID]
	if ok && len(ms.Queue) > 0 {
		res := ""
		for i, song := range ms.Queue {
			res += fmt.Sprintf("\n**%d: **%s", i+1, song.Track.Info.Title)
		}
		_, err = s.ChannelMessageSend(m.ChannelID, "Song queue:\n" + res)
	} else {
		_, err = s.ChannelMessageSend(m.ChannelID, "There is no music playing in this server.")
	}
	return
}