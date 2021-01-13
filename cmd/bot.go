package main

import (
	"github.com/gofree-wtf/discord-cloud-status-bot/pkg"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("discord-cloud-status-bot - start")

	pkg.Test()

	log.Info().Msg("discord-cloud-status-bot - end")
}
