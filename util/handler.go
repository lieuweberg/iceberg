package util

import "github.com/bwmarrin/discordgo"

// Command is the command struct with the data that should be filled in
type Command struct {
	Name string
	Aliases []string
	Run func(s *discordgo.Session, m *discordgo.MessageCreate) (err error)
}

var commands map[string]Command
var aliases map[string]string

func init() () {
	commands = make(map[string]Command)
	aliases = make(map[string]string)
}

// RegisterCommand registers a Command struct to be used with the bot
func RegisterCommand(c Command) {
	commands[c.Name] = c
	for _, alias := range c.Aliases {
		aliases[alias] = c.Name
	}
}

// RunCommand runs the command (specified by name) and also looks for any aliases.
// It will return and error if there is any in command execution (essentially forwarding it),
// but if a command is not found it will not return an error since that is not needed.
func RunCommand(commandName string, s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	if c, ok := commands[commandName]; ok {
		err = c.Run(s, m)
		return
	}

	if c, ok := aliases[commandName]; ok {
		err = commands[c].Run(s, m)
		return
	}
	
	return
}