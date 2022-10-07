package chattingredis

import (
	"context"
	"log"

	"github.com/go-redis/redis"
)

var (
	client *redis.Client
	ctx    = context.Background()
)

const (
	ConstStream = "chatting-stream-1"
)

func Init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func GetNFromStream(stream string, n int64) (redis.XStream, error) {
	streams := make([]string, 0, 1)
	streams = append(streams, stream)
	streams = append(streams, "0")
	args := &redis.XReadArgs{Streams: streams, Count: n, Block: -1}
	cmd := client.XRead(args)

	result, err := cmd.Result()

	if err != nil {
		return *new(redis.XStream), err
	}

	for _, value := range result {
		if value.Stream == stream {
			return value, nil
		}
	}

	return *new(redis.XStream), err
}

func AddToStream(data []byte) {
	strData := string(data)
	rdata := make(map[string]interface{})
	rdata["data"] = strData

	args := &redis.XAddArgs{Stream: ConstStream, Values: rdata}
	cmd := client.XAdd(args)
	log.Println(cmd.Result())
}
