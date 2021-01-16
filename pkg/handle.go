package pkg

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	. "github.com/gofree-wtf/discord-cloud-status-bot/pkg/common"
	"github.com/gofree-wtf/discord-cloud-status-bot/pkg/util"
	"strings"
)

func (c *BotController) handleMessage(session *discordgo.Session, createMsg *discordgo.MessageCreate) {
	Logger.Debug().Interface("createMsg", createMsg).Msg("")

	if createMsg.Author.ID == session.State.User.ID {
		// ignore bot self message
		return
	} else if !strings.HasPrefix(createMsg.Content, c.commandPrefix+" ") {
		// ignore message without prefix
		return
	}

	command := util.SubstringAfter(createMsg.Content, c.commandPrefix+" ")
	Logger.Info().Str("command", command).Msg("received command")

	// TODO handle command
	sendContent := fmt.Sprintf("I got it. %s command !", command)
	_, err := session.ChannelMessageSend(createMsg.ChannelID, sendContent)
	if err != nil {
		Logger.Error().Err(err).Str("content", sendContent).Msg("failed to send message")
	}
}
