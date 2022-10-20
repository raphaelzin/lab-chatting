package user

import (
	"encoding/json"
	"main/models"
	redisInstance "main/redis"
)

func saveUser(username string) (models.User, error) {
	user := models.NewUser(username)

	data, err := json.Marshal(user)
	if err != nil {
		return *user, err
	}

	_, err = redisInstance.Client.Set(user.Id, string(data), 0).Result()

	return *user, err
}

func GetUserWithToken(token string) (models.User, error) {
	data, err := redisInstance.Client.Get(token).Result()
	var user models.User
	if err != nil {
		return user, err
	}

	err = json.Unmarshal([]byte(data), &user)
	return user, err
}
