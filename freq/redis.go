package freq

import (
	"CoinRobot/logger"
	"context"
	"github.com/go-redis/redis/v8"
)

var freqer *redis.Client
var reidsCtx = context.Background()
var log = logger.NewLog()

func init() {
	//连接redis
	freqer = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       2,
	})
}
