package pkg

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gofree-wtf/discord-cloud-status-bot/pkg/command"
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
	Logger.Info().Str("inMsg", inMsg).Msg("received message")

	cmd := command.NewBotCommand(c.commandPrefix)

	outMsg, errMsg, err := cmd.Execute(inMsg)

	var sendMsg string
	if err != nil {
		if errMsg != "" {
			sendMsg = fmt.Sprintf("무언가 문제가 발생하였습니다.\n%s", errMsg)
		} else {
			sendMsg = fmt.Sprintf("무언가 문제가 발생하였습니다.\nError: %s", err.Error())
		}
	} else {
		sendMsg = outMsg
	}
	Logger.Debug().Str("sendMsg", sendMsg).Msg("")

	_, err = session.ChannelMessageSend(createMsg.ChannelID, sendMsg)
	if err != nil {
		Logger.Error().Err(err).Str("sendMsg", sendMsg).Msg("failed to send message")
	}
}
