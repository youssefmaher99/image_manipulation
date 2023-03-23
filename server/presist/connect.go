package presist

import (
	"context"
	"server/logger"

	"github.com/go-redis/redis/v9"
)

var rds *redis.Client

func init() {
	rds = ConnectToRedis()
	// initiateDestroyerWorker()
	// initiateDeadRefsWorker()
}

func ConnectToRedis() *redis.Client {
	conn := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if _, err := conn.Ping(context.Background()).Result(); err != nil {
		logger.MyLog.Fatal(err)
	}

	logger.MyLog.Println("REDIS CONNECTED")
	return conn
}
