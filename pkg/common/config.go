package common

import (
	"fmt"
	"github.com/Netflix/go-env"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
	"time"
)

const DefaultBotCommandPrefix = "!cs"

type _Config struct {
	Log Log `yaml:"log"`
	Bot Bot `yaml:"bot"`
}

type Log struct {
	Level        string `yaml:"level" env:"LOG_LEVEL,default=info"`
	Format       string `yaml:"format" env:"LOG_FORMAT,default=console"`
	SessionLevel string `yaml:"session_level" env:"LOG_SESSION_LEVEL,default=warn"`
}

type Bot struct {
	Token         string `yaml:"token"  env:"BOT_TOKEN" json:"-"`
	CommandPrefix string `yaml:"command_prefix" env:"BOT_COMMAND_PREFIX"`
	TimeZone      string `yaml:"time_zone" env:"BOT_TIMEZONE,default=Asia/Seoul"`

	Location *time.Location `json:"-"`
}

var Config = &_Config{}
var Logger = log.With().Caller().Logger()

func init() {
	_, err := env.UnmarshalFromEnviron(Config)
	if err != nil {
		Logger.Fatal().Err(err).Msg("failed to load environment")
	}

	err = setLogger()
	if err != nil {
		Logger.Error().Err(err).Msg("failed to set logger")
	}

	err = setBotValues()
	if err != nil {
		Logger.Error().Err(err).Msg("failed to set bot values")
	}
}

func ParseConfigFile(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, Config)
	if err != nil {
		return err
	}

	err = setLogger()
	if err != nil {
		return err
	}

	err = setBotValues()
	if err != nil {
		return err
	}
	return nil
}

func ValidateConfig() error {
	if Config.Bot.Token == "" {
		return fmt.Errorf("set your env:BOT_TOKEN or yaml:bot.token")
	}
	return nil
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

func setBotValues() error {
	Config.Bot.CommandPrefix = strings.Trim(Config.Bot.CommandPrefix, " ")

	location, err := time.LoadLocation(Config.Bot.TimeZone)
	if err != nil {
		return err
	}
	Config.Bot.Location = location

	return nil
}
