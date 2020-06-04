package commands

import (
	"fmt"
	"math"

	"github.com/bwmarrin/discordgo"
	"github.com/lieuweberg/iceberg/util"
)

func init() {
	util.RegisterCommand(util.Command{
		Name: "skip",
		Run: skip,
	})
}

func skip(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	if ms, ok := util.Music[m.GuildID]; ok {
		guild, err := s.State.Guild(m.GuildID)
		if err != nil {
			return nil
		}

		if isInVoice := util.IsInVoice(guild, m.Author.ID); !isInVoice {
			_, err = s.ChannelMessageSend(m.ChannelID, "You're not in the music channel.")
			return nil
		}

		usersInVoice := math.Floor(float64(util.GetUsersInVoice(guild) / 2))
		skips := ms.Queue[0].Skips
		requirementFloat := (skips + 1) / float64(usersInVoice)
		if usersInVoice <= 2 || requirementFloat >= 0.4 {
			ms.Player.Stop()
		} else {
			skips++
			_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vote added! Need %d more (%d/%d).", int(usersInVoice-skips), int(skips), int(usersInVoice)))
		}
	} else {
		_, err = s.ChannelMessageSend(m.ChannelID, "There is no music playing.")
	}
	return
}