package worker

import (
	"fmt"
	"os"
	"server/logger"
	"server/notification"
	"server/presist"
	"server/queue"
	"server/util"
	"time"
)

var workersPool = make(chan struct{}, 8)

func SpawnWorkers(queue *queue.Queue) {
	for {
		if queue.IsEmpty() {
			continue
		}
		workersPool <- struct{}{}
		go func(job util.Job) {
			fmt.Println("new worker started processing")
			worker(job)
			<-workersPool
		}(queue.Dequeue())

	}
}

func worker(job util.Job) {

	// mark job as started to process in redis
	presist.UpdateJobKey(job.Uid, "started-processing", "1")

	// apply filter
	for i := 0; i < len(job.Images); i++ {
		image, err := os.Open(job.Images[i].Path)
		if err != nil {
			logger.MyLog.Fatal(err)
		}

		err = util.ApplyFilter(image, job.Filter, job.Uid, job.Images[i].Name)
		if err != nil {
			logger.MyLog.Fatal(err)
		}
	}

	getImageNames := func() []string {
		images := []string{}
		for _, imageName := range job.Images {
			images = append(images, imageName.Name)
		}
		return images
	}

	// archive images
	err := util.Archive(getImageNames(), job.Uid)
	if err != nil {
		logger.MyLog.Fatal(err)
	}

	fmt.Println("worker finished processing")

	// mark job as completed in redis
	presist.UpdateJobKey(job.Uid, "completed", "1")

	// add expiration to job
	presist.AddExpirationToJob(job.Uid, 30*time.Second)

	// notify user that work is done if user is online
	notification.NotificationChans[job.Uid] <- struct{}{}

}
