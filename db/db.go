package db

import (
	"context"
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"sync"
)

var client *qmgo.Client
var once sync.Once
var ctx context.Context
var dataBase *qmgo.Database

type Collection struct {
	col *qmgo.Collection
}

// 指定了数据库
func Init() {
	once.Do(func() {
		var err error
		ctx = context.Background()
		client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: "mongodb://localhost:27017"})
		dataBase = client.Database("blakeyi")
		if err != nil {
			logrus.Fatal(err)
		}
	})
}

func NewCollection(col string) *Collection {
	return &Collection{
		dataBase.Collection(col),
	}
}


// list 不需要content
func (c *Collection) FindAll(filter interface{}, opts ...opts.FindOptions) qmgo.QueryI {
	return c.col.Find(ctx, filter, opts...).Select(bson.M{"content":0})
}

func (c *Collection) FindOne(filter interface{}, opts ...opts.FindOptions) qmgo.QueryI {
	return c.col.Find(ctx, filter, opts...)
}

func (c *Collection) InsertOne(doc interface{}, opts ...opts.InsertOneOptions) (result *qmgo.InsertOneResult, err error) {
	return c.col.InsertOne(ctx, doc, opts...)
}

func (c *Collection) UpdateOne(filter interface{}, update interface{}, opts ...opts.UpdateOptions) (err error) {
	return c.col.UpdateOne(ctx, filter, update, opts...)
}

func (c *Collection) DeleteOne(filter interface{}, opts ...opts.RemoveOptions) (err error) {
	return c.col.Remove(ctx, filter, opts...)
}


