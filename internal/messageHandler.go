package internal

import (
	"log/slog"
	"strings"

	"github.com/Diamon0/fated-roll/internal/commands"
	"github.com/fluxergo/fluxergo/events"
	"github.com/fluxergo/fluxergo/fluxer"
)

func MessageHandler(e *events.MessageCreate) {
	if e.Message.Author.Bot {
		return
	}

	if e.Message.Content[0:1] != "!" {
		return
	}

	cmd := strings.Split(e.Message.Content, " ")
	var message string
	allowedMentions := fluxer.AllowedMentions{}

	switch cmd[0] {
	case "!r":
		if len(cmd) < 2 {
			message = "To roll, you must do `!r [NumberOfDie]d[NumberOfFaces]`, and optionally `kh` and/or `kl` (in that order) to keep the highest or lowest rolls. Additionally, you may add a math expression at the end for modifying the result, such as +5, or -3*6+5 (Following PEMDAS)"
			break
		}

		allowedMentions.Users = append(allowedMentions.Users, e.Message.Author.ID)
		mention := fluxer.UserMention(e.Message.Author.ID)
		commandMessage := commands.Roll(cmd[1:])
		message = mention + ":game_die:\n" + commandMessage

	default:
		return
	}

	if message != "" {
		if _, err := e.Client().Rest.CreateMessage(e.ChannelID, fluxer.NewMessageCreate().WithContent(message)); err != nil {
			slog.Error("Failed to send roll message", slog.Any("error", err))
		}
	}
}
