package command

import (
	"encoding/json"
	"github.com/gofree-wtf/discord-cloud-status-bot/pkg/client/discord"
	. "github.com/gofree-wtf/discord-cloud-status-bot/pkg/common"
	"github.com/spf13/cobra"
)

const (
	statusCmd        = "status"
	statusDiscordCmd = "discord"
)

func statusCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   statusCmd,
		Short: "각 서비스의 상태를 확인할 수 있습니다.",
	}
	c.AddCommand(statusDiscordCommand())
	return c
}

func statusDiscordCommand() *cobra.Command {
	return &cobra.Command{
		Use:   statusDiscordCmd,
		Short: "Discord 상태를 확인할 수 있습니다.",
		Run: func(cmd *cobra.Command, args []string) {
			status, err := discord.GetDiscordStatus()
			if err != nil {
				writeErrMsg(cmd, err.Error())
				return
			}
			Logger.Debug().Interface("discordStatus", status).Msg("")

			jsonStr, err := json.MarshalIndent(status, "", "  ")
			if err != nil {
				writeErrMsg(cmd, err.Error())
				return
			}

			writeOutMsg(cmd, string(jsonStr))
		},
	}
}
