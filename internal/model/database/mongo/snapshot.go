package mongo

import (
	"context"

	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *MongoService) InsertSnapshot(ctx context.Context, snapshot *schema.Snapshot) error {
	_, err := s.snapShotColl.InsertOne(ctx, snapshot)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrItemExist
		}
		return err
	}

	return nil
}

func (s *MongoService) GetSnapshotByRequestId(ctx context.Context, requestId string) (*schema.Snapshot, error) {
	var snapshot schema.Snapshot
	filter := bson.M{"request_id": requestId}
	err := s.snapShotColl.Find(ctx, filter).One(&snapshot)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoSuchItem
		}
		return nil, err
	}

	return &snapshot, nil
}

func (s *MongoService) GetLastestSnapshot(ctx context.Context) (*schema.Snapshot, error) {
	var snapshot schema.Snapshot
	err := s.snapShotColl.Find(ctx, bson.M{}).Sort("-created_at").Limit(1).One(&snapshot)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoSuchItem
		}
		return nil, err
	}

	return &snapshot, nil
}

func (s *MongoService) GetSnapshotCount(ctx context.Context) (int64, error) {
	count, err := s.snapShotColl.Find(ctx, bson.M{}).Count()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// func (s *MongoService) InsertSnapshotList(ctx context.Context, log *zerolog.Logger, snapshotList []*schema.Snapshot) error {
// 	_, err := s.snapShotColl.InsertMany(ctx, snapshotList)

// 	return nil
// }
