package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/rs/zerolog"
)

type GetUserResp struct {
	UID            string `json:"uid"`
	IdentityNumber string `json:"identity_number"`
	FullName       string `json:"full_name"`
	AvatarURL      string `json:"avatar_url"`
	Biography      string `json:"biography"`
	ClientID       string `json:"client_id"`
	IsMvmUser      bool   `json:"is_mvm_user"`
}

func (uh *UserHandler) Get(c *gin.Context) {
	logger := c.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)
	uid := c.GetString("uid")
	if uid == "" {
		logger.Error().Msg("check args error: uid is \"\"")
		handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
		return
	}

	logger.Info().Str("uid", uid).Send()

	user, err := uh.userSrv.GetUserByUid(c, &logger, uid)
	if err != nil {
		logger.Error().Err(err).Msg("[uh.Get][GetUserByUid] err")
		handler.SendResponse(c, errmsg.ERROR_GET_USER, nil)
		return
	}

	resp := new(GetUserResp)
	copier.Copy(resp, user)
	handler.SendResponse(c, errmsg.SUCCSE, resp)
}
