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

func (c *BotController) Run() (closeFn func(), err error) {
	stopCh := make(chan struct{}, 1)

	go func() {
		for {
			Logger.Info().Msg("connect session")
			err := c.connect()
			if err != nil {
				Logger.Error().Err(err).Msg("failed to connect session. retry connect session")
				time.Sleep(sessionReconnectionPeriod)
				continue
			}

			<-stopCh
			break
		}
	}()

	return func() {
		Logger.Info().Msg("waiting for close session")

		stopCh <- struct{}{}
		c.close()

		Logger.Info().Msg("closed session")
	}, nil
}

func (c *BotController) connect() error {
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

func (c *BotController) close() {
	if c.session == nil {
		return
	}

	err := c.session.Close()
	if err != nil {
		Logger.Warn().Err(err).Msg("failed to close session")
	}
}
