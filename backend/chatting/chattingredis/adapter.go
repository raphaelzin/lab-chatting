package chattingredis

import (
	redisInstance "main/redis"

	"github.com/go-redis/redis"
)

const (
	// Stream with past messages for recovery on login
	constStream = "chatting-stream"

	// Stream with all messages
	constAuditStream = "chatting-audit-channel"

	// Pub sub message-only channel
	constChannel = "message-channel"
)

func AddMessage(data []byte) {
	saveMessageToStream(data, constStream)
}

func saveMessageToStream(data []byte, channel string) {
	strData := string(data)
	rdata := make(map[string]interface{})
	rdata["data"] = strData
	cmd := redisInstance.Client.XAdd(&redis.XAddArgs{Stream: channel, MaxLen: 10, Values: rdata})
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

func PublishMessage(data []byte, isText bool) {
	saveMessageToStream(data, constAuditStream)

	if isText {
		AddMessage(data)
		err := redisInstance.Client.Publish(constChannel, string(data)).Err()
		if err != nil {
			panic(err)
		}
	}
}

func GetSubscriptionToChannel() *redis.PubSub {
	return redisInstance.Client.Subscribe(constChannel)
}
