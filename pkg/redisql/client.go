package redisql

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/juricake/redisql/pkg/redisql/pkg/keys"
	"github.com/juricake/redisql/pkg/redisql/pkg/schema_util"
	"strings"
)

type Client struct {
	redis *redis.Client
}

type Options struct {
	Host string
	Port int
}

func NewClient(opt Options) (*Client, error) {
	redisClient := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", opt.Host, opt.Port)})
	_, err := redisClient.Ping().Result()
	if err != nil {
		return nil, err
	}

	c := &Client{
		redis: redisClient,
	}
	return c, nil
}

func (c *Client) CreateTable(name string, schema interface{}) error {
	exists, err := c.redis.HExists(keys.Schemas, name).Result()
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("redisql: table '%s' already exists", name)
	}

	encodedSchema, err := schema_util.Encode(schema)
	if err != nil {
		return fmt.Errorf("redisql: error encoding the schema: %s", err)
	}

	_, err = c.redis.HSet(keys.Schemas, name, encodedSchema).Result()
	if err != nil {
		return fmt.Errorf("redisql: error creating table '%s", name)
	}
	return nil
}

func (c *Client) PingTable(name string) error {
	exists, err := c.redis.HExists(keys.Schemas, name).Result()
	if err != nil {
		return fmt.Errorf("redisql: ping table command err: %s", err.Error())
	}

	if !exists {
		return fmt.Errorf("redisql: table '%s' not found", name)
	}

	return nil
}

func (c *Client) InsertInto(table string, data interface{}) error {
	schema, err := c.redis.HGet(keys.Schemas, table).Result()
	if err != nil {
		return fmt.Errorf("redisql: insert failed, table '%s' not found", table)
	}

	if err := schema_util.Validate(data, schema); err != nil {
		return fmt.Errorf("redisql: insert failed, data doesn't match the schema")
	}

	jsonValue, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("redisql: insert failed, unable to marshal the data: %s", err.Error())
	}

	_, err = c.redis.SAdd(table, jsonValue).Result()
	if err != nil {
		return fmt.Errorf("redisql: insert failed: %s", err.Error())
	}

	return nil
}

func (c *Client) SelectFrom(table string, results interface{}) error {
	resultsJson, err := c.redis.SMembers(table).Result()
	if err != nil {
		return fmt.Errorf("redisql: select failed: %s", err.Error())
	}

	resultsGroup := fmt.Sprintf("[%s]", strings.Join(resultsJson, ","))
	if err := json.Unmarshal([]byte(resultsGroup), results); err != nil {
		return fmt.Errorf("redisql: select failed, unable to unmarshal the data: %s", err.Error())
	}

	return nil
}
