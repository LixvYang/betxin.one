package convert

import (
	"fmt"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

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

func Marshal[T any](j T) (string, error) {
	a, err := jsoniter.Marshal(j)
	if err != nil {
		return "", err
	}
	return string(a), nil
}

func Unmarshal[T any](s string, data T) error {
	return jsoniter.Unmarshal([]byte(s), &data)
}
