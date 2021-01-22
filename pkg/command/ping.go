package command

import (
	"fmt"
	. "github.com/gofree-wtf/discord-cloud-status-bot/pkg/common"
	"github.com/spf13/cobra"
	"strings"
)

const (
	pingCmd          = "ping"
	pingOutMsgPrefix = "pong !"
)

func pingCommand() *cobra.Command {
	return &cobra.Command{
		Use:   pingCmd,
		Short: "간단한 응답 메세지를 줍니다. 테스트 목적으로 사용할 수 있습니다.",
		Run: func(cmd *cobra.Command, args []string) {
			outMsg := pingOutMsgPrefix
			if len(args) > 0 {
				outMsg = fmt.Sprintf("%s\nArgs: %s", outMsg, strings.Join(args, " "))
			}
			l := Logger.With().Str("outMsg", outMsg).Logger()

			createMsgObj := cmd.Context().Value(CreateMessageKey)
			if createMsgObj == nil {
				writeOutMsg(cmd, outMsg)
				return
			}

			createMsg, ok := createMsgObj.(CreateMessage)
			if !ok {
				l.Error().Interface("createMsgObj", createMsgObj).Msg("invalid message context")
				return
			}

			_, err := createMsg.Session.ChannelMessageSend(createMsg.Message.ChannelID, outMsg)
			if err != nil {
				l.Error().Err(err).Msg("failed to send message")
				return
			}
			l.Info().Msg("success to send message")
		},
	}
}
