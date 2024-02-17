package mongo

import (
	"context"
	"fmt"

	"github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoService struct {
	userColl     *qmgo.Collection
	categoryColl *qmgo.Collection
}

func NewMongoService(logger *zerolog.Logger, conf *config.AppConfig) *MongoService {
	ctx := context.Background()
	mongoConf := conf.MongoConfig
	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: fmt.Sprintf("mongodb://%s:%d", mongoConf.Host, mongoConf.Port)})
	if err != nil {
		panic(err)
	}

	client.Database(mongoConf.DB).Collection(consts.UserCollection).CreateOneIndex(
		context.Background(),
		opts.IndexModel{
			Key:          []string{"uid"},
			IndexOptions: options.Index().SetUnique(true),
		},
	)

	client.Database(mongoConf.DB).Collection(consts.UserCollection).CreateIndexes(ctx, []opts.IndexModel{
		{
			Key:          []string{"cid", "name"},
			IndexOptions: options.Index().SetUnique(true),
		},
		{
			Key:          []string{"name"},
			IndexOptions: options.Index().SetUnique(true),
		},
	})

	return &MongoService{
		userColl:     client.Database(mongoConf.DB).Collection(consts.UserCollection),
		categoryColl: client.Database(mongoConf.DB).Collection(consts.CategoryCollection),
	}
}

func Close() {

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
