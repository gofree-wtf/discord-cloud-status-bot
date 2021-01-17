package command

import (
	"bytes"
	"context"
	. "github.com/gofree-wtf/discord-cloud-status-bot/pkg/common"
	"github.com/google/shlex"
	"github.com/spf13/cobra"
)

type BotCommand struct {
	*cobra.Command
}

func NewBotCommand(commandPrefix string) *BotCommand {
	c := &BotCommand{
		Command: &cobra.Command{
			Use:  commandPrefix,
			Long: "여러 클라우드 서비스의 상태 대시보드를 주기적으로 체크하여, 변경점이 있을 때 알림을 주는 Discord 봇 입니다.",
		},
	}
	c.AddCommand(
		pingCommand(),
		statusCommand(),
	)
	return c
}

func (c *BotCommand) Execute(inMsg string) (outMsg, errMsg string, err error) {
	args, err := shlex.Split(inMsg)
	if err != nil {
		return "", "", err
	}

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	c.SetOut(outBuf)
	c.SetErr(errBuf)
	c.SetArgs(args)

	err = c.Command.Execute()
	return outBuf.String(), errBuf.String(), err
}

func (c *BotCommand) ExecuteWithContext(ctx context.Context, inMsg string) (outMsg, errMsg string, err error) {
	args, err := shlex.Split(inMsg)
	if err != nil {
		return "", "", err
	}

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	c.SetOut(outBuf)
	c.SetErr(errBuf)
	c.SetArgs(args)

	err = c.ExecuteContext(ctx)
	return outBuf.String(), errBuf.String(), err
}

func writeOutMsg(cmd *cobra.Command, outMsg string) {
	_, err := cmd.OutOrStdout().Write([]byte(outMsg))
	if err != nil {
		Logger.Error().Err(err).Str("outMsg", outMsg).Msg("failed to write out message")
	}
}

func writeErrMsg(cmd *cobra.Command, errMsg string) {
	_, err := cmd.OutOrStderr().Write([]byte(errMsg))
	if err != nil {
		Logger.Error().Err(err).Str("errMsg", errMsg).Msg("failed to write err message")
	}
}
