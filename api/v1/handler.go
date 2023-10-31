package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/db"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
)

type BetxinHandler struct {
	db *db.Database
}

func NewH()  {}


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
