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

			_, err := cmd.OutOrStdout().Write([]byte(outMsg))
			if err != nil {
				Logger.Error().Err(err).Str("outMsg", outMsg).Msg("failed to write out message")
			}
		},
	}
}
