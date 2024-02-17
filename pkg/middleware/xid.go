package middleware

import (
	"encoding/base64"
	"encoding/binary"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/rs/zerolog"
)

func GinXid(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		xid := genReqId()

		log := *logger
		log.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str(consts.DefaultXid, xid)
		})

		c.Header(consts.DefaultXid, xid)
		c.Set(consts.DefaultLoggerKey, log)
		c.Set(consts.DefaultXid, xid)

		c.Next()
	}
}

var (
	pid = uint32(os.Getpid())
)

func genReqId() string {
	var b [12]byte
	binary.LittleEndian.PutUint32(b[:], pid)
	binary.LittleEndian.PutUint64(b[4:], uint64(time.Now().UnixNano()))
	return base64.URLEncoding.EncodeToString(b[:])
}
