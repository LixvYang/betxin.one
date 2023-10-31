package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/api/v1/topic"
	"github.com/lixvyang/betxin.one/api/v1/user"
	"github.com/lixvyang/betxin.one/internal/model/db"
	"github.com/lixvyang/betxin.one/internal/model/redis"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
)

type BetxinHandler struct {
	user.IUserHandler
	topic.ITopicHandler
}

func NewBetxinHandler() *BetxinHandler {
	db := db.NewDatabse()
	rds := redis.NewIRedis()

	return &BetxinHandler{
		user.NewUserHandler(db, rds),
		topic.NewUserHandler(db, rds),
	}
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: errmsg.GetErrMsg(code),
		Data:    data,
	})
}
