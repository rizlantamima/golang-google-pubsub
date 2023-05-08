# Golang handle Google Pub/Sub

Google Pub/Sub is an asynchronous messages service, or maybe you can call it as message broker like RabbitMQ / Kafka 

in this project I made 2 different cli applications, one of which acts as a message **publisher**, application that sends messages on the topic, the other one as a subscriber, application that subscribing messages from a topic.

## Preparation : 

1. Make sure that you already have go on you pc
2. You must have google projects that enable google pubsub Api.
3. You must create at least 1 topic and 1 subscription on google pubsub api.
4. Save the **name** of your google project on `.env` file as `GOOGLE_PROJECT_ID`


## Publisher Apps Installation

1. Change your working directory to the *publisher*.
```bash
cd publisher
```
2. Install any go module on this app
```bash
go get
```
3. Run this publisher app
```bash
go run main.go
```

## Subscriber Apps Installation

2. Change your working directory to the *subscriber*.
```bash
cd subscriber
```
2. Install any go module on this app
```bash
go get
```
3. Run this subscriber app
```bash
go run main.go
```