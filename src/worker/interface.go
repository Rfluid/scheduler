package worker

import (
	"time"
)

// Any data structure that is a Worker must implement this interface.
type Worker[T, U any] interface {
	InsertSortedByDate(data any, executeAt time.Time, ctx T) error // Inserts data in a sorted manner and schedules a new dequeue if necessary.
	TryDequeue(ctx T) error                                        // If the date is past the first list element date, it dequeues, executes the callback, and schedules the next dequeue. Otherwise, it just schedules the dequeue.
	Dequeue(ctx T) error                                           // Dequeue operation executed by TryDequeue.
	First(ctx T) (U, error)

	ScheduleDequeue(executeAt time.Time, ctx T) error // Replaces the current timer with a new time until dequeue.

	Callback(storedData U) error                    // The callback executed for each dequeue. You might want to use this callback to propagate the data to another service that will effectively use the data. It is not recommended to do heavy processing here.
	SetCallback(cbk func(storedData U) error) error // Sets the callback executed for each dequeue.
}

// TODO: Add more methods for manipulating the Worker queue.
