package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lieuweberg/iceberg/util"
)

func init() {
	util.RegisterCommand(util.Command{
		Name: "pause",
		Aliases: []string{"unpause", "resume"},
		Run: pause,
	})
}

func pause(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	if ms, ok := util.Music[m.GuildID]; ok {
		guild, err := s.State.Guild(m.GuildID)
		if err != nil {
			return err
		}
		if isInVoice := util.IsInVoice(guild, m.Author.ID); !isInVoice {
			_, err = s.ChannelMessageSend(m.ChannelID, "You're not in the music channel.")
			return nil
		}
		err = ms.Player.Pause(!ms.Player.Paused())
		
		if ms.Player.Paused() {
			_, err = s.ChannelMessageSend(m.ChannelID, "Player paused.")
		} else {
			_, err = s.ChannelMessageSend(m.ChannelID, "Player resumed.")
		}
	} else {
		_, err = s.ChannelMessageSend(m.ChannelID, "There is no music playing.")
	}
	return
}