package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
)

func main() {
	log.Println("Subscriber running ...")

	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	projectID := os.Getenv("GOOGLE_PROJECT_ID")

	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		err = fmt.Errorf("pubsub: NewClient: %v", err)
		panic(err)
	}
	defer pubsubClient.Close()

	subs, err := listSubscriptions(ctx, pubsubClient)
	if err != nil {
		err = fmt.Errorf("pubsub: NewClient: %v", err)
		panic(err)
	}
	fmt.Printf("\nList of available subscriptions on your project :\n\n")

	for key, value := range subs {
		fmt.Println("  ", key+1, value)
	}

	var subChoosenNumber int
	fmt.Printf("\n\nplease select one of the subscriptions numbers above : 1 - %v \n", len(subs))
	_, err = fmt.Scan(&subChoosenNumber)
	if err != nil {
		panic(err)
	}

	if subChoosenNumber > len(subs) {
		err := fmt.Errorf("you are choosing invalid / unavailable number")
		panic(err)
	}

	subChoosen := subs[subChoosenNumber-1]

	stream(ctx, pubsubClient, subChoosen)
}

func listSubscriptions(ctx context.Context, client *pubsub.Client) (subs []string, err error) {

	it := client.Subscriptions(ctx)
	for {
		s, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Next: %v", err)
		}

		subNameArr := strings.Split(s.String(), "/")
		subName := subNameArr[len(subNameArr)-1]
		subs = append(subs, subName)
	}
	return subs, nil
}

func stream(ctx context.Context, client *pubsub.Client, subID string) error {
	fmt.Printf("\n subscribing message from %s ... \n \n", subID)
	sub := client.Subscription(subID)

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
	fmt.Printf("\n Got message: id(%s) %q, this message are published at %s \n", msg.ID, string(msg.Data), msg.PublishTime)
	atomic.AddInt32(&received, 1)

	var acknowledge string
	fmt.Printf("\n Do you want to acknowledge this message from pubsub ? y/n \n")

	_, err := fmt.Scan(&acknowledge)
	if err != nil {
		panic(err)
	}

	if strings.ToUpper(strings.TrimSpace(acknowledge)) == "Y" {
		fmt.Printf("\n Done \n")
		msg.Ack()
	}
}
