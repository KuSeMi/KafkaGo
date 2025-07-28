package main

import (
	"kafkagolang/internal/handler"
	"kafkagolang/internal/kafka"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

const (
	topic         = "my-topic"          // Назва топіка з якого читаємо повідомлення
	consumerGroup = "my-consumer-group" // Група споживачів - дозволяє розподіляти навантаження між кількома consumer'ами
)

// Адреси брокерів Kafka - той самий кластер що і для producer'а
var address = []string{"localhost:9091", "localhost:9092", "localhost:9093"}

func main() {
	// Створюємо обробник повідомлень
	h := handler.NewHandler()
	
	// Створюємо 3 consumer'и в одній групі - вони автоматично розподілять партиції між собою
	// Кожен consumer отримає свій номер для ідентифікації в логах
	c1, err := kafka.NewConsumer(h, address, topic, consumerGroup, 1)
	if err != nil {
		logrus.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	c2, err := kafka.NewConsumer(h, address, topic, consumerGroup, 2)
	if err != nil {
		logrus.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	c3, err := kafka.NewConsumer(h, address, topic, consumerGroup, 3)
	if err != nil {
		logrus.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	// Запускаємо всіх consumer'ів в окремих горутинах для паралельної обробки
	go c1.Start()
	go c2.Start()
	go c3.Start()
	
	// Налаштовуємо graceful shutdown - чекаємо сигнал завершення
	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)

	// Блокуємося до отримання сигналу завершення
	<-signChan
	
	// Коректно зупиняємо всіх consumer'ів з комітом offset'ів
	logrus.Fatal(c1.Stop(), c2.Stop(), c3.Stop())
}
