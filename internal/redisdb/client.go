package redisdb

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var (
	Ctx = context.Background()
	rdb *redis.Client
)

func Init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	if err := rdb.Ping(Ctx).Err(); err != nil {
		log.Fatalf("❌ Redis connection error: %v", err)
	}
	log.Println("✅ Connected to Redis")
}

func GetClient() *redis.Client {
	return rdb
}
