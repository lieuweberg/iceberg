package commands

import (
	"database/sql"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/lieuweberg/iceberg/util"
	"github.com/lieuweberg/iceberg/util/db"
)

func init() {
	util.RegisterCommand(util.Command{
		Name: "verify",
		Run: verify,
	})
}

func verify(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	if u, err := db.GetUser(m.Author.ID); err != nil {
		if err == sql.ErrNoRows {
			s.ChannelMessageSend(m.ChannelID, "Oof, you're not verified.")
		} else {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Oof, an error occured: %s", err))
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Yep, you're verified as %s.", u.McName))
	}

	return
}