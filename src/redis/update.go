package redis_scheduler

import (
	"github.com/redis/go-redis/v9"
)

// Updates redisClient and isn't thread safe.
// You can't do basic operations like insert sorted if your client is nil. Note that while correct uses of nil client do exist, they are rare.
func (r *RedisData) UpdateRedisClient(cmdable redis.Cmdable) {
	r.redisClient = cmdable
}

// Updates redisSetKey and isn't thread safe.
func (r *RedisData) UpdateSetKey(redisSetKey string) {
	r.RedisSetKey = redisSetKey
}

// Executes thread safe stop at timer.
func (w *Worker) StopTimer() bool {
	w.timerMu.Lock()
	defer w.timerMu.Unlock()
	w.timerMu.TryLock()

	return w.timer.Stop()
}
