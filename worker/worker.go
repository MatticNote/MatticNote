package worker

import (
	"github.com/MatticNote/MatticNote/config"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var (
	Enqueue *work.Enqueuer
	Worker  *work.WorkerPool
	Client  *work.Client
)

const (
	workerNamespace = "mn_worker"
)

const (
	JobInboxProcess = "inbox_process"
	JobDeliver      = "deliver"
)

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
	Enqueue = work.NewEnqueuer(workerNamespace, redisPool)
	Client = work.NewClient(workerNamespace, redisPool)
}

func InitWorker() {
	redisPool := getRedisPool()
	Worker = work.NewWorkerPool(Context{}, uint(config.Config.Job.MaxActive), workerNamespace, redisPool)

	Worker.JobWithOptions(
		JobInboxProcess,
		work.JobOptions{
			MaxFails: 10,
		},
		(*Context).ProcessInbox,
	)
	Worker.JobWithOptions(
		JobDeliver,
		work.JobOptions{
			MaxFails: 10,
		},
		(*Context).Deliver,
	)
}
