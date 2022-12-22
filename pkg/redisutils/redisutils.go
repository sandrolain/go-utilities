package redisutils

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/gob"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	DefaultHost = "localhost"
	DefaultPort = 6379
)

func DefaultAddress() string {
	return fmt.Sprintf("%v:%v", DefaultHost, DefaultPort)
}

func createContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout*time.Second)
}

type Key []string

func (k *Key) String() string {
	return strings.Join(*k, ":")
}

type Client struct {
	client  *redis.Client
	timeout time.Duration
}

func NewClient(address string, password string, tls *tls.Config, timeout time.Duration) (*Client, error) {
	if address == "" {
		return nil, fmt.Errorf("empty Redis address")
	}
	if password == "" {
		return nil, fmt.Errorf("empty Redis password")
	}
	if timeout == 0 {
		return nil, fmt.Errorf("empty Redis timeout")
	}

	client := redis.NewClient(&redis.Options{
		Addr:      address,
		Password:  password,
		DB:        0, // use default DB
		TLSConfig: tls,
	})

	ctx, cancel := createContext(timeout)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	res := Client{
		client:  client,
		timeout: timeout,
	}
	return &res, nil
}

func (c *Client) Set(key Key, value interface{}, ttl time.Duration) error {
	ctx, cancel := createContext(c.timeout)
	defer cancel()
	var gobBuff bytes.Buffer
	enc := gob.NewEncoder(&gobBuff)
	err := enc.Encode(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key.String(), gobBuff.Bytes(), ttl).Err()
}

func (c *Client) SetNoTtl(key Key, value interface{}, ttl time.Duration) error {
	return c.Set(key, value, 0)
}

func (c *Client) Get(key Key, value interface{}) (bool, error) {
	ctx, cancel := createContext(c.timeout)
	defer cancel()
	gobValue, err := c.client.Get(ctx, key.String()).Bytes()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	var gobBuff bytes.Buffer
	gobBuff.Write(gobValue)
	enc := gob.NewDecoder(&gobBuff)
	return true, enc.Decode(value)
}

func (c *Client) GetAll(key Key, value interface{}) ([]interface{}, error) {
	ctx, cancel := createContext(c.timeout)
	defer cancel()
	key = append(key, "*")
	iter := c.client.Scan(ctx, 0, key.String(), 0).Iterator()
	res := make([]interface{}, 0)
	for iter.Next(ctx) {
		key := iter.Val()
		ok, err := c.Get(Key{key}, value)
		if err != nil {
			return res, err
		}
		if ok {
			value = append(res, value)
		}
	}
	if err := iter.Err(); err != nil {
		return res, err
	}
	return res, nil
}

func AllAsType[T interface{}](all []interface{}) ([]T, error) {
	res := make([]T, len(all))
	for i, v := range all {
		val, ok := v.(T)
		if !ok {
			return res, fmt.Errorf("item %v not convertible", i)
		}
		res[i] = val
	}
	return res, nil
}

func (c *Client) Delete(key Key) error {
	ctx, cancel := createContext(c.timeout)
	defer cancel()
	return c.client.Del(ctx, key.String()).Err()
}
