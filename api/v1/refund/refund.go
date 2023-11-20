package refund

import (
	"github.com/lixvyang/betxin.one/internal/model/database"
)

type RefundHandler struct {
	storage database.IBonuse
}

func NewHandler(db database.IBonuse) *RefundHandler {
	return &RefundHandler{
		storage: db,
	}
}

type IRefundHandler interface {
}
