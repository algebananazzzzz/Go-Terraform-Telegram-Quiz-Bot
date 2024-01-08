package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type User struct {
	Id            int64          `json:"id"`
	UserConvState int            `json:"conv_state"`
	UserData      UserDataStruct `json:"user_data"`
}

type UserDataStruct struct {
	QuizId       int  `json:"quiz_id"`
	Score        int  `json:"score"`
	PrevQnAnswer int8 `json:"prev_qn_answer"`
	NextQnId     int  `json:"next_qn_id"`
}

var ctx = context.Background()

var redisClient = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_ADDR"),
	Password: "", // no password set
	DB:       0,  // use default DB
})

func getUserData(userID int64) (User, error) {
	userData, err := redisClient.HGet(ctx, os.Getenv("REDIS_KEY"), fmt.Sprintf("%s%d", "user:", userID)).Result()
	if err == redis.Nil {
		return User{Id: userID, UserConvState: DefaultState}, nil
	}
	if err != nil {
		return User{}, err
	}
	var data User
	if err := json.Unmarshal([]byte(userData), &data); err != nil {
		return User{}, err
	}
	data.Id = userID

	return data, nil
}

func dumpUserData(data User) error {
	userData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = redisClient.HSet(ctx, os.Getenv("REDIS_KEY"), fmt.Sprintf("%s%d", "user:", data.Id), userData).Err()
	if err != nil && err != redis.Nil {
		return err
	}
	return nil
}
