package v1

import (
	"github.com/lixvyang/betxin.one/api/v1/user"
	"github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/internal/model/database"
	userz "github.com/lixvyang/betxin.one/internal/model/database/mysql/service/user"
	"github.com/rs/zerolog"
)

type BetxinHandler struct {
	user.IUserHandler
}

func NewBetxinHandler(logger *zerolog.Logger, conf *config.AppConfig) *BetxinHandler {
	db := database.New(logger, conf)
	userz := userz.New(db)

	// categorys, err := db.ListCategories(context.Background())
	// if err != nil {
	// 	panic(err)
	// }
	// categoryMap := make(map[int64]*core.Category)
	// for _, category := range categorys {
	// 	categoryMap[category.ID] = category
	// }

	return &BetxinHandler{
		IUserHandler: user.NewHandler(db, userz),
	}
}
