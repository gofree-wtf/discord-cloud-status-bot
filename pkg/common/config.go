package common

import (
	"fmt"
	"github.com/Netflix/go-env"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
	"time"
)

const DefaultBotCommandPrefix = "!cs"

var Logger = log.With().Caller().Logger()

var Config = &_Config{
	// default values
	Log: Log{
		Level:        "info",
		Format:       "console",
		SessionLevel: "warn",
	},
	Bot: Bot{
		CommandPrefix: DefaultBotCommandPrefix,
		TimeZone:      "Asia/Seoul",
		Location:      time.Local,
	},
	Api: Api{
		SelfHealthcheckEnabled:       true,
		SelfHealthcheckPeriodMinutes: 5,
	},
}

type _Config struct {
	Log Log `yaml:"log"`
	Bot Bot `yaml:"bot"`
	Api Api `yaml:"api"`
}

type Log struct {
	Level        string `yaml:"level"         env:"LOG_LEVEL"`
	Format       string `yaml:"format"        env:"LOG_FORMAT"`
	SessionLevel string `yaml:"session_level" env:"LOG_SESSION_LEVEL"`

	SessionLogLevel int `json:"-"`
}

type Bot struct {
	Token         string `yaml:"token"          env:"BOT_TOKEN"           json:"-"`
	CommandPrefix string `yaml:"command_prefix" env:"BOT_COMMAND_PREFIX"`
	TimeZone      string `yaml:"time_zone"      env:"BOT_TIMEZONE"`

	Location *time.Location `json:"-"`
}

type Api struct {
	Port                         int    `yaml:"port"                            env:"API_PORT"`
	SelfHealthcheckEnabled       bool   `yaml:"self_healthcheck_enabled"        env:"API_SELF_HEALTHCHECK_ENABLED"`
	SelfHealthcheckUrl           string `yaml:"self_healthcheck_url"            env:"API_SELF_HEALTHCHECK_URL"`
	SelfHealthcheckPeriodMinutes int    `yaml:"self_healthcheck_period_minutes" env:"API_SELF_HEALTHCHECK_PERIOD_MINUTES"`

	HerokuPort int `env:"PORT"`
}

func (a Api) GetPort() int {
	if a.Port != 0 {
		return a.Port
	} else if a.HerokuPort != 0 {
		return a.HerokuPort
	} else {
		return 8080 // default values
	}
}

func init() {
	err := setLogger()
	if err != nil {
		Logger.Error().Err(err).Msg("failed to set logger")
	}
}

func LoadConfig(configFile string) error {
	err := parseConfigFile(configFile)
	if err != nil {
		return err
	}

	_, err = env.UnmarshalFromEnviron(Config)
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

	return validateConfig()
}

func parseConfigFile(configFile string) error {
	if configFile == "" {
		return nil
	}

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(file, Config)
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
	case "json":
		// use default json formatter
	default:
		return fmt.Errorf("invalid log.format: %s", Config.Log.Format)
	}

	switch strings.Trim(Config.Log.SessionLevel, " ") {
	case "debug":
		Config.Log.SessionLogLevel = discordgo.LogDebug
	case "info":
		Config.Log.SessionLogLevel = discordgo.LogInformational
	case "", "warn":
		Config.Log.SessionLogLevel = discordgo.LogWarning
	case "error":
		Config.Log.SessionLogLevel = discordgo.LogError
	default:
		return fmt.Errorf("invalid log.session_level: %s", Config.Log.SessionLevel)
	}
	return nil
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

func validateConfig() error {
	if Config.Bot.Token == "" {
		return fmt.Errorf("set your env:BOT_TOKEN or yaml:bot.token")
	}
	return nil
}
