package category

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lixvyang/betxin.one/configs"
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/query"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Test(t *testing.T) {
	var err error
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"123456",
		"127.0.0.1",
		3306,
		"betxin",
	)
	db, err = gorm.Open(mysql.Open(dns))
	assert.Nil(t, err)
	db = db.Debug()

	query.SetDefault(db)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	cache := cache.New(&configs.RedisConfig{
		Host:         "127.0.0.1",
		Password:     "123456",
		Port:         6379,
		DB:           0,
		PoolSize:     128,
		MinIdleConns: 100,
	})

	model := NewCategoryModel(query.Q, cache)
	TestDemo := "测试种类"
	err = model.CreateCategory(ctx, &log.Logger, TestDemo)
	assert.Nil(t, err)
	model.CheckCategory(ctx, &log.Logger, TestDemo)

	cate, err := model.GetCategoryByName(ctx, &log.Logger, TestDemo)
	assert.Nil(t, err)

	cateC, err := model.GetCategoryById(ctx, &log.Logger, cate.ID)
	assert.Nil(t, err)

	cates, err := model.ListCategories()
	assert.Nil(t, err)
	fmt.Println(cates)

	err = model.UpdateCategory(ctx, &log.Logger, cateC.ID, "测试更改名称")
	assert.Nil(t, err)
	cateCC, err := model.GetCategoryById(ctx, &log.Logger, cate.ID)
	assert.Nil(t, err)
	assert.Equal(t, cateCC.Name, "测试更改名称")

	err = model.DeleteCategory(ctx, &log.Logger, cateCC.ID)
	assert.Nil(t, err)

	cate, err = model.GetCategoryById(ctx, &log.Logger, cateCC.ID)
	assert.Error(t, err)
}
