package common

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

const defaultBotCommandPrefix = "!cs"

var Logger = log.With().Caller().Logger()

type Config struct {
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

func ParseConfigFile(path string) (*Config, error) {
	c := &Config{}

	// parse config file
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, c)
	if err != nil {
		return nil, err
	}

	// config log module
	err = c.setLogLevel()
	if err != nil {
		return nil, err
	}

	err = c.setLogFormatter()
	if err != nil {
		return nil, err
	}

	// validate values
	err = c.validateBotValues()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) setLogLevel() error {
	logLevel := c.Log.Level
	if logLevel == "" {
		logLevel = zerolog.InfoLevel.String()
	}

	lv, err := zerolog.ParseLevel(c.Log.Level)
	if err != nil {
		Logger.Error().Err(err).Str("log.level", c.Log.Level).Msg("failed to parse log level")
		return err
	}

	zerolog.SetGlobalLevel(lv)
	return nil
}

func (c *Config) setLogFormatter() error {
	switch c.Log.Format {
	case "console":
		Logger = Logger.Output(zerolog.NewConsoleWriter())
		return nil
	case "json":
		// use default json formatter
		return nil
	default:
		return fmt.Errorf("invalid log.format: %s", c.Log.Format)
	}
}

func (c *Config) validateBotValues() error {
	if c.Bot.Token == "" {
		return fmt.Errorf("set your bot.token")
	}

	commandPrefix := strings.Trim(c.Bot.CommandPrefix, " ")
	if commandPrefix == "" {
		Logger.Warn().Msgf("empty bot.command_prefix. use default '%s'", defaultBotCommandPrefix)
		commandPrefix = defaultBotCommandPrefix
	}
	c.Bot.CommandPrefix = commandPrefix

	return nil
}
