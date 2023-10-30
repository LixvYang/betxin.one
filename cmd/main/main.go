package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/lixvyang/betxin.one/configs"
	"github.com/lixvyang/betxin.one/internal/router"
	"github.com/lixvyang/betxin.one/pkg/logger"
	"github.com/rs/zerolog/log"
)

var (
	configFile string
	signalChan = make(chan os.Signal, 1)
)

func main() {
	flag.StringVar(&configFile, "f", "./configs/configs.yaml", "config file")
	if err := configs.Init(configFile); err != nil {
		log.Error().Err(err).Msgf("[configs.Init] err: %+v", err)
	}

	log.Info().Any("Conf", configs.Conf).Msg("初始化配置成功")
	logger.InitLogger(*configs.Conf.LogConfig)

	srv := router.NewService()

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	signalType := <-signalChan
	signal.Stop(signalChan)

	if err := srv.Shutdown(); err != nil {
		logger.Lg.Fatal().Msgf("Server ShutDown: %+v", err)
	}

	log.Info().Msgf("On Signal: <%s>", signalType)
	log.Info().Msg("Exit command received. Exiting...")
}
