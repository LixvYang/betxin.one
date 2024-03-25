package topicpurchase

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/database/mongo"
)

type TopicPurchaseHandler struct {
	storage *mongo.MongoService
}

func NewHandler(db *mongo.MongoService) *TopicPurchaseHandler {
	tp := &TopicPurchaseHandler{
		storage: db,
	}
	return tp
}

type ITopicPurchaseHandler interface {
	Create(*gin.Context)
}
