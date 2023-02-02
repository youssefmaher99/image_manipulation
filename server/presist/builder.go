package presist

import (
	"server/data"
	"server/models"
	"server/queue"
)

func Builder(queue *queue.Queue[models.Job]) {
	data.InMemoryUUID = GetAllUUID()
	data.InMemoryArchives = GetAllArchives()
	jobs := GetAllJobs()
	for _, job := range jobs {
		queue.Enqueue(job)
	}
}
