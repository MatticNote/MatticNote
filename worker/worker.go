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

func getRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   config.Config.Job.MaxIdle,
		MaxActive: config.Config.Job.MaxActive,
		Wait:      true,
		Dial:      config.GetRedisDial,
	}
}

func InitEnqueue() {
	redisPool := getRedisPool()
	Enqueue = work.NewEnqueuer(workerName, redisPool)
}

func InitWorker() {
	redisPool := getRedisPool()
	Worker = work.NewWorkerPool(Context{}, uint(config.Config.Job.MaxActive), workerName, redisPool)

	Worker.JobWithOptions("inbox_worker", work.JobOptions{
		MaxFails: 10,
	}, (*Context).ProcessInbox)
}
