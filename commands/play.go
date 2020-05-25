package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/foxbot/gavalink"
	"github.com/lieuweberg/iceberg/util"
)

func init() {
	util.RegisterCommand(util.Command{
		Name: "play",
		Aliases: []string{"p"},
		Run: play,
	})
}

func play(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		return
	}

	query := strings.Join(strings.Split(m.Content, " ")[1:], " ")
	if query == "" {
		_, err = s.ChannelMessageSend(m.ChannelID, "Please provide a url or search query (ytsearch:QUERY or scsearch:QUERY).")
		return
	}

	var userVoiceState discordgo.VoiceState
	for _, vs := range guild.VoiceStates {
		if vs.UserID == m.Author.ID {
			userVoiceState = *vs
		}
	}
	if userVoiceState.UserID == "" {
		_, err = s.ChannelMessageSend(m.ChannelID, "You are not in a voice channel.")
		return
	}

	alreadyInSameVoice, alreadyInVoice := false, false
	for _, vc := range s.VoiceConnections {
		if vc.ChannelID == userVoiceState.ChannelID {
			alreadyInSameVoice = true
		} else if vc.GuildID == m.GuildID {
			alreadyInVoice = true
		}
	}
	if !alreadyInSameVoice {
		if !alreadyInVoice {
			err := s.ChannelVoiceJoinManual(m.GuildID, userVoiceState.ChannelID, false, false)
			if err != nil {
				return err
			}
		} else {
			_, err = s.ChannelMessageSend(m.ChannelID, "I am already in a different voice channel, not switching.")
			return
		}
	}

	node, err := util.Lavalink.BestNode()
	if err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, "Error finding music node. Please try again.\n" + err.Error())
		return
	}

	tracks, err := node.LoadTracks(query)
	if err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, "Error with query. Please try again or try a different query.\n" + err.Error())
		return
	}

	if tracks.Type != gavalink.TrackLoaded {
		if tracks.Type == gavalink.LoadFailed {
			_, err = s.ChannelMessageSend(m.ChannelID, "Track failed to load. Please try again.")
			return
		} else if tracks.Type == gavalink.NoMatches {
			_, err = s.ChannelMessageSend(m.ChannelID, "No matches for that query. Please try a different query.")
			return
		}
	}
	track := tracks.Tracks[0]
	err = util.LavaPlayers[m.GuildID].Play(track.Data)
	if err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, "Error playing that song. If nothing happens, please try again.\n" + err.Error())
	}

	return
}