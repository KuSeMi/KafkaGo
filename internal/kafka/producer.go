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

// NewProducer створює новий Kafka producer з мінімальною конфігурацією
func NewProducer(address []string) (*Producer, error) {
	conf := &kafka.ConfigMap{
		// bootstrap.servers - список брокерів для початкового підключення
		// Producer автоматично дізнається про інші брокери в кластері
		"bootstrap.servers": strings.Join(address, ","),
	}

	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{producer: p}, nil
}

// Produce відправляє повідомлення до Kafka топіка
func (p *Producer) Produce(message string, topic string, key string, t time.Time) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny, // Дозволяємо Kafka автоматично вибрати партицію на основі ключа
		},
		Value:     []byte(message), // Тіло повідомлення
		Key:       []byte(key),     // Ключ для партиціонування - повідомлення з однаковим ключем потраплять в одну партицію
		Timestamp: t,               // Час створення повідомлення
	}
	
	// Створюємо канал для отримання результату відправки
	kafkaChan := make(chan kafka.Event)
	
	// Асинхронно відправляємо повідомлення
	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}
	
	// Чекаємо результат відправки
	e := <-kafkaChan

	switch ev := e.(type) {
	case *kafka.Message:
		// Повідомлення успішно відправлено
		return nil
	case kafka.Error:
		// Сталася помилка при відправці
		return fmt.Errorf("kafka error: %w", ev)
	default:
		return errUnknowType
	}
}

// Close коректно закриває producer
func (p *Producer) Close() {
	// Flush чекає поки всі повідомлення в черзі будуть відправлені
	// Блокує виконання до flushTimeout або до завершення відправки всіх повідомлень
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}
