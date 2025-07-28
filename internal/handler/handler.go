package handler

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

// HandleMessage обробляє отримане повідомлення з Kafka
func (h *Handler) HandleMessage(message []byte, topic kafka.TopicPartition, consumerNumber int) error {
	// Тут відбувається основна бізнес-логіка обробки повідомлення
	// В реальному додатку тут може бути:
	// - Парсинг JSON/XML
	// - Збереження в базу даних
	// - Виклик зовнішніх API
	// - Валідація даних
	// - Трансформація повідомлень
	
	// Поки що просто логуємо отримане повідомлення з інформацією про партицію
	logrus.Infof("Received message: %s at topic: %s, partition: %d, consumer №: %d", 
		message, *topic.Topic, topic.Partition, consumerNumber)
	return nil
}
