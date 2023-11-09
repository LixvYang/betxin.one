package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/jinzhu/copier"
	"github.com/lixvyang/betxin.one/configs"
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
	flag.StringVar(&configFile, "f", "./configs/configs.yaml", "config file")
	conf := &configs.AppConfig{}
	if err := configs.Init(configFile, conf); err != nil {
		log.Error().Err(err).Msgf("[configs.Init] err: %+v", err)
	}
	log.Info().Any("conf", conf).Msg("init config succes")

	logConf := new(logger.LogConfig)
	copier.Copy(logConf, conf.LogConfig)
	logger.InitLogger(logConf)

	if err := snowflake.Init(conf.StartTime, conf.MachineID); err != nil {
		logger.Lg.Panic().Err(err).Msg("[snowflake.Init] err")
		panic(err)
	}

	srv := router.NewService(conf)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	signalType := <-signalChan
	signal.Stop(signalChan)

	if err := srv.Shutdown(); err != nil {
		logger.Lg.Fatal().Msgf("Server ShutDown: %+v", err)
	}

	log.Info().Msgf("On Signal: <%s>", signalType)
	log.Info().Msg("Exit command received. Exiting...")
}
