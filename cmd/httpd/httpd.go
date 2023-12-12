package httpd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/lixvyang/betxin.one/internal/router"
	"github.com/lixvyang/betxin.one/internal/session"
	"github.com/spf13/cobra"
)

func NewCmdHttpd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "httpd [port]",
		Short: "start the httpd daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			s := session.From(cmd.Context())
			srv := router.NewService(s.Logger, s.Conf)

			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
			signalType := <-signalChan
			signal.Stop(signalChan)

			if err := srv.Shutdown(); err != nil {
				s.Logger.Fatal().Msgf("Server ShutDown: %+v", err)
			}

			s.Logger.Info().Msgf("On Signal: <%s>", signalType)
			s.Logger.Info().Msg("Exit command received. Exiting...")
			return nil
		},
	}

	return cmd
}
