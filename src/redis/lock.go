package redis_scheduler

import "context"

// Acquires Redis lock.
func (d *RedisData) acquireLock(ctx context.Context) (bool, error) {
	status, err := d.redisClient.SetNX(ctx, d.redisLockKey, "locked", 0).Result()
	return status, err
}

// Releases Redis lock.
func (d *RedisData) releaseLock(ctx context.Context) error {
	_, err := d.redisClient.Del(ctx, d.redisLockKey).Result()
	return err
}
