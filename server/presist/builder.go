package presist

import (
	"server/data"
	"server/queue"
)

func Builder(queue *queue.Queue) {
	data.InMemoryUUID = GetAllUUID()
	data.InMemoryArchives = GetAllArchives()
	jobs := GetAllJobs()
	for _, job := range jobs {
		queue.Enqueue(job)
	}
	// push each of these jobs to queue and worker will begin working
}
