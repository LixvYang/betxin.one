package gen

import (
	"github.com/spf13/cobra"

	"github.com/lixvyang/betxin.one/internal/model/database/mysql/store"
	_ "github.com/lixvyang/betxin.one/internal/model/database/mysql/store/user"
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

	return cmd
}
