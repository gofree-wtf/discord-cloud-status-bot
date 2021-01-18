package main

import (
	"flag"
	"github.com/gofree-wtf/discord-cloud-status-bot/pkg"
	"github.com/gofree-wtf/discord-cloud-status-bot/pkg/api"
	. "github.com/gofree-wtf/discord-cloud-status-bot/pkg/common"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// program arguments
	var configFile string

	flag.StringVar(&configFile, "config.file", "", "config file path")
	flag.Parse()

	Logger.Info().Str("-config.file", configFile).Msg("program argument")

	// load config
	err := LoadConfig(configFile)
	if err != nil {
		Logger.Fatal().Err(err).Msg("failed to load config")
	}

	Logger.Info().Interface("config", Config).Msg("")

	// run bot controller
	ctrl, err := pkg.NewBotController(Config.Bot.Token, Config.Bot.CommandPrefix, Config.Log.SessionLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create controller")
	}
	ctrlCloseFn := ctrl.Run()

	// run api server
	apiServer := api.NewApiServer()
	apiServerCloseFn := apiServer.Run()

	// handle os signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sc

	// graceful shutdown
	ctrlCloseFn()
	apiServerCloseFn()
}
