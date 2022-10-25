package util

import (
	"log"
	"os"
)

var workersPool = make(chan struct{}, 8)

func SpawnWorkers(queue *Queue) {
	for {
		if queue.IsEmpty() {
			continue
		}
		workersPool <- struct{}{}
		go func(job any) {

			worker(job)
			<-workersPool
		}(queue.Dequeue())

	}
}

func worker(job any) {
	// apply filter
	// TODO : fix queue to accept data of type handlers.job so that i don't need type assertion

	for i := 0; i < len(job.(Job).Images); i++ {
		image, err := os.Open(job.(Job).Images[i].Path)
		if err != nil {
			log.Fatal(err)
		}

		err = ApplyFilter(image, job.(Job).Filter, job.(Job).Uid, job.(Job).Images[i].Name)
		if err != nil {
			log.Fatal(err)
		}
	}

	getImageNames := func() []string {
		images := []string{}
		for _, imageName := range job.(Job).Images {
			images = append(images, imageName.Name)
		}
		return images
	}

	// archive images
	err := Archive(getImageNames(), job.(Job).Uid)
	if err != nil {
		log.Fatal(err)
	}

}
