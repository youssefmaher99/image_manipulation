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
}
