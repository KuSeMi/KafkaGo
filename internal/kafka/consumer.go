package kafka

import (
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

const (
	sessionTimeout = 8000 // ms
	noTimeout      = -1
)

type Handler interface {
	HandleMessage(message []byte, topic kafka.TopicPartition, consumerNumber int) error
}

type Consumer struct {
	consumer       *kafka.Consumer
	handler        Handler
	stop           bool
	consumerNumber int // number of consumers in the group
}

// NewConsumer створює новий Kafka consumer з налаштуваннями для групової роботи
func NewConsumer(handler Handler, address []string, topic, consumerGroup string, consumerNumber int) (*Consumer, error) {
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(address, ","), // Список брокерів
		"group.id":                 consumerGroup,              // Ідентифікатор групи споживачів
		"session.timeout.ms":       sessionTimeout,            // Таймаут сесії (8 сек)
		"enable.auto.offset.store": false,                     // Вимикаємо автоматичне збереження offset'ів
		"enable.auto.commit":       true,                      // Увімкнути автоматичний коміт offset'ів
		"auto.commit.interval.ms":  5000,                      // Інтервал автокоміту (5 сек)
		"auto.offset.reset":        "earliest",                // Починати з найраніших повідомлень при першому запуску
	}
	
	c, err := kafka.NewConsumer(cfg)
	if err != nil {
		return nil, err
	}

	// Підписуємося на топік
	if err = c.Subscribe(topic, nil); err != nil {
		return nil, err
	}

	return &Consumer{
		handler:        handler,
		consumer:       c,
		consumerNumber: consumerNumber,
	}, nil
}

// Start запускає основний цикл споживання повідомлень
func (c *Consumer) Start() {
	for {
		// Перевіряємо чи потрібно зупинити consumer
		if c.stop {
			break
		}

		// Читаємо повідомлення з Kafka (блокуючий виклик)
		kafkaMsg, err := c.consumer.ReadMessage(noTimeout)
		if err != nil {
			logrus.Error(err)
		}

		// Пропускаємо пусті повідомлення
		if kafkaMsg == nil {
			continue
		}

		// Обробляємо повідомлення через handler
		if err = c.handler.HandleMessage(kafkaMsg.Value, kafkaMsg.TopicPartition, c.consumerNumber); err != nil {
			logrus.Error(err)
			continue
		}

		// Зберігаємо offset оброблених повідомлень для можливості відновлення
		if _, err = c.consumer.StoreMessage(kafkaMsg); err != nil {
			logrus.Error(err)
			continue
		}
	}
}

// Stop коректно зупиняє consumer з комітом offset'ів
func (c *Consumer) Stop() error {
	c.stop = true
	
	// Комітимо всі збережені offset'и перед закриттям
	if _, err := c.consumer.Commit(); err != nil {
		return err
	}

	logrus.Infof("Commited offset")
	return c.consumer.Close()
}
