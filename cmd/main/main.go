package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/lixvyang/betxin.one/configs"
)

var (
	signalChan = make(chan os.Signal, 1)
	configFile = flag.String("f", "./configs/configs.yaml", "config file")
)

func main() {
	flag.Parse()
	if err := configs.Init(*configFile); err != nil {
		log.Error().Err(err).Msgf("[configs.Init] err: %+v", err)
	}

	log.Info().Any("Conf", configs.Conf).Send()

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	signalType := <-signalChan
	signal.Stop(signalChan)
	log.Info().Msgf("On Signal: <%s>", signalType)
	log.Info().Msg("Exit command received. Exiting...")
}
