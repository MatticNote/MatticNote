package worker

import (
	"github.com/MatticNote/MatticNote/config"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var (
	Enqueue *work.Enqueuer
	Worker  *work.WorkerPool
)

const workerName = "mn_worker"

func InitWorker() {
	redisPool := &redis.Pool{
		MaxIdle:   config.Config.Job.MaxIdle,
		MaxActive: config.Config.Job.MaxActive,
		Wait:      true,
		Dial:      config.GetRedisPool,
	}
	Enqueue = work.NewEnqueuer(workerName, redisPool)
	Worker = work.NewWorkerPool(Context{}, uint(config.Config.Job.MaxActive), workerName, redisPool)

	Worker.JobWithOptions("inbox_worker", work.JobOptions{
		MaxFails: 10,
	}, (*Context).ProcessInbox)
}
