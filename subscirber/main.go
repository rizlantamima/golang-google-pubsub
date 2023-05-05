package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Subscriber running ...")

	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	projectID := os.Getenv("GOOGLE_PROJECT_ID")
	subID := os.Getenv("GOOGLE_PUBSUB_SUBSCRIPTION_ID")

	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		err = fmt.Errorf("pubsub: NewClient: %v", err)
		panic(err)
	}
	defer pubsubClient.Close()
	stream(ctx, pubsubClient, subID)
}

func stream(ctx context.Context, client *pubsub.Client, subID string) error {

	sub := client.Subscription(subID)

	// Receive messages for 10 seconds, which simplifies testing.
	// Comment this out in production, since `Receive` should
	// be used as a long running operation.
	// ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()

	var received int32
	err := sub.Receive(ctx, pullHandler)
	if err != nil {
		return fmt.Errorf("sub.Receive: %v", err)
	}
	fmt.Printf("Received %d messages\n", received)

	return nil
}

func pullHandler(_ context.Context, msg *pubsub.Message) {
	var received int32
	fmt.Printf("Got message: %q\n", string(msg.Data))
	atomic.AddInt32(&received, 1)
	msg.Ack()
}
