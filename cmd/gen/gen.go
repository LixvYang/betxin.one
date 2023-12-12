package gen

import (
	"fmt"

	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/lixvyang/betxin.one/internal/model/database/mysql/store"
	_ "github.com/lixvyang/betxin.one/internal/model/database/mysql/store/user"
	"github.com/lixvyang/betxin.one/internal/session"
)

func NewCmdGen() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "generate database operation code",
		// RunE: func(cmd *cobra.Command, args []string) error {
		// 	outPath := os.Getenv("BETXIN_DB_OUTPATH")
		// 	if outPath == "" {
		// 		outPath = defaultOutPath
		// 	}
		// 	s := session.From(cmd.Context())

		// 	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		// 		s.Conf.MySQLConfig.User,
		// 		s.Conf.MySQLConfig.Password,
		// 		s.Conf.MySQLConfig.Host,
		// 		s.Conf.MySQLConfig.Port,
		// 		s.Conf.MySQLConfig.DB,
		// 	)

		// 	// 指定生成代码的具体相对目录(相对当前文件)，默认为：./query
		// 	// 默认生成需要使用WithContext之后才可以查询的代码，但可以通过设置gen.WithoutContext禁用该模式
		// 	g := gen.NewGenerator(gen.Config{
		// 		ModelPkgPath: "sqlmodel",
		// 		// 默认会在 OutPath 目录生成CRUD代码，并且同目录下生成 model 包
		// 		// 所以OutPath最终package不能设置为model，在有数据库表同步的情况下会产生冲突
		// 		// 若一定要使用可以通过ModelPkgPath单独指定model package的名称
		// 		OutPath: outPath,
		// 		/* ModelPkgPath: "dal/model"*/

		// 		// gen.WithoutContext：禁用WithContext模式
		// 		// gen.WithDefaultQuery：生成一个全局Query对象Q
		// 		// gen.WithQueryInterface：生成Query接口
		// 		Mode: gen.WithDefaultQuery | gen.WithQueryInterface,
		// 	})

		// 	// 通常复用项目中已有的SQL连接配置db(*gorm.DB)
		// 	// 非必需，但如果需要复用连接时的gorm.Config或需要连接数据库同步表信息则必须设置
		// 	g.UseDB(connectDB(dns))

		// 	// 从连接的数据库为所有表生成Model结构体和CRUD代码
		// 	// 也可以手动指定需要生成代码的数据表
		// 	g.ApplyBasic(g.GenerateAllTable()...)

		// 	// 执行并生成代码
		// 	g.Execute()

		// 	return nil
		// },

		Run: func(cmd *cobra.Command, args []string) {
			s := session.From(cmd.Context())
			h := store.MustInit(s.Conf)
			h.Generate()
		},
	}

	return cmd
}

func connectDB(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(fmt.Errorf("connect db fail: %w", err))
	}
	return db
}

const (
	defaultOutPath = "./internal/model/database/mysql/dal/query"
)
