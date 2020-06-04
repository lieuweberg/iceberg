package commands

import (
	"fmt"
	"strconv"
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

	query := strings.Replace(strings.Join(strings.Split(m.Content, " ")[1:], " "), " ", "%20", -1)
	if query == "" {
		_, err = s.ChannelMessageSend(m.ChannelID, "Please provide a url or search query.")
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

	var ms *util.MusicStruct
	alreadyInSameVoice, alreadyInVoice := false, false
	if temp, ok := util.Music[m.GuildID]; ok {
		ms = temp;
		for _, vs := range guild.VoiceStates {
			if ms.ChannelID == vs.ChannelID {
				if vs.UserID == m.Author.ID {
					alreadyInSameVoice = true
					break
				}
				alreadyInVoice = true
			}
		}
	}
	if !alreadyInSameVoice {
		if !alreadyInVoice {
			err := s.ChannelVoiceJoinManual(m.GuildID, userVoiceState.ChannelID, false, false)
			if err != nil {
				return err
			}

			util.Music[m.GuildID] = &util.MusicStruct{
				ChannelID: userVoiceState.ChannelID,
				Queue: make([]util.Song, 0),
				SongEnd: make(chan string),
			}
			ms = util.Music[m.GuildID]
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

	if !strings.HasPrefix(query, "http") && !strings.Contains(query, "://") {
		query = "ytsearch:" + query
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
	err = queueSong(s, m, track, ms)
	if err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, "Error adding song to queue. Please try again.\n" + err.Error())
		return
	}

	return
}


func queueSong(s *discordgo.Session, m *discordgo.MessageCreate, track gavalink.Track, ms *util.MusicStruct) (err error) {
	ms.Queue = append(ms.Queue, util.Song{
		Requester: m.Author.Username,
		Track: track,
	})

	if len(ms.Queue) == 1 {
		err = playSong(s, m, ms.Queue[0], ms, -1)
		if err != nil {
			_, err = s.ChannelMessageSend(m.ChannelID, "An error occured during song playback. Please try again.\n" + err.Error())
			return
		}
	} else {
		_, err = s.ChannelMessageSend(m.ChannelID, "Queued for playback: **" + track.Info.Title + "**")
	}

	return
}

func playSong(s *discordgo.Session, m *discordgo.MessageCreate, song util.Song, ms *util.MusicStruct, startTime int) (err error) {
	fmt.Println("attempting to play", song.Track.Info.Title)
	if startTime < 0 {
		err = ms.Player.Play(song.Track.Data)
		fmt.Println("playing", song.Track.Info.Title)
	} else {
		err = ms.Player.PlayAt(song.Track.Data, startTime, song.Track.Info.Length)
	}
	if err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, "Error playing *" + song.Track.Info.Title + "*. Skipping to next song.\n" + err.Error())
		ms.Queue = ms.Queue[1:]
		if len(ms.Queue) != 0 {
			playSong(s, m, ms.Queue[0], ms, -1)
		}
		return
	}
	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":musical_note: Now playing: **%s** *(requested by %s)*", song.Track.Info.Title, song.Requester))

	end := <- ms.SongEnd
	fmt.Println(end)
	if end == "next" {
		err = playSong(s, m, ms.Queue[0], ms, -1)
	} else if end == "end" {
		_, err = s.ChannelMessageSend(m.ChannelID, ":musical_note: Playback finished.")
	} else if strings.HasPrefix(end, "resume:") {
		end = end[7:]
		time, err := strconv.Atoi(end)
		if err != nil {
			return err
		}
		err = playSong(s, m, song, ms, time)
	} else {
		_, err = s.ChannelMessageSend(m.ChannelID, "An error occured with the track")
	}

	return
}