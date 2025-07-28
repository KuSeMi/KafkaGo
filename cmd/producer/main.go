package main

import (
	"fmt"
	k "kafkagolang/internal/kafka"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	topic = "my-topic" // Назва топіка - логічна категорія для групування повідомлень
	// topic2 = "my-topic2" // Можна використовувати різні топіки для різних типів даних
	numberOfKeys = 20 // Кількість унікальних ключів для розподілу повідомлень по партиціях
)

// Адреси брокерів Kafka - кластер з 3 брокерів для забезпечення відмовостійкості
var address = []string{"localhost:9091", "localhost:9092", "localhost:9093"}

func main() {
	// Створюємо producer для відправки повідомлень до Kafka
	p, err := k.NewProducer(address)
	if err != nil {
		logrus.Fatalf("failed to create producer: %v", err)
		return
	}
	// Обов'язково закриваємо producer після завершення роботи
	defer p.Close()

	// Генеруємо унікальні UUID ключі для розподілу повідомлень по партиціях
	keys := generateUUIDString()
	
	// Відправляємо 100 повідомлень
	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("Kafka message %d", i)
		
		// Використовуємо модуло для циклічного вибору ключа
		// Ключ визначає, в яку партицію потрапить повідомлення (через хешування)
		key := keys[i%numberOfKeys]
		
		// Відправляємо повідомлення з поточним часом
		if err := p.Produce(msg, topic, key, time.Now()); err != nil {
			logrus.Errorf("failed to produce message: %v", err)
		}
	}
}

// generateUUIDString генерує масив унікальних UUID для використання як ключі повідомлень
// Ключі важливі для партиціонування - повідомлення з однаковим ключем завжди потраплять в одну партицію
func generateUUIDString() [numberOfKeys]string {
	var uuids [numberOfKeys]string
	for i := 0; i < numberOfKeys; i++ {
		uuids[i] = uuid.NewString() // Генеруємо новий UUID рядок
	}
	return uuids
}
