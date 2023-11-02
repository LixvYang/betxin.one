package model

import (
	"context"
	"fmt"

	"github.com/qiniu/qmgo"
)

type MongoService struct {
}

func Init() {
	ctx := context.Background()
	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: "mongodb://localhost:27017"})
	if err != nil {
		panic(err)
	}
	// cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://localhost:27017", Database: "class", Coll: "user"})

	db := client.Database("class")
	coll := db.Collection("user")
	fmt.Println(coll)
}

func Close() {

}
