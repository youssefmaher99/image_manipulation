package util

import (
	"fmt"
	"log"
	"os"
	"server/notification"
)

var workersPool = make(chan struct{}, 8)

func SpawnWorkers(queue *Queue) {
	for {
		if queue.IsEmpty() {
			continue
		}
		workersPool <- struct{}{}
		go func(job Job) {
			fmt.Println("new worker started processing")
			worker(job)
			<-workersPool
		}(queue.Dequeue())

	}
}

func worker(job Job) {
	// apply filter

	for i := 0; i < len(job.Images); i++ {
		image, err := os.Open(job.Images[i].Path)
		if err != nil {
			log.Fatal(err)
		}

		err = ApplyFilter(image, job.Filter, job.Uid, job.Images[i].Name)
		if err != nil {
			log.Fatal(err)
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
	err := Archive(getImageNames(), job.Uid)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("worker finished processing")

	// notify user that work is done
	notification.NotificationChans[job.Uid] <- struct{}{}

	// mark job as done in redis
	// presist.UpdateJob(job.Uid)
}
