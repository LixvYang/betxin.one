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
		// FROM category
		// WHERE id = @id;
		GetCategoryById(ctx context.Context, id int64) (*Category, error)

		// SELECT
		// *
		// FROM category;
		ListCategories(ctx context.Context) ([]*Category, error)

		// SELECT
		// 	*
		// FROM category
		// WHERE name = @name;
		GetCategoryByName(ctx context.Context, name string) (*Category, error)

		// INSERT INTO category
		// (name)
		// VALUES
		// (@name);
		CreateCategory(ctx context.Context, name string) error

		// DELETE
		// FROM category
		// WHERE
		// name = @name;
		DeleteCategory(ctx context.Context, name string) error

		// Update category SET
		// name = @category.Name
		// WHERE id = @category.ID;
		UpdateCategory(ctx context.Context, category *Category) error
	}
)
