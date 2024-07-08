package redis_scheduler

import (
	"github.com/redis/go-redis/v9"
)

// Updates redisClient and isn't thread safe.
// You can't do basic operations like insert sorted if your client is nil. Note that while correct uses of nil client do exist, they are rare.
func (r *RedisData) UpdateRedisClient(redisClient *redis.Client) {
	r.redisClient = redisClient
}

// Updates redisListKey and isn't thread safe.
func (r *RedisData) UpdateListKey(redisListKey string) {
	r.redisListKey = redisListKey
}

// Updates redisLockKey and isn't thread safe.
func (r *RedisData) UpdateLockKey(redisLockKey string) {
	r.redisLockKey = redisLockKey
}

// Executes thread safe stop at timer.
func (w *Worker) StopTimer() bool {
	w.timerMu.Lock()
	defer w.timerMu.Unlock()
	w.timerMu.TryLock()

	return w.timer.Stop()
}
