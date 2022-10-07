package chattingredis

import (
	redisInstance "main/redis"

	"github.com/go-redis/redis"
)

const (
	constStream = "chatting-stream-1"
)

func AddMessage(data []byte) {
	strData := string(data)
	rdata := make(map[string]interface{})
	rdata["data"] = strData
	cmd := redisInstance.Client.XAdd(&redis.XAddArgs{Stream: constStream, Values: rdata})
	cmd.Result()
}

func GetLastN(n int64) ([][]byte, error) {
	streams := make([]string, 0, 1)
	streams = append(streams, constStream)
	streams = append(streams, "0")
	args := &redis.XReadArgs{Streams: streams, Count: n, Block: -1}
	cmd := redisInstance.Client.XRead(args)

	result, err := cmd.Result()

	if err != nil {
		return make([][]byte, 0), err
	}

	stream := *new(redis.XStream)
	for _, value := range result {
		if value.Stream == constStream {
			stream = value
			break
		}
	}

	messages := make([][]byte, 0)
	streamMessages := stream.Messages
	for _, streamMessage := range streamMessages {
		rawData := streamMessage.Values["data"].(string)
		messages = append(messages, []byte(rawData))
	}

	return messages, nil
}

func PublishMessage(data []byte) {
	err := redisInstance.Client.Publish("message-channel", string(data)).Err()
	if err != nil {
		panic(err)
	}
}

func GetSubscriptionToChannel(channel string) *redis.PubSub {
	return redisInstance.Client.Subscribe(channel)
}
