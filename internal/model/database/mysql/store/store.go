package store

import (
	"fmt"

	"github.com/lixvyang/betxin.one/config"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var defaultHandler *Handler

const (
	OutDirPrefix = "internal/model/database/mysql/store/"
)

type Handler struct {
	*gorm.DB
}

func MustInit(conf *config.AppConfig) *Handler {
	h, err := Init(conf)
	if err != nil {
		panic(err)
	}

	return h
}

func Init(conf *config.AppConfig) (*Handler, error) {
	if defaultHandler != nil {
		return defaultHandler, nil
	}

	var (
		err error
		db  *gorm.DB
	)
	switch conf.Driver {
	case "mysql":

		dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			conf.MySQLConfig.User,
			conf.MySQLConfig.Password,
			conf.MySQLConfig.Host,
			conf.MySQLConfig.Port,
			conf.MySQLConfig.DB,
		)

		// info, err := query.Use(m.db).Topic
		db, err = gorm.Open(mysql.Open(dns), &gorm.Config{})
	default:
		panic("unknown driver")
	}
	if err != nil {
		return nil, err
	}

	defaultHandler = &Handler{
		DB: db,
	}
	return defaultHandler, err
}

type generateModel struct {
	cfg gen.Config
	f   func(g *gen.Generator)
}

var generateModels []*generateModel

func RegistGenerate(cfg gen.Config, f func(g *gen.Generator)) {
	generateModels = append(generateModels, &generateModel{
		cfg: cfg,
		f:   f,
	})
}

func (h *Handler) Generate() {
	for _, gm := range generateModels {
		if gm.cfg.Mode == 0 {
			gm.cfg.Mode = gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface
		}
		g := gen.NewGenerator(gm.cfg)
		g.UseDB(h.DB)
		gm.f(g)
		g.Execute()
	}
}
