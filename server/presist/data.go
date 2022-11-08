package presist

import (
	"context"
	"fmt"
	"log"
	"server/util"
)

func AddJob(job util.Job) {
	ctx := context.Background()
	for _, image := range job.Images {

		// create hash for each image
		err := rds.HSet(ctx, job.Uid+"-"+image.Name, "name", image.Name, "path", image.Path).Err()
		if err != nil {
			// panic(err)
			fmt.Println(err)
		}

		// create list and push image hashes
		err = rds.RPush(ctx, "images:"+job.Uid, job.Uid+"-"+image.Name).Err()
		if err != nil {
			panic(err)
		}

	}

	// create hash and add list
	err := rds.HSet(ctx, "job:"+job.Uid, "uuid", job.Uid, "filter", job.Filter, "ttl", job.TTl, "images", "images:"+job.Uid, "completed", "0").Err()
	if err != nil {
		panic(err)
	}
}

func markJobAsCompleted(jobId string) {
	ctx := context.Background()
	err := rds.HSet(ctx, "job:"+jobId, "completed", "1").Err()
	if err != nil {
		log.Fatal(err)
	}
}
