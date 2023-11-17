package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	assert := assert.New(t)
	p := Page{
		PreID:         1233,
		NextTimeAtUTC: time.Now().UnixMilli(),
		PageSize:      10,
	}
	encodeStr := p.Encode()
	newPage := encodeStr.Decode()
	assert.Equal(newPage.PreID, p.PreID)
	assert.Equal(newPage.NextTimeAtUTC, p.NextTimeAtUTC)
	assert.Equal(newPage.PageSize, p.PageSize)
}
