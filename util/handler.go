package util

import "github.com/bwmarrin/discordgo"

// Command is the command struct with the data that should be filled in
type Command struct {
	Name string
	Aliases []string
	Run func(s *discordgo.Session, m *discordgo.MessageCreate) (err error)
}

// Commands is the map of command names to their Command struct (thus their function)
var Commands map[string]Command
// Aliases is the map of aliases to command names
var Aliases map[string]string

func init() () {
	Commands = make(map[string]Command)
	Aliases = make(map[string]string)
}

// RegisterCommand registers a Command struct to be used with the bot
func RegisterCommand(c Command) {
	Commands[c.Name] = c
	for _, alias := range c.Aliases {
		Aliases[alias] = c.Name
	}
}

// RunCommand runs the command (specified by name) and also looks for any aliases.
// It will return and error if there is any in command execution (essentially forwarding it),
// but if a command is not found it will not return an error since that is not needed.
func RunCommand(commandName string, s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	if c, ok := Commands[commandName]; ok {
		err = c.Run(s, m)
		return
	}

	if c, ok := Aliases[commandName]; ok {
		err = Commands[c].Run(s, m)
		return
	}
	
	return
}