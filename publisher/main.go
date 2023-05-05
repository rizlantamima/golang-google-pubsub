package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
)

func main() {

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

	fmt.Println("****************************************")
	fmt.Println("*                                      *")
	fmt.Println("*                                      *")
	fmt.Println("*      Wellcome Message Publisher      *")
	fmt.Println("*                                      *")
	fmt.Println("*                                      *")
	fmt.Println("****************************************")

	var retry string
	for {
		err = ask(ctx, pubsubClient)
		if err != nil {
			panic(err)
		}

		fmt.Printf("You want to sent message again ? y/n : ")
		fmt.Scanln(&retry)
		retry = strings.TrimSpace(strings.ToUpper(retry))
		if retry != "Y" {
			break
		}
		retry = ""
	}
	fmt.Println("Thank you! bye!")
}

func ask(ctx context.Context, client *pubsub.Client) error {
	topics, err := listTopics(ctx, client)
	if err != nil {
		err = fmt.Errorf("pubsub: NewClient: %v", err)
		return err
	}

	fmt.Printf("\nList of available topics on your project :\n\n")

	for key, value := range topics {
		fmt.Println("  ", key+1, value)
	}

	var topicNumberChoosen int

	var message string
	fmt.Print("\n\nplease select one of the topic numbers above : ")
	_, err = fmt.Scan(&topicNumberChoosen)
	if err != nil {
		return err
	}

	if topicNumberChoosen > len(topics) {
		err := fmt.Errorf("you are choosing invalid / unavailable number")
		return err
	}

	topicChoosen := topics[topicNumberChoosen-1]

	fmt.Printf("Whats message you wanto to sent to the topic %s : ", topicChoosen)
	_, err = fmt.Scanln(&message)
	if err != nil {
		return err
	}

	publish(ctx, client, topicChoosen, message)
	return nil
}

func listTopics(ctx context.Context, client *pubsub.Client) (topics []string, err error) {

	it := client.Topics(ctx)
	for {
		topic, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		topicNameArr := strings.Split(topic.String(), "/")
		topicName := topicNameArr[len(topicNameArr)-1]
		topics = append(topics, topicName)
	}
	return topics, err

}
func publish(ctx context.Context, client *pubsub.Client, topicID, msg string) error {

	t := client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})

	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("pubsub: result.Get: %v", err)
	}
	fmt.Printf("Published a message; msg ID: %v\n", id)
	return nil
}
