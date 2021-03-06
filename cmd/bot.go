package main

import (
	"flag"
	"github.com/gofree-wtf/discord-cloud-status-bot/pkg"
	"github.com/gofree-wtf/discord-cloud-status-bot/pkg/api"
	. "github.com/gofree-wtf/discord-cloud-status-bot/pkg/common"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	// os environments
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		Logger.Info().Str(pair[0], pair[1]).Msg("os environment")
	}

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

	Logger.Info().Interface("config", Config).Msg("success to load config")

	// run bot controller
	ctrl, err := pkg.NewBotController(Config.Bot.Token, Config.Bot.CommandPrefix, Config.Log.SessionLogLevel)
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
