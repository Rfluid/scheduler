package redis_scheduler

import (
	"context"
	"time"
)

// Schedules the next TryDequeue to given date.
func (w *Worker) ScheduleDequeue(
	executeAt time.Time, // Next TryDequeue date.
	ctx context.Context, // Redis context passed to TryDequeue.
) error {
	w.timerMu.Lock()
	defer w.timerMu.Unlock()

	// Stop the existing timer scheduler if running
	if w.timer != nil {
		w.timer.Stop()
	}

	// Calculate the duration until the next execution time
	durationUntilExecuteAt := time.Until(executeAt)

	// Replace the old timer with the new one
	w.timer = time.AfterFunc(durationUntilExecuteAt, func() {
		go w.TryDequeue(ctx)
	})

	return nil
}
