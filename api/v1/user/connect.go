package user

import (
	"context"
	"time"

	"github.com/lixvyang/betxin.one/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/lixvyang/betxin.one/pkg/jwt"
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
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

func (uh *UserHandler) Connect(c *gin.Context) {
	logger := c.MustGet(consts.LoggerKey).(*zerolog.Logger)
	var req SigninReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error().Int("errmsg", errmsg.ERROR_BIND).Msgf("bind args error: %+v", err)
		handler.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	logger.Info().Any("req", req).Send()

	err := uh.checkConnectArg(logger, &req)
	if err != nil {
		logger.Error().Int("errmsg", errmsg.ERROR_INVAILD_ARGV).Msgf("check args error: %+v", err)
		handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	authorizer := auth.New([]string{
		"30aad5a5-e5f3-4824-9409-c2ff4152724e",
	}, []string{
		"localhost:3000",
		"localhost:3000/*",
	})

	switch req.LoginMethod {
	case "mixin_token":
		userInfo, err := authorizer.Authorize(c, &auth.AuthorizationParams{
			Method:     auth.AuthMethodMixinToken,
			MixinToken: req.MixinToken,
		})
		if err != nil {
			logger.Error().Err(err).Str("req.LoginMethod", req.LoginMethod).Msg("login_method request user info failed")
			handler.SendResponse(c, errmsg.ERROR_OAUTH, nil)
			return
		}

		user := &schema.User{
			AvatarURL:      userInfo.AvatarURL,
			IdentityNumber: userInfo.IdentityNumber,
			SessionID:      userInfo.SessionID,
			UID:            userInfo.UserID,
			FullName:       userInfo.FullName,
			Biography:      userInfo.Biography,
		}

		err = uh.storage.CheckUser(ctx, logger, userInfo.UserID)
		if err != nil {
			err = uh.storage.CreateUser(ctx, logger, user)
			if err != nil {
				logger.Error().Err(err).Msg("[CheckUser][CreateUser] error")
				handler.SendResponse(c, errmsg.ERROR_GET_USER, nil)
				return
			}
		} else {
			err = uh.storage.UpdateUser(ctx, logger, user)
			if err != nil {
				logger.Error().Err(err).Msg("[CheckUser][UpdateUser] error")
			}
		}

		jwtToken, err := jwt.GenToken(user.UID)
		if err != nil {
			logger.Error().Err(err).Str("jwtToken", jwtToken).Err(err).Msg("[CheckUser][UpdateUser] error")
			handler.SendResponse(c, errmsg.SUCCSE, nil)
			return
		}

		logger.Info().Str("uid", user.UID).Msg("oauth success")

		handler.SendResponse(c, errmsg.SUCCSE, &SigninResp{
			Token: jwtToken,
		})
		return
	case "mvm":
		message, err := eip4361.Parse(req.SignedMsg)
		if err != nil {
			handler.SendResponse(c, errmsg.ERROR_AUTH, nil)
			return
		}

		if err := message.Validate(time.Now()); err != nil {
			handler.SendResponse(c, errmsg.ERROR_AUTH, nil)
			return
		}

		if err := eip4361.Verify(message, req.Sign); err != nil {
			handler.SendResponse(c, errmsg.ERROR_AUTH, nil)
			return
		}

		// get the public key from the message, and use it to login
		jwtToken, err := uh.loginMvm(c, logger, message.Address)
		if err != nil {
			logger.Error().Err(err).Msgf("[uh.loginMvm] error")
			handler.SendResponse(c, errmsg.ERROR, nil)
			return
		}

		handler.SendResponse(c, errmsg.SUCCSE, &SigninResp{
			Token: jwtToken,
		})
		return
	default:
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}
}

func (uh *UserHandler) loginMvm(c *gin.Context, logger *zerolog.Logger, pubkey string) (string, error) {
	addr := common.HexToAddress(pubkey)
	mvmUser, err := mvm.GetBridgeUser(c, addr)
	if err != nil {
		logger.Error().Err(err).Msgf("[loginMvm][mvm.GetBridgeUser] err")
		return "", err
	}

	contractAddr, err := mvm.GetUserContract(c, mvmUser.UserID)
	if err != nil {
		logger.Error().Err(err).Msgf("[loginMvm][mvm.GetUserContract] err")
		return "", err
	}

	logger.Info().Any("contractAddr", contractAddr).Msg("userInfo")

	// if contractAddr is not 0x000..00, it means the user has already registered a mvm account
	emptyAddr := common.Address{}
	if contractAddr == emptyAddr {
		logger.Error().Err(err).Msgf("[loginMvm][mvm.emptyAddr] err")
		return "", err
	}

	logger.Info().Msgf("user: %+v", mvmUser)

	user := &schema.User{
		IsMvmUser:  true,
		UID:        mvmUser.UserID,
		FullName:   mvmUser.FullName,
		Contract:   mvmUser.Contract,
		PrivateKey: mvmUser.Key.PrivateKey,
		ClientID:   mvmUser.Key.ClientID,
		SessionID:  mvmUser.Key.SessionID,
	}
	logger.Info().Any("user", user).Msg("userInfo")
	// err = uh.storage.CreateUser(context.Background(), logger, user)
	// if err != nil {
	// 	logger.Error().Err(err).Msgf("[loginMvm][db.CreateUser(user)] err")
	// 	return "", err
	// }

	err = uh.storage.CheckUser(c, logger, user.UID)
	if err != nil {
		err = uh.storage.CreateUser(c, logger, user)
		if err != nil {
			logger.Error().Err(err).Msg("[CheckUser][CreateUser] error")
			handler.SendResponse(c, errmsg.ERROR_GET_USER, nil)
			return "", errors.New("create user err")
		}
	} else {
		err = uh.storage.UpdateUser(c, logger, user)
		if err != nil {
			logger.Error().Err(err).Msg("[CheckUser][UpdateUser] error")
		}
	}

	jwtToken, err := jwt.GenToken(user.UID)
	if err != nil {
		logger.Error().Err(err).Msgf("[loginMvm][jwt.GenToken] err")
		return "", err
	}

	return jwtToken, nil
}

func (uh *UserHandler) checkConnectArg(logger *zerolog.Logger, req *SigninReq) error {
	if req.LoginMethod != "mixin_token" && req.LoginMethod != "mvm" {
		logger.Error().Str("req.LoginMethod", req.LoginMethod).Msg("login_method invaild")
		return errors.New("arg error")
	}
	return nil
}
