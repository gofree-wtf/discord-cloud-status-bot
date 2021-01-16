package main

import (
	"flag"
	"github.com/gofree-wtf/discord-cloud-status-bot/pkg"
	. "github.com/gofree-wtf/discord-cloud-status-bot/pkg/common"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// program arguments
	var configFile string

	flag.StringVar(&configFile, "config.file", "./configs/app.yaml", "config file path")
	flag.Parse()

	Logger.Info().Str("-config.file", configFile).Msg("")

	// parse config file
	config, err := ParseConfigFile(configFile)
	if err != nil {
		Logger.Fatal().Err(err).Str("config.file", configFile).Msg("failed to parse config file")
	}
	Logger.Debug().Interface("config", config).Msg("")

	// run bot controller
	ctrl, err := pkg.NewBotController(config.Bot.Token, config.Bot.CommandPrefix, config.Log.SessionLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create controller")
	}

	shutdownFn, err := ctrl.Run()
	if err != nil {
		Logger.Fatal().Err(err).Msg("failed to run controller")
	}

	// handle os signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sc
	shutdownFn()
}
