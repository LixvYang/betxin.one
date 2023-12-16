package topicpurchase

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/database"
)

type TopicPurchaseHandler struct {
	topic          database.ITopic
	topicPurchaase database.ITopicPurchase
}

func NewHandler(db database.Database) ITopicPurchaseHandler {
	return nil
	// return &TopicPurchaseHandler{db, db}
}

type ITopicPurchaseHandler interface {
	Create(*gin.Context)
}
