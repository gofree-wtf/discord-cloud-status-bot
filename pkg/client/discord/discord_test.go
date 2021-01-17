package discord

import (
	. "github.com/gofree-wtf/discord-cloud-status-bot/pkg/common"
	. "github.com/stretchr/testify/assert"
	"testing"
)

func TestDiscordGetStatus(t *testing.T) {
	status, err := GetDiscordStatus()
	NoError(t, err)
	Logger.Info().Interface("status", status).Msg("")

	updatedAt, err := status.Page.UpdatedAt()
	NoError(t, err)
	Logger.Info().Time("updatedAt", updatedAt).Msg("")
}
