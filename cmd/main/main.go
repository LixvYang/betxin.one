package main

import (
	"flag"
	"os"

	"github.com/jinzhu/copier"
	"github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/pkg/logger"
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
	logger.New(logConf)

	// if err := snowflake.Init(conf.StartTime, conf.MachineID); err != nil {
	// 	logger.Lg.Panic().Err(err).Msg("[snowflake.Init] err")
	// }

	// srv := router.NewService(conf)

	// signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	// signalType := <-signalChan
	// signal.Stop(signalChan)

	// if err := srv.Shutdown(); err != nil {
	// 	logger.Lg.Fatal().Msgf("Server ShutDown: %+v", err)
	// }

	// log.Info().Msgf("On Signal: <%s>", signalType)
	log.Info().Msg("Exit command received. Exiting...")
}
