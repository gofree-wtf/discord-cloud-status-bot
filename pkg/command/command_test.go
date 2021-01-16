package command

import (
	"fmt"
	. "github.com/gofree-wtf/discord-cloud-status-bot/pkg/common"
	. "github.com/stretchr/testify/assert"
	"testing"
)

func TestHelp(t *testing.T) {
	c := NewBotCommand(DefaultBotCommandPrefix)
	outMsg, _, err := c.Execute("help")
	NoError(t, err)
	fmt.Println(outMsg)
}

func TestPing(t *testing.T) {
	c := NewBotCommand(DefaultBotCommandPrefix)
	outMsg, _, err := c.Execute(pingCmd)
	NoError(t, err)
	Equal(t, pingOutMsgPrefix, outMsg)
}

func TestPingHelp(t *testing.T) {
	c := NewBotCommand(DefaultBotCommandPrefix)
	outMsg, _, err := c.Execute(pingCmd + " -h")
	NoError(t, err)
	fmt.Println(outMsg)
}

func TestPingWithArgs(t *testing.T) {
	c := NewBotCommand(DefaultBotCommandPrefix)
	outMsg, _, err := c.Execute(pingCmd + ` test 'test' "test"`)
	NoError(t, err)
	fmt.Println(outMsg)
}
