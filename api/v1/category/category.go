package category

import (
	"github.com/lixvyang/betxin.one/internal/model/database"
)

type CategoryHandler struct {
	storage database.ICategoty
}

func NewHandler(db database.ICategoty) ICategoryHandler {
	return &CategoryHandler{
		storage: db,
	}
}

type ICategoryHandler interface {
}
