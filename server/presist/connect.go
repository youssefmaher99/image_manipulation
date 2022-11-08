package presist

import "github.com/go-redis/redis/v9"

var rds *redis.Client = ConnectToRedis()

func ConnectToRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
