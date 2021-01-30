package pkg

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gofree-wtf/discord-cloud-status-bot/pkg/command"
	. "github.com/gofree-wtf/discord-cloud-status-bot/pkg/common"
	"time"
)

const sessionReconnectionPeriod = 10 * time.Second

type BotController struct {
	commandPrefix string
	command       *command.BotCommand
	session       *discordgo.Session
}

func NewBotController(token, commandPrefix string, sessionLogLevel int) (*BotController, error) {
	c := &BotController{
		commandPrefix: commandPrefix,
		command:       command.NewBotCommand(commandPrefix),
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	session.LogLevel = sessionLogLevel
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	session.AddHandler(c.handleMessage)

	c.session = session
	return c, nil
}

func (c *BotController) Run() (closeFn func()) {
	stopCh := make(chan struct{}, 1)

	go func() {
		for {
			Logger.Info().Msg("connect bot session")

			err := c.session.Open()
			if err != nil {
				Logger.Error().Err(err).Msg("failed to connect bot session. retry connect session")
				time.Sleep(sessionReconnectionPeriod)
				continue
			}

			<-stopCh
			break
		}
	}()

	return func() {
		Logger.Info().Msg("waiting for close bot session")

		stopCh <- struct{}{}

		defer Logger.Info().Msg("closed bot session")

		if c.session == nil {
			return
		}

		err := c.session.Close()
		if err != nil {
			Logger.Warn().Err(err).Msg("failed to close bot session")
		}
	}
}
