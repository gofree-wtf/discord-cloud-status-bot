package pkg

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	. "github.com/gofree-wtf/discord-cloud-status-bot/pkg/common"
	"strings"
	"time"
)

const sessionReconnectionPeriod = 10 * time.Second

type BotController struct {
	token           string
	commandPrefix   string
	sessionLogLevel int

	session *discordgo.Session
}

func NewBotController(token, commandPrefix, sessionLogLevel string) (*BotController, error) {
	c := &BotController{
		token:         token,
		commandPrefix: commandPrefix,
	}

	switch strings.Trim(sessionLogLevel, " ") {
	case "debug":
		c.sessionLogLevel = discordgo.LogDebug
	case "info":
		c.sessionLogLevel = discordgo.LogInformational
	case "", "warn":
		c.sessionLogLevel = discordgo.LogWarning
	case "error":
		c.sessionLogLevel = discordgo.LogError
	default:
		return nil, fmt.Errorf("invalid session log level: %s", sessionLogLevel)
	}

	return c, nil
}

func (c *BotController) Run() (shutdownFn func(), err error) {
	stopCh := make(chan struct{}, 1)

	go func() {
		for {
			Logger.Info().Msg("start session")
			err := c.startSession()
			if err != nil {
				Logger.Error().Err(err).Msg("failed to start session. retry start session")
				time.Sleep(sessionReconnectionPeriod)
				continue
			}

			<-stopCh
			break
		}
	}()

	return func() {
		Logger.Info().Msg("waiting for stop session")

		stopCh <- struct{}{}
		c.stopSession()

		Logger.Info().Msg("stopped session")
	}, nil
}

func (c *BotController) startSession() error {
	session, err := discordgo.New("Bot " + c.token)
	if err != nil {
		return err
	}
	session.LogLevel = c.sessionLogLevel

	session.AddHandler(c.handleMessage)
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	c.session = session

	err = session.Open()
	if err != nil {
		return err
	}
	return nil
}

func (c *BotController) stopSession() {
	if c.session == nil {
		return
	}

	err := c.session.Close()
	if err != nil {
		Logger.Warn().Err(err).Msg("failed to close session")
	}
}
