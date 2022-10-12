package accessors

import (
	"encoding/json"
	"main/redis"
	"time"
)

type BaseAccessor[T any] struct {
}

func (a BaseAccessor[T]) cacheSave(id string, model T) error {
	data, err := json.Marshal(model)
	if err != nil {
		return err
	}
	err = redis.Client.Set(id, string(data), time.Hour*2).Err()
	return err
}

func (a BaseAccessor[T]) cacheGet(id string) (T, error) {
	result, err := redis.Client.Get(id).Result()

	if err != nil {
		return *new(T), err
	}

	var model T
	json.Unmarshal([]byte(result), &model)

	return model, nil
}
