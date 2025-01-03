package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

func SendTopic(topic string, msg []byte) {

	writer := &kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092"),
		Topic:                  topic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           1 * time.Second,
		RequiredAcks:           kafka.RequireNone,
		AllowAutoTopicCreation: true,
	}
	defer writer.Close()

	err := writer.WriteMessages(
		context.Background(),
		kafka.Message{Value: msg},
	)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	fmt.Printf("topic:%s 消息发送成功 \n", topic)
}

func main() {
	SendTopic("test11", []byte("枫枫1"))
	SendTopic("test11", []byte("知道1"))
}
