package category

import (
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/query"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
)

type CategoryModel struct {
	db    *query.Query
	cache *cache.Cache
}

func NewUserModel(query *query.Query, cache *cache.Cache) CategoryModel {
	return CategoryModel{
		db:    query,
		cache: cache,
	}
}

func (cm *CategoryModel) CheckCategory(name string) (err error) {
	return nil
}

func (cm *CategoryModel) CreateCategory(name string) (err error) {
	return nil
}

func (cm *CategoryModel) GetCategoryById(id int64) (*schema.Category, error) {
	return nil, nil
}

func (cm *CategoryModel) ListCategories() ([]*schema.Category, error) {
	return nil, nil
}

func (cm *CategoryModel) UpdateCategory(id int64, name string) (err error) {
	return nil
}

func (cm *CategoryModel) DeleteCategory(id int64) {

}
