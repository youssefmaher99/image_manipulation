package presist

import (
	"context"
	"server/logger"
	"server/util"
	"time"
)

func AddJob(job util.Job) {
	ctx := context.Background()
	for _, image := range job.Images {

		// create hash for each image
		err := rds.HSet(ctx, job.Uid+"-"+image.Name, "name", image.Name, "path", image.Path).Err()
		if err != nil {
			logger.MyLog.Fatal(err)
		}

		// create list and push image hashes
		err = rds.RPush(ctx, "images:"+job.Uid, job.Uid+"-"+image.Name).Err()
		if err != nil {
			logger.MyLog.Fatal(err)
		}

	}

	// create hash and add list
	err := rds.HSet(ctx, "job:"+job.Uid, "uuid", job.Uid, "filter", job.Filter, "images", "images:"+job.Uid, "completed", "0", "started-processing", "0").Err()
	if err != nil {
		logger.MyLog.Fatal(err)
	}
}

func AddExpirationToJob(jobId string, timeToExpire time.Duration) {
	ctx := context.Background()
	err := rds.Expire(ctx, "job:"+jobId, timeToExpire).Err()
	if err != nil {
		logger.MyLog.Fatal(err)
	}
}

func UpdateJobKey(jobId string, key string, value string) {
	ctx := context.Background()
	err := rds.HSet(ctx, "job:"+jobId, key, value).Err()
	if err != nil {
		logger.MyLog.Fatal(err)
	}
}

func DeleteJob(jobId string) {
	ctx := context.Background()

	err := rds.Del(ctx, "job:"+jobId).Err()
	if err != nil {
		logger.MyLog.Fatal(err)
	}

	// get all items in list then remove one by one then remove list
	images, err := rds.LRange(ctx, "images:"+jobId, 0, -1).Result()
	if err != nil {
		logger.MyLog.Fatal(err)
	}

	// delete each image
	for _, image := range images {
		err = rds.Del(ctx, image).Err()
		if err != nil {
			logger.MyLog.Fatal(err)
		}
	}

	// remove list
	err = rds.Del(ctx, "images:"+jobId).Err()
	if err != nil {
		logger.MyLog.Fatal(err)
	}

}
