package common

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

const DefaultBotCommandPrefix = "!cs"

type _Config struct {
	Log Log `yaml:"log"`
	Bot Bot `yaml:"bot"`
}

type Log struct {
	Level        string `yaml:"level"`
	Format       string `yaml:"format"`
	SessionLevel string `yaml:"session_level"`
}

type Bot struct {
	Token         string `yaml:"token"`
	CommandPrefix string `yaml:"command_prefix"`
}

var Config = &_Config{
	// default values
	Log: Log{
		Level:        "info",
		Format:       "console",
		SessionLevel: "warn",
	},
	Bot: Bot{
		CommandPrefix: DefaultBotCommandPrefix,
	},
}

var Logger = log.With().Caller().Logger()

func init() {
	err := setLogger()
	if err != nil {
		Logger.Error().Err(err).Msg("failed to set logger")
	}
}

func setLogger() error {
	lv, err := zerolog.ParseLevel(Config.Log.Level)
	if err != nil {
		Logger.Error().Err(err).Str("log.level", Config.Log.Level).Msg("failed to parse log level")
		return err
	}

	zerolog.SetGlobalLevel(lv)

	switch Config.Log.Format {
	case "console":
		Logger = Logger.Output(zerolog.NewConsoleWriter())
		return nil
	case "json":
		// use default json formatter
		return nil
	default:
		return fmt.Errorf("invalid log.format: %s", Config.Log.Format)
	}
}

func ParseConfigFile(path string) error {
	// parse config file
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, Config)
	if err != nil {
		return err
	}

	// config log module
	err = setLogger()
	if err != nil {
		return err
	}

	// validate values
	if Config.Bot.Token == "" {
		return fmt.Errorf("set your bot.token")
	}
	Config.Bot.CommandPrefix = strings.Trim(Config.Bot.CommandPrefix, " ")

	return nil
}
