package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	assert := assert.New(t)
	now := time.Now()
	p := Page{
		CreatedAt:     now,
		NextTimeAtUTC: time.Now().UnixMilli(),
		PageSize:      10,
	}
	encodeStr := p.Encode()
	newPage := encodeStr.Decode()
	assert.True(newPage.CreatedAt.Equal(p.CreatedAt))
	assert.Equal(newPage.NextTimeAtUTC, p.NextTimeAtUTC)
	assert.Equal(newPage.PageSize, p.PageSize)
}
