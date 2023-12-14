package store

import (
	"embed"
	"fmt"

	"github.com/lixvyang/betxin.one/config"
	"github.com/pressly/goose/v3"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

var defaultHandler *Handler

const migrationsDir = "migrations"

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

	if err := goose.SetDialect(conf.Driver); err != nil {
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

func (h *Handler) MigrationUp() error {
	db, _ := h.DB.DB()
	goose.SetBaseFS(embedMigrations)
	return goose.Up(db, migrationsDir)
}

func (h *Handler) MigrationUpTo(version int64) error {
	db, _ := h.DB.DB()
	goose.SetBaseFS(embedMigrations)
	return goose.UpTo(db, migrationsDir, version)
}

func (h *Handler) MigrationDown() error {
	db, _ := h.DB.DB()
	goose.SetBaseFS(embedMigrations)
	return goose.Down(db, migrationsDir)
}

func (h *Handler) MigrationDownTo(version int64) error {
	db, _ := h.DB.DB()
	goose.SetBaseFS(embedMigrations)
	return goose.DownTo(db, migrationsDir, version)
}

func (h *Handler) MigrationRedo() error {
	db, _ := h.DB.DB()
	goose.SetBaseFS(embedMigrations)
	return goose.Redo(db, migrationsDir)
}

func (h *Handler) MigrationCreate(name string) error {
	db, _ := h.DB.DB()
	goose.SetBaseFS(nil)
	return goose.Create(db, "store/"+migrationsDir, name, "sql")
}

func (h *Handler) MigrationStatus() error {
	db, _ := h.DB.DB()
	goose.SetBaseFS(embedMigrations)
	return goose.Status(db, migrationsDir)
}

func Transaction(f func(tx *Handler) error) error {
	return defaultHandler.Transaction(func(db *gorm.DB) error {
		return f(&Handler{
			DB: db,
		})
	})
}
