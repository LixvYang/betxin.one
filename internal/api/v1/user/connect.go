package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/mongo"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/lixvyang/betxin.one/pkg/jwt"
	"github.com/pandodao/passport-go/auth"
	"github.com/pandodao/passport-go/mvm"
	"github.com/rs/zerolog"
)

type SigninReq struct {
	LoginMethod string `json:"login_method"`
	MixinToken  string `json:"mixin_token"`
	Sign        string `json:"sign"`
	SignedMsg   string `json:"sign_msg"`
}

type UserResp struct {
	UID            string    `json:"uid"`
	IdentityNumber string    `json:"identity_number"`
	FullName       string    `json:"full_name"`
	AvatarURL      string    `json:"avatar_url"`
	Biography      string    `json:"biography"`
	MvmPublicKey   string    `json:"mvm_public_key"`
	CreatedAt      time.Time `json:"created_at"`
}

type SigninResp struct {
	Token string    `json:"token"`
	User  *UserResp `json:"user"`
}

func (uh *UserHandler) Connect(c *gin.Context) {
	logger := c.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)
	var req SigninReq
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Error().Int("errmsg", errmsg.ERROR_BIND).Msgf("bind args error: %+v", err)
		handler.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	logger.Info().Any("req", req).Send()

	err = uh.checkConnectArg(&logger, &req)
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

	var isMvmUser bool
	authUser := new(auth.User)
	switch req.LoginMethod {
	case "mixin_token":
		authUser, err = authorizer.Authorize(c, &auth.AuthorizationParams{
			Method:     auth.AuthMethodMixinToken,
			MixinToken: req.MixinToken,
		})
		if err != nil {
			logger.Error().Err(err).Str("req.LoginMethod", req.LoginMethod).Msg("login_method request user info failed")
			handler.SendResponse(c, errmsg.ERROR_OAUTH, nil)
			return
		}
	case "mvm":
		authUser, err = authorizer.Authorize(ctx, &auth.AuthorizationParams{
			Method:           auth.AuthMethodMvm,
			MvmSignature:     req.Sign,
			MvmSignedMessage: req.SignedMsg,
		})
		if err != nil {
			logger.Error().Err(err).Str("req.LoginMethod", req.LoginMethod).Msg("login_method request user info failed")
			handler.SendResponse(c, errmsg.ERROR_OAUTH, nil)
			return
		}
		isMvmUser = true
		logger.Info().Str("uid", authUser.UserID).Msg("oauth success")
	default:
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	user, jwtToken, err := uh.LoginWithMixin(c, &logger, authUser, isMvmUser)
	if err != nil {
		logger.Error().Err(err).Str("req.LoginMethod", req.LoginMethod).Msg("[LoginWithMixin] failed")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	respUser := new(UserResp)
	copier.Copy(&respUser, &user)

	handler.SendResponse(c, errmsg.SUCCSE, &SigninResp{
		Token: jwtToken,
		User:  respUser,
	})
}

func (uh *UserHandler) checkConnectArg(logger *zerolog.Logger, req *SigninReq) error {
	if req.LoginMethod != "mixin_token" && req.LoginMethod != "mvm" {
		logger.Error().Str("req.LoginMethod", req.LoginMethod).Msg("login_method invaild")
		return errors.New("arg error")
	}
	return nil
}

func (uh *UserHandler) LoginWithMixin(ctx context.Context, logger *zerolog.Logger, authUser *auth.User, isMvmUser bool) (*schema.User, string, error) {
	user, err := uh.loginWithMixin(ctx, logger, authUser, isMvmUser)
	if err != nil {
		return nil, "", err
	}

	jwtToken, err := jwt.GenToken(user.UID)
	if err != nil {
		logger.Error().Err(err).Msgf("[loginMvm][jwt.GenToken] err")
		return nil, "", err
	}

	return user, jwtToken, nil
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

	jwtToken, err := jwt.GenToken(user.UID)
	if err != nil {
		logger.Error().Err(err).Msgf("[loginMvm][jwt.GenToken] err")
		return "", err
	}

	return jwtToken, nil
}

func (uh *UserHandler) loginWithMixin(ctx context.Context, logger *zerolog.Logger, authUser *auth.User, isMvmUser bool) (*schema.User, error) {
	var user = &schema.User{
		UID:            authUser.UserID,
		IdentityNumber: authUser.IdentityNumber,
		FullName:       authUser.FullName,
		AvatarURL:      authUser.AvatarURL,
		MvmPublicKey:   authUser.MvmAddress.Hex(),
		Biography:      authUser.Biography,
		IsMvmUser:      isMvmUser,
		MixinCreatedAt: authUser.CreatedAt,
	}

	existing, err := uh.userSrv.GetUserByUid(ctx, logger, user.UID)
	if err != nil && err != mongo.ErrNoSuchUser {
		logger.Error().Err(err).Msgf("[LoginWithMixin][GetUserByUid] err")
		return nil, err
	}

	// create
	if err == mongo.ErrNoSuchUser {
		user.CreatedAt = time.Now()
		user.UpdatedAt = user.CreatedAt
		err = uh.userSrv.CreateUser(ctx, logger, user)
		if err != nil && err != mongo.ErrUserExist {
			logger.Error().Err(err).Msgf("[LoginWithMixin][CreateUser] err")
			return nil, err
		}
		return user, nil
	}

	user.UpdatedAt = time.Now()
	// update
	err = uh.userSrv.UpdateUser(ctx, logger, existing.UID, user)
	if err != nil {
		fmt.Printf("err users.Updates: %v\n", err)
		return nil, err
	}

	newUser, err := uh.userSrv.GetUserByUid(ctx, logger, user.UID)
	if err != nil {
		logger.Error().Err(err).Msgf("[LoginWithMixin][GetUserByUid] err")
		return nil, err
	}

	return newUser, nil
}
