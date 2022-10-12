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
	cmd := redisInstance.Client.XAdd(&redis.XAddArgs{Stream: constStream, MaxLen: 10, Values: rdata})
	cmd.Result()
}

func GetLastN(n int64) ([][]byte, error) {
	cmd := redisInstance.Client.XRevRangeN(constStream, "+", "-", n)
	streamMessages, err := cmd.Result()

	if err != nil {
		return make([][]byte, 0), err
	}

	messages := make([][]byte, 0)
	for i := len(streamMessages) - 1; i >= 0; i-- {
		rawData := streamMessages[i].Values["data"].(string)
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
