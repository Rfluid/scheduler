package redis_scheduler

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	common_service "github.com/Rfluid/scheduler/src/common/service"
	"github.com/redis/go-redis/v9"
)

// Inserts element sorted by date into the Redis list.
//
// The date is converted to Unix using common_service.DateToFloat so we can manage Redis scores.
//
// If the executeAt score is smaller than the score of the first element, the next dequeue is rescheduled using ScheduleDequeue.
func (w *Worker) InsertSortedByDate(
	data any, // Data to schedule propagation.
	executeAt time.Time, // Execution date of data propagation.
	ctx context.Context, // Used in redis operations.
) error {
	firstElementBeforeScoreCh := make(chan *float64)
	scoreCh := make(chan float64, 2)
	dataStrCh := make(chan string)
	errCh := make(chan error)
	var errWg sync.WaitGroup

	// Acquiring redis lock
	acquired, err := w.acquireLock(ctx)
	if err != nil {
		return err
	}
	if !acquired {
		return errors.New("could not acquire lock")
	}
	// Ensure the lock is released at the end
	defer w.releaseLock(ctx)

	// Get the score of the first element before insertion
	errWg.Add(1)
	go func() {
		firstElementBeforeScore, err := w.FirstScore(ctx)
		errCh <- err
		firstElementBeforeScoreCh <- firstElementBeforeScore
	}()

	// Serialize the data to JSON
	errWg.Add(1)
	go func() {
		dataBytes, err := json.Marshal(data)
		errCh <- err

		// Use the timestamp as the score
		go func() {
			for i := 0; i < 2; i++ {
				scoreCh <- common_service.DateToFloat(executeAt)
			}
		}()

		// Convert the JSON bytes to string
		dataStrCh <- string(dataBytes)
	}()

	// Check if the score of the first element will change
	errWg.Add(1)
	go func() {
		var err error = nil
		prevScoreP := <-firstElementBeforeScoreCh
		if prevScoreP == nil || *prevScoreP > <-scoreCh {
			err = w.ScheduleDequeue(executeAt, ctx)
		}
		errCh <- err
	}()

	// Insert the data into the sorted set
	errWg.Add(1)
	go func() {
		_, err = w.redisClient.ZAdd(ctx, w.redisListKey, redis.Z{
			Score:  <-scoreCh,
			Member: <-dataStrCh,
		}).Result()
		errCh <- err
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

// Returns the score pointer of the first element if the queue has any elements. Otherwise, returns nil.
func (w *Worker) FirstScore(ctx context.Context) (*float64, error) {
	firstElement, err := w.redisClient.ZRangeWithScores(ctx, w.redisListKey, 0, 0).Result()
	if err != nil {
		return nil, err
	}
	if len(firstElement) == 0 {
		return nil, nil
	}
	return &firstElement[0].Score, nil
}

// Returns the first element of the queue. If the length is zero, returns an empty element struct.
func (w *Worker) First(ctx context.Context) (redis.Z, error) {
	firstElement, err := w.redisClient.ZRangeWithScores(ctx, w.redisListKey, 0, 0).Result()
	if err != nil {
		return redis.Z{}, err
	}
	if len(firstElement) == 0 {
		return redis.Z{}, nil
	}
	return firstElement[0], nil
}
