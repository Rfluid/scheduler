package example_redis

import (
	"context"
	"fmt"
	"os"
	"time"

	redis_scheduler "github.com/Rfluid/scheduler/src/redis"
	"github.com/redis/go-redis/v9"
)

var (
	ctx                        = context.Background() // Context passed to Redis client
	futureOffset time.Duration = 5                    // First callback time offset.
	sleepTime    time.Duration = 8
)

func main() {
	fmt.Println("Connecting to redis...")
	client := connectToRedis() // Returns Redis client
	if client == nil {
		fmt.Println("client is nil")
		os.Exit(1)
	}
	fmt.Println("Connected to redis")

	fmt.Println("Creating worker and insert...")
	// Create the worker
	worker := redis_scheduler.Create(client, "wList", "wLock")

	// Set appropriate callback.
	// This callback will be called when the worker is ready to process the data.
	//
	// You might want to use this callback to propagate the data to another service that will effectively use the data. It is not recommended to do heavy processing here.
	worker.SetCallback(
		func(data redis.Z) error {
			fmt.Println("propagating data.")
			return nil
		},
	)

	// Setting execution date.
	currentTime := time.Now()
	future := currentTime.Add(futureOffset * time.Second)

	// Scheduling the first dequeue and callback.
	worker.InsertSortedByDate(
		map[string]string{"first": "data"}, future, ctx,
	)

	// Scheduling the second dequeue and callback 1 second after the first one.
	worker.InsertSortedByDate(
		map[string]string{"another": "object"}, future.Add(1*time.Second), ctx,
	)
	fmt.Println("Created worker and inserted 2 objects.")

	fmt.Println("Sleeping to wait for the callbacks.")
	time.Sleep(sleepTime * time.Second)
	fmt.Println("Done sleeping. Terminating the program with 2 tasks executed.")
}

func connectToRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "redis",
		DB:       0, // use default DB
	})

	// Test the connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Could not connect to Redis: %v\n", err)
		return nil
	}

	return rdb
}
