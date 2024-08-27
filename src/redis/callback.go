package redis_scheduler

import "github.com/redis/go-redis/v9"

// Sets the callback to be executed after the dequeue.
func (w *Worker) SetCallback(
	cbk func(data redis.Z) error,
) error {
	w.callback = cbk
	return nil
}

// Callback executed for each dequeue. You might want to use this callback to propagate the data to another service that will effectively use the data. It is not recommended to do heavy processing here.
func (w *Worker) Callback(data redis.Z) error {
	return w.callback(data)
}

// Sets the callback executed for each dequeue error.
func (w *Worker) SetTryDequeueErrCallback(
	cbk func(error),
) error {
	w.tryDequeueErrCallback = cbk
	return nil
}

func (w *Worker) TryDequeueErrCallback(err error) {
	w.tryDequeueErrCallback(err)
}
