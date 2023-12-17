package core

import (
	"context"
	"time"
)

type (
	Collect struct {
		ID        int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
		UID       string    `json:"uid"`
		Tid       string    `json:"tid"`
		Status    bool      `json:"status"` // 状态
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	CollectStore interface {
		// SELECT
		// 	*
		// FROM collect
		// WHERE `id` = @id;
		GetCollectById(ctx context.Context, id int64) (*Collect, error)

		// SELECT
		// *
		// FROM collect
		// WHERE `tid` = @tid AND
		// `status` = 1;
		GetCollectByTid(ctx context.Context, tid string) ([]*Collect, error)

		// SELECT
		// *
		// FROM collect
		// WHERE `uid` = @uid AND
		// `status` = 1;
		GetCollectByUid(ctx context.Context, uid string) ([]*Collect, error)

		// SELECT
		// *
		// FROM collect
		// WHERE `uid` = @uid AND
		// `tid` = @tid AND
		// `status` = 1;
		GetCollectByUidTid(ctx context.Context, uid, tid string) (*Collect, error)

		// UPDATE  collect SET
		// `status` = 0
		// WHERE
		// `tid` = @tid AND `uid` = @uid;
		DeleteCollect(ctx context.Context, tid, uid string) error

		// INSERT INTO collect 
		// (`tid`, `uid`, `status`, `created_at`, `updated_at`)
		// VALUES
		// (@tid, @uid, 1, NOW(), NOW());
		CreateCollect(ctx context.Context, tid, uid string) error
	}
)
