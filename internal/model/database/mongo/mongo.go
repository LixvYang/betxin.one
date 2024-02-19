package mongo

import (
	"context"
	"fmt"

	"github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoService struct {
	userColl     *qmgo.Collection
	categoryColl *qmgo.Collection
	bonuseColl *qmgo.Collection
	collectColl  *qmgo.Collection
	topicColl    *qmgo.Collection
}

func NewMongoService(logger *zerolog.Logger, conf *config.AppConfig) *MongoService {
	ctx := context.Background()
	mongoConf := conf.MongoConfig
	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: fmt.Sprintf("mongodb://%s:%d", mongoConf.Host, mongoConf.Port), Auth: &qmgo.Credential{
		Username: conf.MongoConfig.Username,
		Password: conf.MongoConfig.Password,
	}})
	if err != nil {
		panic(err)
	}

	client.Database(mongoConf.DB).Collection(consts.UserCollection).CreateIndexes(ctx, []opts.IndexModel{
		{
			Key:          []string{"uid"},
			IndexOptions: options.Index().SetUnique(true),
		},
		{
			Key:          []string{"identity_number"},
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
			Key: []string{"cid"},
		},
		{
			Key:          []string{"tid"},
			IndexOptions: options.Index().SetUnique(true),
		},
	})

	ms := &MongoService{
		userColl:     client.Database(mongoConf.DB).Collection(consts.UserCollection),
		categoryColl: client.Database(mongoConf.DB).Collection(consts.CategoryCollection),
		collectColl:  client.Database(mongoConf.DB).Collection(consts.CollectCollection),
		topicColl:    client.Database(mongoConf.DB).Collection(consts.TopicCollection),
	}
	ms.initCategory()
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
			Name: "Buisiness",
		},
		{
			ID:   2,
			Name: "Crypto",
		},
		{
			ID:   3,
			Name: "Sports",
		},
		{
			ID:   4,
			Name: "Games",
		},
		{
			ID:   5,
			Name: "News",
		},
		{
			ID:   6,
			Name: "Trending",
		},
		{
			ID:   7,
			Name: "Others",
		},
	}

	for _, category := range categorys {
		m.upsertCategory(context.Background(), &log.Logger, category.ID, category.Name)
	}
}
