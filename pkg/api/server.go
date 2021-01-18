package api

import (
	"fmt"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	. "github.com/gofree-wtf/discord-cloud-status-bot/pkg/common"
	"net/http"
	"time"
)

type ApiServer struct {
	router *gin.Engine
	server *http.Server
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func NewApiServer() *ApiServer {
	r := gin.New()
	r.Use(
		logger.SetLogger(logger.Config{
			Logger: &Logger,
		}),
		gin.Recovery(),
	)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, discord-cloud-status-bot !",
		})
	})

	return &ApiServer{
		router: r,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", Config.Api.GetPort()),
			Handler: r,
		},
	}
}

func (s *ApiServer) Run() (closeFn func()) {
	var (
		running               = true
		selfHealthcheckTicker *time.Ticker
	)

	go func() {
		l := Logger.With().Int("port", Config.Api.GetPort()).Logger()

		for running {
			l.Info().Msg("run api server")

			err := s.server.ListenAndServe()
			if err == http.ErrServerClosed {
				return
			} else if err != nil {
				l.Error().Err(err).Msg("failed to run api server. retry run server")
				time.Sleep(time.Second)
			}
		}
	}()

	if Config.Api.SelfHealthcheckEnabled {
		selfHealthcheckTicker = time.NewTicker(time.Duration(Config.Api.SelfHealthcheckPeriodMinutes) * time.Minute)
		s.runSelfHealthcheck(selfHealthcheckTicker)
	}

	return func() {
		Logger.Info().Msg("waiting for close api server")

		running = false
		if selfHealthcheckTicker != nil {
			selfHealthcheckTicker.Stop()
		}

		defer Logger.Info().Msg("closed api server")

		err := s.server.Close()
		if err != nil {
			Logger.Warn().Err(err).Msg("failed to close api server")
		}
	}
}

func (s *ApiServer) runSelfHealthcheck(ticker *time.Ticker) {
	l := Logger.With().
		Str("selfHealthcheckUrl", Config.Api.SelfHealthcheckUrl).
		Int("selfHealthcheckPeriodMinutes", Config.Api.SelfHealthcheckPeriodMinutes).
		Logger()

	healthcheck := func() {
		resp, err := HttpClient.Get(Config.Api.SelfHealthcheckUrl)
		if err != nil {
			l.Error().Err(err).Msg("failed to self healthcheck")
		} else if resp.StatusCode != http.StatusOK {
			l.Error().Err(err).Int("responseCode", resp.StatusCode).
				Msg("failed to self healthcheck")
		} else {
			l.Info().Msg("success to self healthcheck")
		}
	}

	go func() {
		l.Info().Msg("start self healthcheck")

		for range ticker.C {
			healthcheck()
		}
	}()
}
