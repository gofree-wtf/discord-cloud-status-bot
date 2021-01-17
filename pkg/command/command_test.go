package command

import (
	"fmt"
	. "github.com/gofree-wtf/discord-cloud-status-bot/pkg/common"
	. "github.com/stretchr/testify/assert"
	"testing"
)

var command = NewBotCommand(DefaultBotCommandPrefix)

func TestHelp(t *testing.T) {
	outMsg, _, err := command.Execute("help")
	NoError(t, err)
	fmt.Println(outMsg)
}

func TestPing(t *testing.T) {
	outMsg, _, err := command.Execute(pingCmd)
	NoError(t, err)
	Equal(t, pingOutMsgPrefix, outMsg)
}

func TestPingHelp(t *testing.T) {
	outMsg, _, err := command.Execute(pingCmd + " -h")
	NoError(t, err)
	fmt.Println(outMsg)
}

func TestPingWithArgs(t *testing.T) {
	outMsg, _, err := command.Execute(pingCmd + ` test 'test' "test"`)
	NoError(t, err)
	fmt.Println(outMsg)
}

func TestStatus(t *testing.T) {
	outMsg, _, err := command.Execute(statusCmd)
	NoError(t, err)
	fmt.Println(outMsg)
}

func TestStatusDiscord(t *testing.T) {
	outMsg, _, err := command.Execute(statusCmd + " " + statusDiscordCmd)
	NoError(t, err)
	fmt.Println(outMsg)
}
