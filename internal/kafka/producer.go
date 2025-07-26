package kafka

import (
	"fmt"
	"strings"
	"time"

	"errors"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	flushTimeout = 5000 // ms
)

var (
	errUnknowType = errors.New("unknow event type")
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(address []string) (*Producer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(address, ","),
	}

	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(message string, topic string, key string, t time.Time) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value:     []byte(message),
		Key:       []byte(key),
		Timestamp: t,
	}
	kafkaChan := make(chan kafka.Event)
	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}
	e := <-kafkaChan

	switch ev := e.(type) {
	case *kafka.Message:
		return nil
		// The client will automatically try to recover from all errors.
	case kafka.Error:
		return fmt.Errorf("kafka error: %w", ev)
	default:
		return errUnknowType
	}

}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout) // чекаєм всі повідомлення в черзі коли вони будуть відправлені (бллокує до flushTimeout або до завершення всіх повідомлень)
	p.producer.Close()
}
