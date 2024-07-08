package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	redis_scheduler "github.com/Rfluid/scheduler/src/redis"
	"github.com/redis/go-redis/v9"
)

var (
	ctx                        = context.Background() // Context passed to Redis client
	futureOffset time.Duration = 5                    // First callback time offset.
	sleepTime    time.Duration = 10
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
	worker := redis_scheduler.Create(client, "wList")

	executionCounter := 0
	var propagateMu sync.Mutex
	// Set appropriate callback.
	// This callback will be called when the worker is ready to process the data.
	//
	// You might want to use this callback to propagate the data to another service that will effectively use the data. It is not recommended to do heavy processing here.
	worker.SetCallback(
		func(data redis.Z) error {
			propagateMu.Lock()
			defer propagateMu.Unlock()
			fmt.Println("propagating data.")
			fmt.Println(data)
			executionCounter++
			return nil
		},
	)

	// Setting execution date.
	currentTime := time.Now()
	future := currentTime.Add(futureOffset * time.Second)

	// Scheduling dequeues.
	worker.InsertSortedByDate(
		map[string]string{"third": "print"}, future, ctx,
	)

	worker.InsertSortedByDate(
		map[string]string{"first": "print"}, currentTime.Add(3*time.Second), ctx,
	)
	worker.InsertSortedByDate(
		map[string]string{"second": "print"}, currentTime.Add(3*time.Second), ctx,
	)

	worker.InsertSortedByDate(
		map[string]string{"inserting": "equalElements"}, currentTime.Add(8*time.Second), ctx,
	)
	worker.InsertSortedByDate(
		map[string]string{"inserting": "equalElements"}, // This element will replace the score of equal element.
		currentTime.Add(4*time.Second),
		ctx,
	)

	worker.InsertSortedByDate(
		map[string]string{"fourth": "print"}, future.Add(1*time.Second), ctx,
	)

	fmt.Println("Created worker.")

	time.Sleep(sleepTime * time.Second)
	fmt.Println()
	fmt.Printf("Terminating the program with %d tasks executed.\n", executionCounter)
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
