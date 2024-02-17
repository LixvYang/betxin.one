package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/jinzhu/copier"
	"github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/internal/router"
	"github.com/lixvyang/betxin.one/pkg/logger"
	"github.com/lixvyang/betxin.one/pkg/snowflake"
	"github.com/rs/zerolog/log"
)

var (
	configFile string
	signalChan = make(chan os.Signal, 1)
)

func main() {
	flag.StringVar(&configFile, "f", "./config/config.yaml", "config file")
	conf := new(config.AppConfig)
	if err := config.Init(configFile, conf); err != nil {
		log.Error().Err(err).Msgf("[configs.Init] err: %+v", err)
	}
	log.Info().Any("conf", conf).Msg("init config succes")

	logConf := new(logger.LogConfig)
	copier.Copy(logConf, conf.LogConfig)
	logger := logger.New(logConf)

	if err := snowflake.Init(conf.StartTime, conf.MachineID); err != nil {
		logger.Panic().Err(err).Msg("[snowflake.Init] err")
	}

	srv := router.NewService(logger, conf)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	signalType := <-signalChan
	signal.Stop(signalChan)

	if err := srv.Shutdown(); err != nil {
		logger.Fatal().Msgf("Server ShutDown: %+v", err)
	}

	logger.Info().Msgf("On Signal: <%s>", signalType)
	logger.Info().Msg("Exit command received. Exiting...")
}
