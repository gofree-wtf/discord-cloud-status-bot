package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/gofree-wtf/discord-cloud-status-bot/pkg/common"
	"net/http"
	"time"
)

type ApiServer struct {
	router *gin.Engine
	server *http.Server
}

func NewApiServer() *ApiServer {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, discord-cloud-status-bot !",
		})
	})

	return &ApiServer{
		router: r,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", Config.Api.Port),
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
		for running {
			Logger.Info().Uint16("port", Config.Api.Port).Msg("run api server")

			err := s.server.ListenAndServe()
			if err == http.ErrServerClosed {
				return
			} else if err != nil {
				Logger.Error().Err(err).Msg("failed to run api server. retry run server")
				time.Sleep(time.Second)
			}
		}
	}()

	if Config.Api.SelfHealthcheckEnabled {
		selfHealthcheckTicker = time.NewTicker(time.Duration(Config.Api.SelfHealthcheckPeriodMinutes) * time.Minute)

		healthcheck := func() {
			resp, err := HttpClient.Get(fmt.Sprintf("http://localhost:%d", Config.Api.Port))
			if err != nil {
				Logger.Error().Err(err).Msg("failed to self healthcheck")
			} else if resp.StatusCode != http.StatusOK {
				Logger.Error().Err(err).Int("responseCode", resp.StatusCode).
					Msg("failed to self healthcheck")
			} else {
				Logger.Info().Msg("success to self healthcheck")
			}
		}

		go func() {
			Logger.Info().Uint32("selfHealthcheckPeriodMinutes", Config.Api.SelfHealthcheckPeriodMinutes).
				Msg("start self healthcheck")

			for range selfHealthcheckTicker.C {
				healthcheck()
			}
		}()
	}

	return func() {
		Logger.Info().Msg("waiting for close api server")

		running = false
		selfHealthcheckTicker.Stop()

		defer Logger.Info().Msg("closed api server")

		err := s.server.Close()
		if err != nil {
			Logger.Warn().Err(err).Msg("failed to close api server")
		}
	}
}
