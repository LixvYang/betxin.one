package user

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	v1 "github.com/lixvyang/betxin.one/api/v1"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/model"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/lixvyang/betxin.one/pkg/jwt"
	"github.com/pandodao/passport-go/auth"
	"github.com/pandodao/passport-go/eip4361"
	"github.com/pandodao/passport-go/mvm"
	"github.com/rs/zerolog"
)

type SigninReq struct {
	LoginMethod string `json:"login_method"`
	MixinToken  string `json:"mixin_token"`
	Sign        string `json:"sign"`
	SignedMsg   string `json:"sign_msg"`
}

type SigninResp struct {
	Token string `json:"token"`
}

func (uh *UserHandler) MixinOauth(c *gin.Context) {
	logger := c.MustGet(consts.LoggerKey).(*zerolog.Logger)
	var req SigninReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error().Int("errmsg", errmsg.ERROR_BIND).Msgf("bind args error: %+v", err)
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	if req.LoginMethod != "mixin_token" && req.LoginMethod != "mvm" {
		logger.Error().Str("req.LoginMethod", req.LoginMethod).Msg("login_method invaild")
		v1.SendResponse(c, errmsg.ERROR_OAUTH, nil)
		return
	}

	// 先在自己的数据里找
	// uh.db.GetUserByUid(req.)

	authorizer := auth.New([]string{
		"30aad5a5-e5f3-4824-9409-c2ff4152724e",
	}, []string{
		"localhost:4000",
		"localhost:4000/*",
	})

	switch req.LoginMethod {
	case "mixin_token":
		// 2. 访问用户信息
		userInfo, err := authorizer.Authorize(c, &auth.AuthorizationParams{
			Method:     auth.AuthMethodMixinToken,
			MixinToken: req.MixinToken,
		})
		if err != nil {
			logger.Error().Str("req.LoginMethod", req.LoginMethod).Msg("login_method invaild")
			v1.SendResponse(c, errmsg.ERROR_OAUTH, nil)
			return
		}

		user := &model.User{
			AvatarUrl:      userInfo.AvatarURL,
			IdentityNumber: userInfo.IdentityNumber,
			SessionId:      userInfo.SessionID,
			Uid:            userInfo.UserID,
			FullName:       userInfo.FullName,
			Biography:      userInfo.Biography,
		}

		code := uh.db.CheckUser(userInfo.UserID)
		if code != errmsg.SUCCSE {
			if code = uh.db.CreateUser(user); err != nil {
				logger.Error().Int("code", code).Msg("[CheckUser][CreateUser] error")
				v1.SendResponse(c, errmsg.ERROR_GET_USER, nil)
				return
			}
		} else {
			code = uh.db.UpdateUser(user)
			if code != errmsg.SUCCSE {
				logger.Error().Int("code", code).Msg("[CheckUser][UpdateUser] error")
			}
		}

		jwtToken, err := jwt.GenToken(user.Uid)
		if err != nil {
			logger.Error().Err(err).Str("jwtToken", jwtToken).Int("code", code).Msg("[CheckUser][UpdateUser] error")
			v1.SendResponse(c, errmsg.SUCCSE, nil)
			return
		}

		logger.Info().Str("uid", user.Uid).Msg("oauth success")

		v1.SendResponse(c, errmsg.SUCCSE, &SigninResp{
			Token: jwtToken,
		})
		return
	case "mvm":
		message, err := eip4361.Parse(req.SignedMsg)
		if err != nil {
			v1.SendResponse(c, errmsg.ERROR_AUTH, nil)
			return
		}

		if err := message.Validate(time.Now()); err != nil {
			v1.SendResponse(c, errmsg.ERROR_AUTH, nil)
			return
		}

		if err := eip4361.Verify(message, req.Sign); err != nil {
			v1.SendResponse(c, errmsg.ERROR_AUTH, nil)
			return
		}

		// get the public key from the message, and use it to login
		jwtToken, code := uh.loginMvm(c, logger, message.Address)
		if code != errmsg.SUCCSE {
			logger.Error().Msgf("[uh.loginMvm] error")
			v1.SendResponse(c, errmsg.ERROR, nil)
			return
		}

		v1.SendResponse(c, errmsg.SUCCSE, &SigninResp{
			Token: jwtToken,
		})
		return
	default:
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}
}

func (uh *UserHandler) loginMvm(c *gin.Context, logger *zerolog.Logger, pubkey string) (string, int) {
	addr := common.HexToAddress(pubkey)
	mvmUser, err := mvm.GetBridgeUser(c, addr)
	if err != nil {
		logger.Error().Err(err).Msgf("[loginMvm][mvm.GetBridgeUser] err")
		return "", errmsg.ERROR
	}

	contractAddr, err := mvm.GetUserContract(c, mvmUser.UserID)
	if err != nil {
		logger.Error().Err(err).Msgf("[loginMvm][mvm.GetUserContract] err")
		fmt.Printf("err mvm.GetUserContract: %v\n", err)
		return "", errmsg.ERROR
	}

	// if contractAddr is not 0x000..00, it means the user has already registered a mvm account
	emptyAddr := common.Address{}
	if contractAddr != emptyAddr {
		return "", errmsg.ERROR
	}

	user := &model.User{
		IsMvmUser:  1,
		Uid:        mvmUser.UserID,
		FullName:   mvmUser.FullName,
		Contract:   mvmUser.Contract,
		PrivateKey: mvmUser.Key.PrivateKey,
		ClientId:   mvmUser.Key.ClientID,
		SessionId:  mvmUser.Key.SessionID,
	}
	logger.Info().Any("user", user).Msg("userInfo")
	code := uh.db.CreateUser(user)
	if code != errmsg.SUCCSE {
		logger.Error().Msgf("[loginMvm][db.CreateUser(user)] err")
		return "", errmsg.ERROR
	}

	jwtToken, err := jwt.GenToken(user.Uid)
	if err != nil {
		logger.Error().Err(err).Msgf("[loginMvm][jwt.GenToken] err")
		return "", errmsg.ERROR
	}

	return jwtToken, errmsg.SUCCSE
}
