package core

import (
	"context"
	"time"
)

type (
	Topic struct {
		ID            int64      `gorm:"column:id;primaryKey;autoIncrement:true;comment:ID" json:"id"` // ID
		Tid           string     `gorm:"column:tid;not null" json:"tid"`
		Cid           int64      `gorm:"column:cid;not null;comment:ID" json:"cid"` // ID
		Title         string     `gorm:"column:title;not null" json:"title"`
		Intro         string     `gorm:"column:intro;not null" json:"intro"`
		Content       string     `gorm:"column:content;not null" json:"content"`
		YesRatio      string     `gorm:"column:yes_ratio;not null;default:50.00" json:"yes_ratio"`
		NoRatio       string     `gorm:"column:no_ratio;not null;default:50.00" json:"no_ratio"`
		YesCount      string     `gorm:"column:yes_count;not null;default:0.00000000" json:"yes_count"`
		NoCount       string     `gorm:"column:no_count;not null;default:0.00000000" json:"no_count"`
		TotalCount    string     `gorm:"column:total_count;not null;default:0.00000000" json:"total_count"`
		CollectCount  int64      `gorm:"column:collect_count;not null" json:"collect_count"`
		ReadCount     int64      `gorm:"column:read_count;not null" json:"read_count"`
		ImgURL        string     `gorm:"column:img_url;not null;comment:URL" json:"img_url"` // URL
		IsStop        bool       `gorm:"column:is_stop" json:"is_stop"`
		RefundEndTime time.Time  `json:"refund_end_time"`
		EndTime       time.Time  `json:"end_time"`
		CreatedAt     time.Time  `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
		UpdatedAt     time.Time  `gorm:"column:updated_at;not null;autoUpdateTime" json:"updated_at"`
		DeletedAt     *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	}

	TopicStore interface {
		// SELECT
		// 	*
		// FROM topic
		// WHERE uid in (@tids) AND deleted_at IS NULL;
		GetTopicsByTids(ctx context.Context, tids []string) ([]*Topic, error)

		// SELECT
		// *
		// FROM topic
		// WHERE cid = @cid AND deleted_at IS NULL;
		GetTopicsByCid(ctx context.Context, cid int64) ([]*Topic, error)

		// SELECT
		// *
		// FROM topic
		// WHERE tid = @tid AND deleted_at IS NULL;
		GetTopicByTid(ctx context.Context, tid string) (*Topic, error)

		// DELETE FROM topic
		// WHERE
		// 	`tid` = @tid;
		DeleteTopic(ctx context.Context, tid string) error

		// UPDATE topic
		// {{set}}
		//   `is_stop` = 1
		// {{end}}
		// WHERE
		// 	`tid` = @tid;
		StopTopic(ctx context.Context, tid string) error

		// INSERT INTO topic
		// 	(
		// 		`tid`,`cid`,`title`,`intro`,
		// 		`content`,`yes_ratio`,`no_ratio`,`yes_count`,
		// 		`no_count`,`total_count`,`collect_count`,`read_count`,
		// 		`img_url`,`is_stop`,`refund_end_time`,`end_time`,
		// 		`created_at`,`updated_at`,`deleted_at`
		// 	)
		// VALUES
		// 	(
		// 		@topic.Tid,@topic.Cid,@topic.Title,@topic.Intro,
		// 		@topic.Content,"50.00","50.00","0",
		// 		"0","0","0","0",
		// 		@topic.ImgURL,0,@topic.RefundEndTime,@topic.EndTime,
		//		NOW(),NOW(),NULL
		// 	);
		CreateTopic(ctx context.Context, topic *Topic) error

		// SELECT
		// *
		// FROM topic
		// WHERE `cid` = @cid AND
		// `created_at` <= @created_at AND
		// `deleted_at` IS NULL AND
		// `is_stop` = 0
		// ORDER BY `created_at` DESC
		// LIMIT @page_size;
		ListTopicsByCid(ctx context.Context, cid int64, created_at time.Time, page_size int64) ([]*Topic, error)

		// UPDATE topic
		// {{set}}
		//   `tid` = @topic.Tid,
		//   `cid` = @topic.Cid,
		//   `title` = @topic.Title,
		//   `intro` = @topic.Intro,
		//   `content` = @topic.Content,
		//   `img_url` = @topic.ImgURL,
		//   `is_stop` = @topic.IsStop,
		//   `refund_end_time` = @topic.RefundEndTime,
		//   `end_time` = @topic.EndTime,
		//   `updated_at` = NOW()
		// {{end}}
		// WHERE
		// 	`tid` = @topic.Tid;
		UpdateTopicInfo(ctx context.Context, topic *Topic) error
	}
)
