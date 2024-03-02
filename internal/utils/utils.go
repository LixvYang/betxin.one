package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/shopspring/decimal"
)

func NewUUID() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}

func VaildUUID(str string) (uuid.UUID, error) {
	return uuid.FromString(str)
}

func GetOrGenUid(c *gin.Context) string {
	xid, ok := c.Get(consts.DefaultXid)
	if !ok {
		xid = NewUUID()
	}
	return xid.(string)
}

func NewDecimalFromString(s string) decimal.Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.NewFromInt(0)
	}
	return d
}

func DecimalAdd(a string, b string) decimal.Decimal {
	return NewDecimalFromString(a).Add(NewDecimalFromString(b))
}
