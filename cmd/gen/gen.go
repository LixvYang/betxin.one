package gen

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/lixvyang/betxin.one/internal/model/database"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/core"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/store"
	_ "github.com/lixvyang/betxin.one/internal/model/database/mysql/store/category"
	_ "github.com/lixvyang/betxin.one/internal/model/database/mysql/store/collect"
	_ "github.com/lixvyang/betxin.one/internal/model/database/mysql/store/topic"
	_ "github.com/lixvyang/betxin.one/internal/model/database/mysql/store/user"
	"github.com/lixvyang/betxin.one/internal/session"
	"github.com/lixvyang/betxin.one/internal/utils"
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

	cmd.AddCommand(&cobra.Command{
		Use:   "gendata",
		Short: "gen data for betxin",
		RunE: func(cmd *cobra.Command, args []string) error {
			s := session.From(cmd.Context())
			store.MustInit(s.Conf)

			// 生成12个测试数据
			db := database.New(s.Logger, s.Conf)

			var wg sync.WaitGroup

			for i := 0; i < 21; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					time.Sleep(time.Second * time.Duration(i) >> 20)
					topic := &core.Topic{
						Tid:           utils.NewUUID(),
						Cid:           1,
						Title:         fmt.Sprintf("Test title %d", i),
						Content:       "Test content",
						Intro:         "Test intro",
						EndTime:       time.Now().Add(-time.Second),
						RefundEndTime: time.Now().Add(-time.Minute * time.Duration(i)),
					}
					db.CreateTopic(context.Background(), topic)
				}(i)
			}
			wg.Wait()

			return nil
		},
	})

	return cmd
}
