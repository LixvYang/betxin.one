package category

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/query"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/sqlmodel"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/rs/zerolog"
)

type CategoryModel struct {
	db    *query.Query
	cache *cache.Cache
}

func NewCategoryModel(query *query.Query, cache *cache.Cache) CategoryModel {
	return CategoryModel{
		db:    query,
		cache: cache,
	}
}

func (cm *CategoryModel) CheckCategory(ctx context.Context, logger *zerolog.Logger, name string) (err error) {
	sqlCategory, err := cm.db.Category.WithContext(ctx).Where(query.Category.Name.Eq(name)).Last()
	if err != nil || sqlCategory == nil {
		return err
	}
	return nil
}

func (cm *CategoryModel) CreateCategory(ctx context.Context, logger *zerolog.Logger, name string) (err error) {
	return cm.db.Category.WithContext(ctx).Create(&sqlmodel.Category{Name: name})
}

func (cm *CategoryModel) GetCategoryByName(ctx context.Context, logger *zerolog.Logger, name string) (schemaCategory *schema.Category, err error) {
	schemaCategory = new(schema.Category)
	sqlCategory, err := cm.db.Category.WithContext(ctx).Where(query.Category.Name.Eq(name)).Last()
	if err != nil {
		return nil, err
	}

	copier.Copy(schemaCategory, &sqlCategory)
	return schemaCategory, nil
}

func (cm *CategoryModel) GetCategoryById(ctx context.Context, logger *zerolog.Logger, id int64) (schemaCategory *schema.Category, err error) {
	sqlCategory, err := cm.db.Category.WithContext(ctx).Where(query.Category.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}

	schemaCategory = new(schema.Category)
	copier.Copy(schemaCategory, sqlCategory)
	return schemaCategory, nil
}

func (cm *CategoryModel) ListCategories() (schemaCategory []*schema.Category, err error) {
	ctx := context.Background()
	sqlCategory, err := cm.db.Category.WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}

	schemaCategory = make([]*schema.Category, len(sqlCategory))
	copier.Copy(&schemaCategory, &sqlCategory)
	return schemaCategory, nil
}

func (cm *CategoryModel) UpdateCategory(ctx context.Context, logger *zerolog.Logger, id int64, name string) (err error) {
	info, err := cm.db.Category.WithContext(ctx).Where(query.Category.ID.Eq(id)).Update(query.Category.Name, name)
	if err != nil || info.RowsAffected == 0 {
		return err
	}
	return nil
}

func (cm *CategoryModel) DeleteCategory(ctx context.Context, logger *zerolog.Logger, id int64) error {
	info, err := cm.db.Category.WithContext(ctx).Where(query.Category.ID.Eq(id)).Delete()
	if err != nil || info.RowsAffected == 0 {
		return err
	}
	return nil
}
