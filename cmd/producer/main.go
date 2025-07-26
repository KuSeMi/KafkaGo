package main

import (
	"fmt"
	k "kafkagolang/internal/kafka"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	topic = "my-topic" // topic name
	// topic2 = "my-topic2" // uncomment this line to use a different topic
	numberOfKeys = 20
)

var address = []string{"localhost:9091", "localhost:9092", "localhost:9093"}

func main() {
	p, err := k.NewProducer(address)
	if err != nil {
		logrus.Fatalf("failed to create producer: %v", err)
		return
	}
	defer p.Close()

	keys := generateUUIDString() // Generate UUIDs for keys
	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("Kafka message %d", i)
		// key := fmt.Sprintf("%d", i%3) // Using modulo to create a key for partitioning // murmur2 hash
		key := keys[i%numberOfKeys] // Use the generated UUIDs as keys
		if err := p.Produce(msg, topic, key, time.Now()); err != nil {
			logrus.Errorf("failed to produce message: %v", err)
		}
	}
}

func generateUUIDString() [numberOfKeys]string {
	var uuids [numberOfKeys]string
	for i := 0; i < numberOfKeys; i++ {
		uuids[i] = uuid.NewString() // Generate a new UUID string
	}
	return uuids
}
