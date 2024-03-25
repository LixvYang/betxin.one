package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/lixvyang/betxin.one/internal/utils/convert"

	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	ErrItemExist  error = errors.New("item already exist")
	ErrNoSuchItem error = errors.New("item no exist")
)

type MongoService struct {
	client *qmgo.Client

	userColl                 *qmgo.Collection
	categoryColl             *qmgo.Collection
	bonuseColl               *qmgo.Collection
	collectColl              *qmgo.Collection
	topicColl                *qmgo.Collection
	refundColl               *qmgo.Collection
	topicPurchaseColl        *qmgo.Collection
	topicPurchaseHistoryColl *qmgo.Collection
	snapShotColl             *qmgo.Collection
}

func NewMongoService(mongoConf *config.MongoConfig) *MongoService {
	ctx := context.Background()

	mgoConf := &qmgo.Config{
		Uri:            fmt.Sprintf("mongodb://%s:%d", mongoConf.Host, mongoConf.Port),
		ReadPreference: &qmgo.ReadPref{Mode: readpref.PrimaryMode},
	}

	if mongoConf.Password != "" {
		mgoConf.Auth = &qmgo.Credential{
			Username: mongoConf.Username,
			Password: mongoConf.Password,
		}
	}

	client, err := qmgo.NewClient(ctx, mgoConf)
	if err != nil {
		panic(err)
	}

	client.Database(mongoConf.DB).Collection(consts.UserCollection).CreateIndexes(ctx, []opts.IndexModel{
		{
			Key:          []string{"uid"},
			IndexOptions: options.Index().SetUnique(true),
		},
	})

	client.Database(mongoConf.DB).Collection(consts.CategoryCollection).CreateIndexes(ctx, []opts.IndexModel{
		{
			Key:          []string{"id", "name"},
			IndexOptions: options.Index().SetUnique(true),
		},
		{
			Key:          []string{"name"},
			IndexOptions: options.Index().SetUnique(true),
		},
	})

	client.Database(mongoConf.DB).Collection(consts.TopicCollection).CreateIndexes(ctx, []opts.IndexModel{
		{
			Key:          []string{"tid"},
			IndexOptions: options.Index().SetUnique(true),
		},
	})

	// ttl 30
	client.Database(mongoConf.DB).Collection(consts.SnapshotCollection).CreateIndexes(ctx, []opts.IndexModel{
		{
			Key:          []string{"request_id"},
			IndexOptions: options.Index().SetUnique(true),
		},
		{
			Key:          []string{"created_at"},
			IndexOptions: options.Index().SetExpireAfterSeconds(30 * 24 * 60 * 60), // 30 days
		},
	})

	// ttl 30
	client.Database(mongoConf.DB).Collection(consts.RefundCollection).CreateIndexes(ctx, []opts.IndexModel{
		{
			Key:          []string{"request_id"},
			IndexOptions: options.Index().SetUnique(true),
		},
		{
			Key:          []string{"created_at"},
			IndexOptions: options.Index().SetExpireAfterSeconds(7 * 24 * 60 * 60), // 30 days
		},
	})

	ms := &MongoService{
		client:                   client,
		userColl:                 client.Database(mongoConf.DB).Collection(consts.UserCollection),
		categoryColl:             client.Database(mongoConf.DB).Collection(consts.CategoryCollection),
		collectColl:              client.Database(mongoConf.DB).Collection(consts.CollectCollection),
		topicColl:                client.Database(mongoConf.DB).Collection(consts.TopicCollection),
		refundColl:               client.Database(mongoConf.DB).Collection(consts.RefundCollection),
		topicPurchaseColl:        client.Database(mongoConf.DB).Collection(consts.TopicPurchaseCollection),
		topicPurchaseHistoryColl: client.Database(mongoConf.DB).Collection(consts.TopicPurchaseHistoryCollection),
		bonuseColl:               client.Database(mongoConf.DB).Collection(consts.BonuseCollection),
		snapShotColl:             client.Database(mongoConf.DB).Collection(consts.MixinUtxoCollection),
	}

	ms.initCategory()
	// 定时同步utxos 并且定时聚合utxos
	return ms
}

func isMongoDupeKeyError(err error) bool {
	e, ok := err.(mongo.WriteException)
	if !ok {
		return false
	}
	for _, writeError := range e.WriteErrors {
		if writeError.Code == 11000 {
			return true
		}
	}
	return false
}

func (m *MongoService) initCategory() {
	var categorys []schema.Category = []schema.Category{
		{
			ID:   1,
			Name: "buisiness",
		},
		{
			ID:   2,
			Name: "crypto",
		},
		{
			ID:   3,
			Name: "sports",
		},
		{
			ID:   4,
			Name: "games",
		},
		{
			ID:   5,
			Name: "news",
		},
		{
			ID:   6,
			Name: "trending",
		},
		{
			ID:   7,
			Name: "others",
		},
	}

	for _, category := range categorys {
		m.upsertCategory(context.Background(), category.ID, category.Name)
	}
}

func (s *MongoService) HandleTopicStopAction(ctx context.Context, stopAction *schema.TopicStopAction) (*schema.StopTopicActionResp, error) {
	if _, err := convert.VaildUUID(stopAction.Tid); err != nil {
		return nil, err
	}

	now := time.Now()
	topic, err := s.GetTopicByTid(ctx, stopAction.Tid)
	if err != nil {
		return nil, err
	}

	// 如果话题已停止 || 话题还没过期
	if topic.IsStop || topic.EndTime.After(now) {
		return nil, errors.New("topic has stopped or not expired")
	}

	session, err := s.client.Session()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)
	// 获取话题信息和话题
	callback := func(sessCtx context.Context) (interface{}, error) {
		topic, err := s.GetTopicByTid(sessCtx, stopAction.Tid)
		if err != nil {
			return nil, qmgo.ErrTransactionRetry
		}

		topicPurchases, err := s.QueryTopicPurchase(ctx, "", stopAction.Tid)
		if err != nil {
			return nil, qmgo.ErrTransactionRetry
		}

		// 停止话题
		err = s.StopTopic(ctx, stopAction.Tid)
		if err != nil {
			return nil, qmgo.ErrTransactionRetry
		}

		resp := &schema.StopTopicActionResp{
			Topic:          *topic,
			TopicPurchases: topicPurchases,
		}
		return resp, nil
	}

	respAny, err := session.StartTransaction(ctx, callback)
	if err != nil {
		return nil, err
	}

	resp, ok := respAny.(*schema.StopTopicActionResp)
	if !ok {
		return nil, err
	}

	return resp, nil
}

func (s *MongoService) HandleMixinTopicRefundAction(ctx context.Context, refundAction *schema.TopicRefundAction) error {
	if refundAction.Amount.IsNegative() {
		return errors.New("invalid amount")
	}

	_, err := convert.VaildUUID(refundAction.Uid)
	if err != nil {
		return err
	}

	_, err = convert.VaildUUID(refundAction.Tid)
	if err != nil {
		return err
	}

	_, err = convert.VaildUUID(refundAction.RequestID)
	if err != nil {
		return err
	}

	session, err := s.client.Session()
	if err != nil {
		return err
	}

	defer session.EndSession(ctx)

	callback := func(sessCtx context.Context) (interface{}, error) {
		topic, err := s.GetTopicByTid(sessCtx, refundAction.Tid)
		if err != nil {
			return nil, err
		}

		var topicPurchase *schema.TopicPurchase
		topicPurchase, err = s.GetTopicPurchase(sessCtx, refundAction.Uid, refundAction.Tid)
		if err != nil {
			return nil, err
		}

		if topicPurchase == nil {
			return nil, errors.New("topic purchase not exist")
		}

		if refundAction.Action {
			yesCount := convert.NewDecimalFromString(topicPurchase.YesAmount)
			yesCount = yesCount.Sub(refundAction.Amount)
			if yesCount.IsNegative() {
				return nil, errors.New("refund topic purchase amount is greater than yes amount")
			}
			topicPurchase.YesAmount = yesCount.String()
		} else {
			noCount := convert.NewDecimalFromString(topicPurchase.NoAmount)
			noCount = noCount.Sub(refundAction.Amount)
			if noCount.IsNegative() {
				return nil, errors.New("refund topic purchase amount is greater than no amount")
			}
			topicPurchase.NoAmount = noCount.String()
		}

		if refundAction.Action {
			yesCount := convert.NewDecimalFromString(topic.YesAmount)
			yesCount = yesCount.Sub(refundAction.Amount)
			if yesCount.IsNegative() {
				return nil, errors.New("refund topic amount is greater than yes amount")
			}
			topic.YesAmount = yesCount.String()
		} else {
			noCount := convert.NewDecimalFromString(topic.NoAmount)
			noCount = noCount.Sub(refundAction.Amount)
			if noCount.IsNegative() {
				return nil, errors.New("refund topic amount is greater than no amount")
			}
			topic.NoAmount = noCount.String()
		}

		err = s.UpdateTopicPurchase(sessCtx, topicPurchase)
		if err != nil {
			return nil, qmgo.ErrTransactionRetry
		}

		err = s.UpdateTopic(sessCtx, topic.Tid, topic)
		if err != nil {
			return nil, qmgo.ErrTransactionRetry
		}

		// MIXIN 发送转账 获取reqid和memo
		err = s.CreateRefund(ctx, &schema.Refund{
			RequestID: refundAction.RequestID,
			Uid:       refundAction.Uid,
			Tid:       refundAction.Tid,
			Amount:    refundAction.Amount.String(),
			Action:    refundAction.Action,
			Memo:      refundAction.Memo,
			CreatedAt: time.Now(),
		})
		if err != nil {
			return nil, qmgo.ErrTransactionRetry
		}

		return nil, nil
	}

	_, err = session.StartTransaction(ctx, callback)
	if err != nil {
		return err
	}
	return nil
}

// 处理入账记录
func (s *MongoService) HandleMixinTopicDepositAction(ctx context.Context, buyAction *schema.TopicBuyAction) error {
	now := time.Now()
	if buyAction.Amount.IsNegative() {
		return errors.New("invalid amount")
	}

	_, err := convert.VaildUUID(buyAction.Uid)
	if err != nil {
		return err
	}

	_, err = convert.VaildUUID(buyAction.RequestID)
	if err != nil {
		return err
	}

	_, err = convert.VaildUUID(buyAction.Tid)
	if err != nil {
		return err
	}

	session, err := s.client.Session()
	if err != nil {
		return err
	}

	defer session.EndSession(ctx)

	callback := func(sessCtx context.Context) (interface{}, error) {
		topic, err := s.GetTopicByTid(sessCtx, buyAction.Tid)
		if err != nil {
			return nil, err
		}

		var topicPurchase *schema.TopicPurchase
		topicPurchase, err = s.GetTopicPurchase(sessCtx, buyAction.Uid, buyAction.Tid)
		if err != nil {
			s.CreateTopicPurchase(sessCtx, buyAction.Uid, buyAction.Tid)
			topicPurchase, err = s.GetTopicPurchase(sessCtx, buyAction.Uid, buyAction.Tid)
			if err != nil {
				return nil, err
			}
		}

		if topicPurchase == nil {
			return nil, errors.New("topic purchase not exist")
		}

		if buyAction.Action {
			yesCount := convert.NewDecimalFromString(topicPurchase.YesAmount)
			yesCount = yesCount.Add(buyAction.Amount)
			topicPurchase.YesAmount = yesCount.String()
		} else {
			noCount := convert.NewDecimalFromString(topicPurchase.NoAmount)
			noCount = noCount.Add(buyAction.Amount)
			topicPurchase.NoAmount = noCount.String()
		}

		if buyAction.Action {
			yesCount := convert.NewDecimalFromString(topic.YesAmount)
			yesCount = yesCount.Add(buyAction.Amount)
			topic.YesAmount = yesCount.String()
		} else {
			noCount := convert.NewDecimalFromString(topic.NoAmount)
			noCount = noCount.Add(buyAction.Amount)
			topic.NoAmount = noCount.String()
		}

		err = s.UpsertTopicPurchase(sessCtx, topicPurchase)
		if err != nil {
			return nil, qmgo.ErrTransactionRetry
		}

		err = s.UpdateTopic(sessCtx, topic.Tid, topic)
		if err != nil {
			return nil, qmgo.ErrTransactionRetry
		}

		// 更新话题购买行为
		err = s.UpsertTopicPurchaseHistory(ctx, &schema.TopicPurchaseHistory{
			Uid:        buyAction.Uid,
			Tid:        buyAction.Tid,
			Action:     buyAction.Action,
			Amount:     buyAction.Amount.String(),
			Memo:       buyAction.Memo,
			Finished:   true,
			FinishedAt: now,
		})
		if err != nil {
			return nil, qmgo.ErrTransactionRetry
		}

		return nil, nil
	}

	_, err = session.StartTransaction(ctx, callback)
	if err != nil {
		return err
	}

	return nil
}
