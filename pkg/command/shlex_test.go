package command

import (
	. "github.com/gofree-wtf/discord-cloud-status-bot/pkg/common"
	"github.com/google/shlex"
	. "github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSingleArgs(t *testing.T) {
	args := []string{"help"}
	line := strings.Join(args, " ")

	splits, err := shlex.Split(line)
	NoError(t, err)
	EqualValues(t, args, splits)
}

func TestMultipleArgs(t *testing.T) {
	args := []string{"ping", "--help"}
	line := strings.Join(args, " ")

	splits, err := shlex.Split(line)
	NoError(t, err)
	EqualValues(t, args, splits)
}

func TestSingleQuoteArgs(t *testing.T) {
	args := []string{"ping", "aa"}
	line := `ping 'aa'`

	splits, err := shlex.Split(line)
	NoError(t, err)
	EqualValues(t, args, splits)
}

func TestDoubleQuoteArgs(t *testing.T) {
	args := []string{"ping", "aa"}
	line := `ping "aa"`

	splits, err := shlex.Split(line)
	NoError(t, err)
	EqualValues(t, args, splits)
}

func TestInvalidSingleQuoteArgs(t *testing.T) {
	line := `ping 'aa''`

	_, err := shlex.Split(line)
	Error(t, err)
	Logger.Info().Msg(err.Error())
}

func TestInvalidDoubleQuoteArgs(t *testing.T) {
	line := `ping "aa""`

	_, err := shlex.Split(line)
	Error(t, err)
	Logger.Info().Msg(err.Error())
}
