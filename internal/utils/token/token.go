package token

import (
	"encoding/base64"

	"github.com/lixvyang/betxin.one/internal/utils/convert"
)

type Token string

type Page struct {
	CreatedAt     int64 `json:"created_at"` // unix毫秒级时间
	NextTimeAtUTC int64 `json:"next_time_at_utc"`
	PageSize      int64 `json:"page_size"`
}

func (p Page) Encode() Token {
	b, err := convert.Marshal(p)
	if err != nil {
		return Token("")
	}
	return Token(base64.StdEncoding.EncodeToString(b))
}

func (t Token) Decode() Page {
	var result Page
	if len(t) == 0 {
		return result
	}
	bytes, err := base64.StdEncoding.DecodeString(string(t))
	if err != nil {
		return result
	}
	err = convert.Unmarshal(bytes, &result)
	if err != nil {
		return result
	}
	return result
}
