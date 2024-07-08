package redis_scheduler

import (
	"context"
	"sync"
	"time"

	common_service "github.com/Rfluid/scheduler/src/common/service"
)

// Dequeues if the current date is after the date of the first element's score.
// Otherwise, schedules the dequeue.
func (w *Worker) TryDequeue(
	ctx context.Context, // Context passed to redis.
) error {
	// Acquiring redis lock
	w.setMu.Lock()

	// Get the first element
	firstElement, err := w.First(ctx)
	if err != nil {
		defer w.setMu.Unlock()
		return err
	}

	// Convert the score to a time.Time object (assuming it's a Unix timestamp)
	firstElementTime := common_service.FloatToDate(firstElement.Score)

	// Check if the time of the first element is before or equal to current time
	currentTime := time.Now()
	if firstElementTime.After(currentTime) {
		defer w.setMu.Unlock()
		return w.ScheduleDequeue(firstElementTime, ctx)
	}

	errCh := make(chan error)
	var errWg sync.WaitGroup
	errWg.Add(3)
	go func() {
		defer w.setMu.Unlock()
		errCh <- w.Dequeue(ctx)

		score, err := w.FirstScore(ctx)
		errCh <- err

		if score == nil {
			errCh <- nil
			return
		}

		firstElementTime := common_service.FloatToDate(*score)
		errCh <- w.ScheduleDequeue(firstElementTime, ctx)
	}()

	errWg.Add(1)
	go func() {
		errCh <- w.Callback(firstElement)
	}()

	go func() {
		defer close(errCh)
		errWg.Wait()
	}()

	for err := range errCh {
		errWg.Done()
		if err != nil {
			return err
		}
	}

	return nil
}

// Removes the first element from the sorted set.
func (w *Worker) Dequeue(ctx context.Context) error {
	_, err := w.redisClient.ZRemRangeByRank(ctx, w.RedisSetKey, 0, 0).Result()
	if err != nil {
		return err
	}

	return nil
}
