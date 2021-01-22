package pkg

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	. "github.com/gofree-wtf/discord-cloud-status-bot/pkg/common"
	"github.com/gofree-wtf/discord-cloud-status-bot/pkg/util"
	"strings"
)

func (c *BotController) handleMessage(session *discordgo.Session, createMsg *discordgo.MessageCreate) {
	if createMsg.Author.ID == session.State.User.ID {
		// ignore bot self message
		return
	} else if !strings.HasPrefix(createMsg.Content, c.commandPrefix+" ") {
		// ignore message without prefix
		return
	}
	Logger.Debug().Interface("createMsg", createMsg).Msg("")

	inMsg := util.SubstringAfter(createMsg.Content, c.commandPrefix+" ")

	l := Logger.With().Str("inMsg", inMsg).Logger()
	l.Info().Msg("received message")

	var ctx = context.WithValue(context.Background(), CreateMessageKey, CreateMessage{
		Session: session,
		Message: createMsg,
	})
	var sendMsg string

	outMsg, errMsg, err := c.command.ExecuteWithContext(ctx, inMsg)
	if err != nil {
		if errMsg != "" {
			sendMsg = fmt.Sprintf("무언가 문제가 발생하였습니다.\n%s", errMsg)
		} else {
			sendMsg = fmt.Sprintf("무언가 문제가 발생하였습니다.\nError: %s", err.Error())
		}
	} else {
		sendMsg = outMsg
	}

	if sendMsg == "" {
		return
	}
	l = l.With().Str("sendMsg", sendMsg).Logger()

	_, err = session.ChannelMessageSend(createMsg.ChannelID, sendMsg)
	if err != nil {
		l.Error().Err(err).Msg("failed to send message")
		return
	}
	l.Info().Msg("success to send message")
}
