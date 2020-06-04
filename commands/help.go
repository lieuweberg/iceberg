package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lieuweberg/iceberg/util"
)

func init() {
	util.RegisterCommand(util.Command{
		Name: "help",
		Run: help,
		Aliases: []string{"h"},
	})
}

func help(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	reply := "**Iceberg help menu**\nCurrent prefix: `" + util.Config.Prefix + "`\n"

	for _, c := range util.Commands {
		reply += "\n - `" + c.Name + "`"
		if len(c.Aliases) > 0 {
			reply += "("
			for i, a := range c.Aliases {
				sep := ""
				if len(c.Aliases) != i+1 {
					sep = ", "
				}
				reply += "`" + a + "`" + sep
			}
			reply += ")"
		}
	}

	_, err = s.ChannelMessageSend(m.ChannelID, reply)

	return
}