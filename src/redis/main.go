package redis_scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/Rfluid/scheduler/src/worker"
	"github.com/redis/go-redis/v9"
)

// Manages Redis data to handle workers.
type RedisData struct {
	redisClient  *redis.Client
	redisListKey string // Key of the sorted list in Redis.
	redisLockKey string // Key to use for locking.
}

type Worker struct {
	timer    *time.Timer              // Used to schedule dequeues. Only one dequeue is scheduled at a time.
	timerMu  sync.Mutex               // Prevents data races when changing the timer.
	callback func(data redis.Z) error // Propagates the data after dequeuing.
	RedisData
}

// Ensures that Redis Worker implements the Worker interface.
var _ worker.Worker[context.Context, redis.Z] = &Worker{}

// Creates a Redis worker.
func Create(
	redisClient *redis.Client,
	redisListKey string, // Key of the sorted list in Redis.
	redisLockKey string, // Key to use for locking.
) *Worker {
	w := Worker{
		RedisData: RedisData{
			redisClient:  redisClient,
			redisListKey: redisListKey,
			redisLockKey: redisLockKey,
		},
	}

	w.SetCallback(func(data redis.Z) error { return nil })

	return &w
}
