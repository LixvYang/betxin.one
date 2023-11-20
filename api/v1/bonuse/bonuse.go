package bonuse

import (
	"github.com/lixvyang/betxin.one/internal/model/database"
)

type BonuseHandler struct {
	storage database.IBonuse
}

func NewHandler(db database.IBonuse) *BonuseHandler {
	return &BonuseHandler{
		storage: db,
	}
}

type IBonuseHandler interface {
}
