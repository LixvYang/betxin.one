package convert

import (
	"fmt"
	"strconv"

	"github.com/gofrs/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
)

func init() {
	decimal.DivisionPrecision = 2
}

func StrToInt64(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func IntToStr(t int64) string {
	return fmt.Sprintf("%d", t)
}

func Marshal[T any](j T) ([]byte, error) {
	a, err := jsoniter.Marshal(j)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func Unmarshal[T any](b []byte, data T) error {
	return jsoniter.Unmarshal(b, &data)
}

func NewUUID() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}

func VaildUUID(str string) (uuid.UUID, error) {
	return uuid.FromString(str)
}

func NewDecimalFromString(s string) decimal.Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.NewFromInt(0)
	}
	return d
}

func DecimalAdd(a string, b string) decimal.Decimal {
	decimal.DivisionPrecision = 8
	return NewDecimalFromString(a).Add(NewDecimalFromString(b))
}

func DecimalDiv(a string, b string) decimal.Decimal {
	if NewDecimalFromString(b).IsZero() {
		return decimal.NewFromInt(0)
	}

	return NewDecimalFromString(a).Div(NewDecimalFromString(b))
}
