package core

import "context"

type (
	Category struct {
		ID   int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
		Name string `gorm:"column:name;not null" json:"name"`
	}

	CategoryStore interface {
		// SELECT
		// 	*
		// FROM @@table
		// WHERE id = @id;
		GetCategoryById(ctx context.Context, id int64) (*Category, error)

		// SELECT
		// *
		// FROM @@table;
		ListCategories(ctx context.Context) ([]*Category, error)

		// SELECT
		// 	*
		// FROM @@table
		// WHERE name = @name;
		GetCategoryByName(ctx context.Context, name string) (*Category, error)
	}
)
