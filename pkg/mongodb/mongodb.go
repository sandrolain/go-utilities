package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DatabaseClient struct {
	client  *mongo.Client
	db      *mongo.Database
	timeout time.Duration
}

func (c *DatabaseClient) Close() {
	defer func() {
		ctx, cancel := createContext(c.timeout)
		defer cancel()
		if err := c.client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func (c *DatabaseClient) Ping() error {
	ctx, cancel := createContext(c.timeout)
	defer cancel()
	if err := c.client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("Database Ping success")
	return nil
}

func (c *DatabaseClient) AssertUniqueIndex(collection string, field string) (string, error) {
	return c.Coll(collection).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: field, Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
}

func (c *DatabaseClient) AssertIndex(collection string, field string) (string, error) {
	return c.Coll(collection).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{{Key: field, Value: 1}},
		},
	)
}

func (c *DatabaseClient) Coll(collection string) *mongo.Collection {
	return c.db.Collection(collection)
}

func (c *DatabaseClient) FindOne(collection string, filter interface{}, v interface{}) (bool, error) {
	ctx, cancel := createContext(c.timeout)
	defer cancel()
	err := c.db.Collection(collection).FindOne(ctx, filter).Decode(v)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *DatabaseClient) FindOneById(collection string, id primitive.ObjectID, v interface{}) (bool, error) {
	filter := bson.M{"_id": id}
	return c.FindOne(collection, filter, v)
}

func (c *DatabaseClient) FindOneByField(collection string, field string, value interface{}, v interface{}) (bool, error) {
	filter := bson.M{field: value}
	return c.FindOne(collection, filter, v)
}

func (c *DatabaseClient) FindMany(collection string, filter interface{}, sort interface{}, limit int64, v interface{}) error {
	ctx, cancel := createContext(c.timeout)
	defer cancel()
	opts := options.Find()
	if sort != nil {
		opts.SetSort(sort)
	}
	if limit > 0 {
		opts.SetLimit(limit)
	}
	cursor, err := c.db.Collection(collection).Find(ctx, filter, opts)
	if err != nil {
		return err
	}
	return cursor.All(ctx, v)
}

func (c *DatabaseClient) FindManyByField(collection string, field string, value interface{}, sort int, limit int64, v interface{}) error {
	filter := bson.M{field: value}
	sortExp := bson.M{field: sort}
	return c.FindMany(collection, filter, sortExp, limit, v)
}

func (c *DatabaseClient) InsertOne(collection string, v interface{}) (*mongo.InsertOneResult, error) {
	ctx, cancel := createContext(c.timeout)
	defer cancel()
	return c.db.Collection(collection).InsertOne(ctx, v)
}

func (c *DatabaseClient) InsertMany(collection string, v []interface{}) (*mongo.InsertManyResult, error) {
	ctx, cancel := createContext(c.timeout)
	defer cancel()
	return c.db.Collection(collection).InsertMany(ctx, v)
}

func (c *DatabaseClient) UpdateOne(collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	opts := options.Update()
	ctx, cancel := createContext(c.timeout)
	defer cancel()
	return c.db.Collection(collection).UpdateOne(ctx, filter, bson.M{"$set": update}, opts)
}

func (c *DatabaseClient) UpdateOneByField(collection string, field string, value interface{}, update interface{}) (*mongo.UpdateResult, error) {
	filter := bson.M{field: value}
	return c.UpdateOne(collection, filter, update)
}

func (c *DatabaseClient) UpdateOneById(collection string, id primitive.ObjectID, update interface{}) (*mongo.UpdateResult, error) {
	return c.UpdateOne(collection, bson.M{"_id": id}, update)
}

func (c *DatabaseClient) UpsertOne(collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(true)
	ctx, cancel := createContext(c.timeout)
	defer cancel()
	return c.db.Collection(collection).UpdateOne(ctx, filter, bson.M{"$set": update}, opts)
}

func (c *DatabaseClient) UpsertOneByField(collection string, field string, value interface{}, update interface{}) (*mongo.UpdateResult, error) {
	filter := bson.M{field: value}
	return c.UpsertOne(collection, filter, update)
}

func (c *DatabaseClient) UpsertOneById(collection string, id primitive.ObjectID, update interface{}) (*mongo.UpdateResult, error) {
	return c.UpsertOne(collection, bson.M{"_id": id}, update)
}

func (c *DatabaseClient) DeleteOne(collection string, filter interface{}) (*mongo.DeleteResult, error) {
	ctx, cancel := createContext(c.timeout)
	defer cancel()
	return c.db.Collection(collection).DeleteOne(ctx, filter)
}

func (c *DatabaseClient) DeleteOneByField(collection string, field string, value interface{}) (*mongo.DeleteResult, error) {
	return c.DeleteOne(collection, bson.M{field: value})
}

func (c *DatabaseClient) DeleteOneById(collection string, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	return c.DeleteOne(collection, bson.M{"_id": id})
}

func createContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout*time.Second)
}

func connect(uri string, timeout time.Duration) (*mongo.Client, error) {
	ctx, _ := createContext(timeout)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, err
}

func NewDatabaseClient(uri string, database string, timeout time.Duration) (*DatabaseClient, error) {
	client, err := connect(uri, timeout)
	if err != nil {
		return nil, err
	}
	return &DatabaseClient{
		client:  client,
		db:      client.Database(database),
		timeout: timeout,
	}, nil
}
