package user

import (
	"context"
	"time"

	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/query"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/sqlmodel"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/lixvyang/betxin.one/internal/utils/convert"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/jinzhu/copier"
)

type UserModel struct {
	db    *query.Query
	cache *cache.Cache
}

func NewUserModel(query *query.Query, cache *cache.Cache) UserModel {
	return UserModel{
		db:    query,
		cache: cache,
	}
}

// CheckUser 查询用户是否存在
func (um *UserModel) CheckUser(ctx context.Context, logger *zerolog.Logger, uid string) (err error) {
	_, err = um.GetUserByUid(ctx, logger, uid)
	if err != nil {
		return err
	}
	return nil
}

// Create user
func (um *UserModel) CreateUser(ctx context.Context, logger *zerolog.Logger, user *schema.User) (err error) {
	if err = um.CheckUser(ctx, logger, user.UID); err == nil {
		// 用户已经存在
		return errors.New(consts.CacheAlreadyExist)
	}
	var sqlUser sqlmodel.User
	if err = copier.Copy(&sqlUser, user); err != nil {
		return err
	}
	if err = um.db.User.WithContext(ctx).Create(&sqlUser); err != nil {
		return err
	}
	go um.encodeUserInfoToCache(ctx, logger, user)
	return nil
}

func (um *UserModel) GetUserByUid(ctx context.Context, logger *zerolog.Logger, uid string) (user *schema.User, err error) {
	user = new(schema.User)
	// 缓存找
	user, err = um.getUserinfoFromCache(ctx, logger, uid)
	if err != nil {
		logger.Info().Msgf("uid: %s, not found in cache", uid)
	} else {
		return user, nil
	}

	user = new(schema.User)
	// 数据库找
	sqlUser, err := um.db.User.WithContext(ctx).Where(query.User.UID.Eq(uid)).Last()
	if err != nil {
		logger.Info().Msgf("uid: %s, not found in mysql", uid)
		return nil, err
	}
	copier.Copy(user, sqlUser)

	go um.encodeUserInfoToCache(ctx, logger, user)
	return user, nil
}

// Delete user
func (um *UserModel) DeleteUser(ctx context.Context, logger *zerolog.Logger, uid string) (err error) {
	// 延时双删除
	defer func() {
		go func() {
			time.Sleep(consts.DelayedDeletionInterval)
			um.cache.HDel(ctx, consts.RdsHashUserInfoKey, uid)
		}()
	}()
	// 缓存找
	_, err = um.getUserinfoFromCache(ctx, logger, uid)
	if err != nil {
		logger.Info().Msgf("uid: %s, not found in cache", uid)
	} else {
		// 删缓存
		um.cache.HDel(ctx, consts.RdsHashUserInfoKey, uid)
	}

	// 数据库找
	_, err = um.db.User.WithContext(ctx).Where(query.User.UID.Eq(uid)).Last()
	if err != nil {
		logger.Info().Msgf("uid: %s, not found in mysql", uid)
		return
	}
	// 数据库删除数据
	_, err = um.db.User.WithContext(ctx).Where(query.User.UID.Eq(uid)).Delete()
	if err != nil {
		logger.Info().Msgf("uid: %s, delete failed in mysql", uid)
		return
	}

	return nil
}

// EditUser 编辑用户信息
func (um *UserModel) UpdateUser(ctx context.Context, logger *zerolog.Logger, user *schema.User) (err error) {
	// 删除缓存
	um.cache.HDel(ctx, consts.RdsHashUserInfoKey, user.UID)

	// 延时双删
	defer func() {
		go func() {
			time.Sleep(consts.DelayedDeletionInterval)
			um.cache.HDel(ctx, consts.RdsHashUserInfoKey, user.UID)
		}()
	}()

	// 更新数据
	um.db.Transaction(func(tx *query.Query) error {
		defer func() {
			if r := recover(); r != nil {
				logger.Info().Any("recover", r).Send()
			}
		}()

		_, err := query.User.WithContext(ctx).
			Where(query.User.UID.Eq(user.UID)).
			Updates(sqlmodel.User{FullName: user.FullName, AvatarURL: user.AvatarURL, SessionID: user.SessionID, Biography: user.Biography, PrivateKey: user.PrivateKey, ClientID: user.ClientID, Contract: user.Contract, IsMvmUser: user.IsMvmUser})
		if err != nil {
			return err
		}

		return nil
	})

	return nil
}

// GetUserByName gets an user by the username.
func (um *UserModel) GetUserByFullName(ctx context.Context, logger *zerolog.Logger, full_name string) (user *schema.User, err error) {
	user = new(schema.User)
	// 数据库找
	sqlUser, err := um.db.User.WithContext(ctx).Where(query.User.FullName.Eq(full_name)).Last()
	if err != nil {
		logger.Info().Msgf("full_name: %s, not found in mysql", full_name)
		return
	}
	copier.Copy(sqlUser, user)
	return user, nil
}

func (um *UserModel) ListUser(ctx context.Context, logger *zerolog.Logger, offset, limit int) ([]*schema.User, int, error) {
	var users []*schema.User
	var count int64
	return users, int(count), nil
}

func (um *UserModel) encodeUserInfoToCache(ctx context.Context, logger *zerolog.Logger, userInfo *schema.User) {
	bytes, err := convert.Marshal(userInfo)
	if err != nil {
		logger.Error().Msgf("encode node to bytes fail, %+v", userInfo)
		return
	}

	err = um.cache.HSet(ctx, consts.RdsHashUserInfoKey, userInfo.UID, bytes)
	if err != nil {
		logger.Error().Msgf("encode user to redis fail, %+v", userInfo)
	}
}

func (um *UserModel) getUserinfoFromCache(ctx context.Context, logger *zerolog.Logger, uid string) (*schema.User, error) {
	bytes, err := um.cache.HGet(ctx, consts.RdsHashUserInfoKey, uid)
	// 找到了数据
	if err == nil && bytes != nil {
		var userInfo schema.User
		convert.Unmarshal(bytes, &userInfo)
		return &userInfo, nil
	}
	// 没有找到数据
	if err == cache.Nil {
		return nil, errors.New(consts.CacheNotFound)
	}
	// redis错误
	return nil, err
}
