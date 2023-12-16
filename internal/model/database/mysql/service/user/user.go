package userz

import (
	"context"
	"fmt"

	"github.com/lixvyang/betxin.one/internal/model/database/mysql/core"
	"github.com/pandodao/passport-go/auth"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Config struct {
	MixinClientSecret string
}

type UserService struct {
	// client *mixin.Client
	user core.UserStore
	// cfg    Config
}

func New(
	// client *mixin.Client,
	user core.UserStore,
	// cfg Config,
) *UserService {
	return &UserService{
		// client: client,
		user: user,
		// cfg:    cfg,
	}
}

func (s *UserService) LoginWithMixin(ctx context.Context, logger *zerolog.Logger, authUser *auth.User, isMvmUser bool) (*core.User, error) {
	var user = &core.User{
		UID:            authUser.UserID,
		IdentityNumber: authUser.IdentityNumber,
		FullName:       authUser.FullName,
		AvatarURL:      authUser.AvatarURL,
		MvmPublicKey:   authUser.MvmAddress.Hex(),
		Biography:      authUser.Biography,
		IsMvmUser:      isMvmUser,
	}

	existing, err := s.user.GetUserByUid(ctx, user.UID)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error().Err(err).Msgf("[LoginWithMixin][GetUserByUid] err")
		return nil, err
	}

	// create
	if err == gorm.ErrRecordNotFound {
		err = s.user.CreateUser(ctx, user)
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.Error().Err(err).Msgf("[LoginWithMixin][CreateUser] err")
			return nil, err
		}
		return user, nil
	}

	// update
	err = s.user.UpdateUserInfo(ctx, existing.UID, user)
	if err != nil {
		fmt.Printf("err users.Updates: %v\n", err)
		return nil, err
	}

	newUser, err := s.user.GetUserByUid(ctx, user.UID)
	if err != nil {
		logger.Error().Err(err).Msgf("[LoginWithMixin][GetUserByUid] err")
		return nil, err
	}

	return newUser, nil
}
