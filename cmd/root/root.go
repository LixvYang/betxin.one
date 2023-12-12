package root

import (
	"github.com/lixvyang/betxin.one/cmd/httpd"
	"github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/internal/session"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func NewCmdRoot(version string) *cobra.Command {
	var opt struct {
		configFile string
	}

	cmd := &cobra.Command{
		Use:           "betxin <command> <subcommand> [flags]",
		Short:         "gb",
		Long:          `A boilerplate for go programe.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			s := session.From(cmd.Context())

			conf := new(config.AppConfig)
			if err := config.Init(opt.configFile, conf); err != nil {
				log.Error().Err(err).Msgf("[configs.Init] err: %+v", err)
			}
			s.WithConf(conf)
			s.WithLogger(conf.LogConfig)

			log.Info().Any("conf", conf).Msg("init config succes")
			return nil
		},
	}

	// load config

	cmd.PersistentFlags().StringVarP(&opt.configFile, "file", "f", "./config/config.yaml", "config file path")

	cmd.AddCommand(httpd.NewCmdHttpd())
	// cmd.AddCommand(wss.NewCmdWss())
	// cmd.AddCommand(echo.NewCmdEcho())
	// cmd.AddCommand(migrate.NewCmdMigrate())
	// cmd.AddCommand(worker.NewCmdWorker())
	// cmd.AddCommand(gen.NewCmdGen())

	return cmd
}
