package mongoutils

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client struct {
	client  *mongo.Client
	db      *mongo.Database
	timeout time.Duration
}

func (c *Client) Close() {
	defer func() {
		ctx, cancel := createContext(c.timeout)
		defer cancel()
		if err := c.client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func (c *Client) Ping() error {
	ctx, cancel := createContext(c.timeout)
	defer cancel()
	if err := c.client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("Database Ping success")
	return nil
}

func (c *Client) AssertUniqueIndex(collection string, field string) (string, error) {
	return c.Coll(collection).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: field, Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
}

func (c *Client) AssertIndex(collection string, field string) (string, error) {
	return c.Coll(collection).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{{Key: field, Value: 1}},
		},
	)
}

func (c *Client) AssertTtlIndex(collection string, field string, expireSeconds int32) (string, error) {
	return c.Coll(collection).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{{Key: field, Value: 1}},
			Options: &options.IndexOptions{
				ExpireAfterSeconds: &expireSeconds,
			},
		},
	)
}

func (c *Client) Coll(collection string) *mongo.Collection {
	return c.db.Collection(collection)
}

func (c *Client) FindOne(collection string, filter interface{}, v interface{}) (bool, error) {
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

func (c *Client) FindOneById(collection string, id interface{}, v interface{}) (bool, error) {
	filter := bson.M{"_id": id}
	return c.FindOne(collection, filter, v)
}

func (c *Client) FindOneByField(collection string, field string, value interface{}, v interface{}) (bool, error) {
	filter := bson.M{field: value}
	return c.FindOne(collection, filter, v)
}

func (c *Client) FindMany(collection string, filter interface{}, sort interface{}, limit int64, v interface{}) error {
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

func (c *Client) FindManyByField(collection string, field string, value interface{}, limit int64, v interface{}) error {
	filter := bson.M{field: value}
	return c.FindMany(collection, filter, nil, limit, v)
}

func (c *Client) InsertOne(collection string, v interface{}) (*mongo.InsertOneResult, error) {
	ctx, cancel := createContext(c.timeout)
	defer cancel()
	return c.db.Collection(collection).InsertOne(ctx, v)
}

func (c *Client) InsertMany(collection string, v []interface{}) (*mongo.InsertManyResult, error) {
	ctx, cancel := createContext(c.timeout)
	defer cancel()
	return c.db.Collection(collection).InsertMany(ctx, v)
}

func (c *Client) UpdateOne(collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	opts := options.Update()
	ctx, cancel := createContext(c.timeout)
	defer cancel()
	return c.db.Collection(collection).UpdateOne(ctx, filter, bson.M{"$set": update}, opts)
}

func (c *Client) UpdateOneByField(collection string, field string, value interface{}, update interface{}) (*mongo.UpdateResult, error) {
	filter := bson.M{field: value}
	return c.UpdateOne(collection, filter, update)
}

func (c *Client) UpdateOneById(collection string, id interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return c.UpdateOne(collection, bson.M{"_id": id}, update)
}

func (c *Client) UpsertOne(collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(true)
	ctx, cancel := createContext(c.timeout)
	defer cancel()
	return c.db.Collection(collection).UpdateOne(ctx, filter, bson.M{"$set": update}, opts)
}

func (c *Client) UpsertOneByField(collection string, field string, value interface{}, update interface{}) (*mongo.UpdateResult, error) {
	filter := bson.M{field: value}
	return c.UpsertOne(collection, filter, update)
}

func (c *Client) UpsertOneById(collection string, id interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return c.UpsertOne(collection, bson.M{"_id": id}, update)
}

func (c *Client) DeleteOne(collection string, filter interface{}) (*mongo.DeleteResult, error) {
	ctx, cancel := createContext(c.timeout)
	defer cancel()
	return c.db.Collection(collection).DeleteOne(ctx, filter)
}

func (c *Client) DeleteOneByField(collection string, field string, value interface{}) (*mongo.DeleteResult, error) {
	return c.DeleteOne(collection, bson.M{field: value})
}

func (c *Client) DeleteOneById(collection string, id interface{}) (*mongo.DeleteResult, error) {
	return c.DeleteOne(collection, bson.M{"_id": id})
}

func (c *Client) DeleteMany(collection string, filter interface{}) (*mongo.DeleteResult, error) {
	ctx, cancel := createContext(c.timeout)
	defer cancel()
	return c.db.Collection(collection).DeleteMany(ctx, filter)
}

func createContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout*time.Second)
}

func connect(uri string, timeout time.Duration) (*mongo.Client, error) {
	ctx, _ := createContext(timeout)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, err
}

func NewClient(uri string, database string, timeout time.Duration) (*Client, error) {
	if uri == "" {
		return nil, fmt.Errorf("empty MongoDB URI")
	}
	if database == "" {
		return nil, fmt.Errorf("empty MongoDB database")
	}
	if timeout == 0 {
		return nil, fmt.Errorf("empty MongoDB timeout")
	}

	client, err := connect(uri, timeout)
	if err != nil {
		return nil, err
	}
	return &Client{
		client:  client,
		db:      client.Database(database),
		timeout: timeout,
	}, nil
}
