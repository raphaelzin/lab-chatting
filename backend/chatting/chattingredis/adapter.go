package chattingredis

import "github.com/go-redis/redis"

func AddMessage(data []byte) {
	AddToStream(data)
}

func GetLastN(n int64) ([][]byte, error) {
	result, err := GetNFromStream(ConstStream, n)
	if err != nil {
		return make([][]byte, 0), err
	}

	messages := make([][]byte, 0)
	streamMessages := result.Messages
	for _, streamMessage := range streamMessages {
		rawData := streamMessage.Values["data"].(string)
		messages = append(messages, []byte(rawData))
	}

	return messages, nil
}

func PublishMessage(data []byte) {
	err := client.Publish("message-channel", string(data)).Err()
	if err != nil {
		panic(err)
	}
}

func GetSubscriptionToChannel(channel string) *redis.PubSub {
	return client.Subscribe(channel)
}
