package util

import "github.com/foxbot/gavalink"

// Lavalink is the lavalink link
var Lavalink *gavalink.Lavalink
// LavaPlayers maps the gavalink.Players to each guild
var LavaPlayers map[string]*gavalink.Player

// LavalinkSetup sets up the lavalink connection and populates Lavalink and Player.
// This function should only be used in the ready event.
func LavalinkSetup(botID string) (err error) {
	Lavalink = gavalink.NewLavalink("1", botID)
	LavaPlayers = make(map[string]*gavalink.Player)

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

// EventHandler handles events that Lavalink may send
type EventHandler struct{}

// OnTrackEnd is raised when a track ends
func (eh EventHandler) OnTrackEnd(player *gavalink.Player, track string, reason string) error {
	return nil
}

// OnTrackException is raised when a track throws an exception
func (eh EventHandler) OnTrackException(player *gavalink.Player, track string, reason string) error {
	return nil
}

// OnTrackStuck is raised when a track gets stuck
func (eh EventHandler) OnTrackStuck(player *gavalink.Player, track string, threshold int) error {
	return nil
}