package topic

import (
	"context"
	"time"

	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/query"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/sqlmodel"
	"github.com/lixvyang/betxin.one/pkg/safe"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/rs/zerolog"
)

func (tm *TopicModel) AfterFind(ctx context.Context, logger *zerolog.Logger, q *query.Query, sqlTopics ...*sqlmodel.Topic) error {
	safe.GoRun(func() {
		for _, sqlTopic := range sqlTopics {
			tm.db.WithContext(ctx).Topic.Where(q.Topic.Tid.Eq(sqlTopic.Tid)).Update(q.Topic.ReadCount, sqlTopic.ReadCount+1)
		}
	})

	return nil
}

func (tm *TopicModel) BeforeUpdate(ctx context.Context, logger *zerolog.Logger, q *query.Query, sqlTopics ...*sqlmodel.Topic) error {
	var errmsg error
	safe.Run(func() {
		for _, t := range sqlTopics {
			if t.IsDeleted {
				errmsg = errors.New("topic already deleted")
				break
			}
			t.UpdatedAt = time.Now().UnixMilli()

			if t.IsStop || time.Now().After(time.UnixMilli(t.EndTime)) {
				errmsg = errors.New("topic already stop")
				break
			}
			decimal.DivisionPrecision = 2
			yesCnt, _ := decimal.NewFromString(t.YesCount)

			totalCnt, err := decimal.NewFromString(t.TotalCount)
			if err != nil {
				errmsg = err
				break
			}
			if totalCnt.Equal(decimal.NewFromFloat(0)) {
				break
			}
			yesRatio := yesCnt.Div(totalCnt)
			t.YesRatio = yesRatio.String()
			t.NoRatio = decimal.NewFromInt(100).Sub(yesRatio).String()
		}
	})

	return errmsg
}
