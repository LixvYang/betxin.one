package gen

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/lixvyang/betxin.one/internal/model/database/mysql/store"
	_ "github.com/lixvyang/betxin.one/internal/model/database/mysql/store/user"
	_ "github.com/lixvyang/betxin.one/internal/model/database/mysql/store/topic"
	"github.com/lixvyang/betxin.one/internal/session"
)

func NewCmdGen() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "generate database operation code",

		Run: func(cmd *cobra.Command, args []string) {
			s := session.From(cmd.Context())
			h := store.MustInit(s.Conf)
			h.Generate()
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "Migrate the DB to the most recent version available",
		RunE: func(cmd *cobra.Command, args []string) error {
			s := session.From(cmd.Context())
			h := store.MustInit(s.Conf)

			return h.MigrationUp()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "up-to VERSION",
		Short: "Migrate the DB to a specific VERSION",
		RunE: func(cmd *cobra.Command, args []string) error {
			s := session.From(cmd.Context())
			h := store.MustInit(s.Conf)
			if len(args) != 1 {
				return fmt.Errorf("up-to requires a version argument")
			}
			version, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid version: %w", err)
			}
			return h.MigrationUpTo(version)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "Roll back the version by 1",
		RunE: func(cmd *cobra.Command, args []string) error {
			s := session.From(cmd.Context())
			h := store.MustInit(s.Conf)

			return h.MigrationDown()
		},
	})

	return cmd
}
