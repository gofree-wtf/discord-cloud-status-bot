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
	stopFlag := false

	go func() {
		for {
			Logger.Info().Uint16("port", Config.Api.Port).Msg("run api server")

			err := s.server.ListenAndServe()
			if err == http.ErrServerClosed || stopFlag {
				return
			} else if err != nil {
				Logger.Error().Err(err).Msg("failed to run api server. retry run server")
				time.Sleep(time.Second)
			}
		}
	}()

	return func() {
		Logger.Info().Msg("waiting for close api server")

		stopFlag = true

		defer Logger.Info().Msg("closed api server")

		err := s.server.Close()
		if err != nil {
			Logger.Warn().Err(err).Msg("failed to close api server")
		}
	}
}
