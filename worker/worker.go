package worker

import (
	"github.com/MatticNote/MatticNote/config"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var Worker *work.WorkerPool

func InitWorker() {
	redisPool := &redis.Pool{
		MaxIdle:   config.Config.Job.MaxIdle,
		MaxActive: config.Config.Job.MaxActive,
		Wait:      true,
		Dial:      config.GetRedisPool,
	}
	Worker = work.NewWorkerPool(struct{}{}, uint(config.Config.Job.MaxActive), "mn_worker", redisPool)
}
