package category

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/database"
)

type CategoryHandler struct {
	storage database.ICategoty
}

func NewHandler(db database.ICategoty) *CategoryHandler {
	return &CategoryHandler{
		storage: db,
	}
}

type ICategoryHandler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
	Delete(c *gin.Context)
	Update(c *gin.Context)
}
