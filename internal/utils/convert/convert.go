package convert

import (
	"fmt"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func StrToInt64(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func IntToStr(t any) string {
	return fmt.Sprintf("%v", t)
}

func Marshal[T any](j T) (string, error) {
	a, err := json.Marshal(j)
	if err != nil {
		return "", err
	}
	return string(a), nil
}

func Unmarshal[T any](s string, data T) error {
	return json.Unmarshal([]byte(s), &data)
}
