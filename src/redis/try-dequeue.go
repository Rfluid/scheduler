package redis_scheduler

import (
	"context"
	"errors"
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
	acquired, err := w.acquireLock(ctx)
	if err != nil {
		return err
	}
	if !acquired {
		return errors.New("could not acquire lock")
	}

	// Get the first element
	firstElement, err := w.First(ctx)
	if err != nil {
		w.releaseLock(ctx)
		return err
	}

	// Convert the score to a time.Time object (assuming it's a Unix timestamp)
	firstElementTime := common_service.FloatToDate(firstElement.Score)

	// Check if the time of the first element is before or equal to current time
	currentTime := time.Now()
	if firstElementTime.After(currentTime) {
		w.releaseLock(ctx)
		// Schedule replacement timer job
		return w.ScheduleDequeue(firstElementTime, ctx)
	}

	errCh := make(chan error)
	var errWg sync.WaitGroup
	errWg.Add(3)
	go func() {
		defer w.releaseLock(ctx)
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
	_, err := w.redisClient.ZRemRangeByRank(ctx, w.redisListKey, 0, 0).Result()
	if err != nil {
		return err
	}

	return nil
}
