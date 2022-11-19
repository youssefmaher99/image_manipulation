package presist

import (
	"context"
	"server/logger"
	"server/util"
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
	err := rds.HSet(ctx, "job:"+job.Uid, "uuid", job.Uid, "filter", job.Filter, "ttl", job.TTl, "images", "images:"+job.Uid, "completed", "0", "started-processing", "0").Err()
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
}
