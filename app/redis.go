package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func createRedisClient(redisURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		fmt.Println("Invalid Redis URL:", err)
		return nil, err
	}

	client := redis.NewClient(opt)
	_, err = client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return nil, err
	}

	fmt.Println("Connected to Redis.")

	_, err = client.ConfigSet(ctx, "notify-keyspace-events", "Ex").Result()
	if err != nil {
		fmt.Println("Failed to set notify-keyspace-events:", err)
		return nil, err
	}

	return client, nil
}
