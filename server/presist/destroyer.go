package presist

import (
	"context"
	"fmt"
	"server/data"
	"server/logger"
	"time"
)

// a destroyer worker is a worker responsible for checking keys expiration in redis and the destruction of any refrences for keys on server such as archives or images

func initiateDestroyerWorker() {
	go func() {
		for {
			// check all job keys in redis against archives
			for k := range data.InMemoryArchives {
				_, err := rds.HGet(context.Background(), fmt.Sprintf("job:%s", k), "uuid").Result()
				if err != nil {

					// key doesn't exist anymore in redis check for other refrences and delete them
					if err.Error() == "redis: nil" {

						// remove from in mem arch
						data.InMemoryArchives.Remove(k)
						// remove from mem uui
						data.InMemoryUUID.Remove(k)

						// remove all matching files in filtered dir
						// remove all matching files in uploaded dir
						// remove it from archives
						data.RemoveFromDisk(k)

						// remove all redis refrences (job,images,image)
						DeleteJob(k)
						continue
					}
					logger.MyLog.Fatal(err)
				}
			}
			time.Sleep(time.Second * 5)
		}
	}()
}
