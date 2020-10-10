package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lieuweberg/iceberg/util"
)

func init() {
	util.RegisterCommand(util.Command{
		Name: "stop",
		Run: stop,
		Aliases: []string{"leave", "disconnect", "dc"},
	})
}

func stop(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		return
	}

	if util.IsInVoiceWithMusic(guild, m.Author.ID) {
		if returnedMessage := util.LeaveAndDestroy(s, m.GuildID); returnedMessage != "" {
			s.ChannelMessageSend(m.ChannelID, "Playback may not have been stopped.\n" + returnedMessage)
		} else {
			s.ChannelMessageSend(m.ChannelID, "Playback stopped.")
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "You're not listening to my music :(")
	}
	return
}