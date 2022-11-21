package presist

import (
	"context"
	"fmt"
	"server/data"
	"server/logger"
	"time"

	"github.com/go-redis/redis/v9"
)

// a destroyer worker is a worker responsible for checking keys expiration in redis and the destruction of any refrences for keys on server such as archives or images

func initiateDestroyerWorker(conn *redis.Client) {
	go func() {
		for {
			// check all job keys in redis against archives
			for k := range data.InMemoryArchives {
				exist, err := conn.Exists(context.Background(), fmt.Sprintf("job:%s", k)).Result()
				if err != nil {
					logger.MyLog.Fatal(err)
				}

				if exist <= 0 {
					fmt.Println("Cleaning")
					// remove from in mem arch
					data.InMemoryArchives.Remove(k)
					// remove from mem uui
					data.InMemoryUUID.Remove(k)

					// remove all matching files in filtered dir
					// remove all matching files in uploaded dir
					// remove it from archives
					data.RemoveFromDisk(k)
				}
			}
			time.Sleep(time.Second * 5)
		}
	}()
}
