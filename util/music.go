package util

import (
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/foxbot/gavalink"
)

// Lavalink is the lavalink connection
var Lavalink *gavalink.Lavalink
// Music maps the music struct to each guild
var Music map[string]*MusicStruct

// MusicStruct is the struct to use in the Music map
type MusicStruct struct {
	ChannelID string
	Queue []Song
	SongEnd chan string
	Player *gavalink.Player
}

// Song is the struct to use in a queue
type Song struct {
	Requester string
	Track gavalink.Track
	Skips float64
}

// LavalinkSetup sets up the lavalink connection and populates Lavalink and Player.
// This function should only be used in the ready event.
func LavalinkSetup(botID string) (err error) {
	Lavalink = gavalink.NewLavalink("1", botID)
	Music = make(map[string]*MusicStruct)

	err = Lavalink.AddNodes(gavalink.NodeConfig{
		REST: "http://localhost:2333",
		WebSocket: "ws://localhost:2333",
		Password: "nicemusicbro",
		
	})
	if err != nil {
		return
	}

	return
}

// GetUsersInVoice returns the number of people in the same voice channel as the bot, excluding the bot
func GetUsersInVoice(guild *discordgo.Guild) int {
	usersInVoice := 0
	for _, vs := range guild.VoiceStates {
		if Music[guild.ID].ChannelID == vs.ChannelID {
			usersInVoice++
		}
	}
	return usersInVoice - 1
}

// IsInVoice checks if a user is in the same voice channel as where music is playing
func IsInVoice(guild *discordgo.Guild, userID string) bool {
	for _, vs := range guild.VoiceStates {
		if Music[guild.ID].ChannelID == vs.ChannelID && userID == vs.UserID {
			return true
		}
	}
	return false
}

// EventHandler handles Lavalink track events
type EventHandler struct{}

// OnTrackEnd is raised when a track ends
func (eh EventHandler) OnTrackEnd(player *gavalink.Player, track string, reason string) (err error) {
	if ms, ok := Music[player.GuildID()]; ok {
		ms.Queue = ms.Queue[1:]
		if len(ms.Queue) != 0 {
			ms.SongEnd <- "next"
		} else {
			ms.SongEnd <- "end"
			ms.Queue = make([]Song, 0)
		}
	}
	return
}

// OnTrackException is raised when a track throws an exception
func (eh EventHandler) OnTrackException(player *gavalink.Player, track string, reason string) (err error) {
	log.Printf("Track exception for %s: %s", player.GuildID, reason)
	if ms, ok := Music[player.GuildID()]; ok {
		ms.SongEnd <- reason
	}
	return
}

// OnTrackStuck is raised when a track gets stuck
func (eh EventHandler) OnTrackStuck(player *gavalink.Player, track string, threshold int) (err error) {
	if ms, ok := Music[player.GuildID()]; ok {
		ms.SongEnd <- "resume:" + strconv.Itoa(threshold)
	}
	return
}