# Apache Kafka Tutorial - Повний гід по основних концепціях

## Зміст
1. [Що таке Apache Kafka?](#що-таке-apache-kafka)
2. [Основні концепції](#основні-концепції)
3. [Архітектура системи](#архітектура-системи)
4. [Як запустити проект](#як-запустити-проект)
5. [Детальний розбір коду](#детальний-розбір-коду)
6. [Блок-схеми взаємодії](#блок-схеми-взаємодії)

## Що таке Apache Kafka?

Apache Kafka - це розподілена платформа для потокової обробки даних, яка дозволяє:
- **Публікувати та підписуватися** на потоки записів (як система повідомлень)
- **Зберігати** потоки записів надійно та відмовостійко
- **Обробляти** потоки записів в реальному часі

### Основні переваги Kafka:
- ⚡ **Висока продуктивність** - мільйони повідомлень на секунду
- 🔄 **Горизонтальне масштабування** - легко додавати нові брокери
- 💾 **Довготривале зберігання** - дані можуть зберігатися днями, тижнями, роками
- 🛡️ **Відмовостійкість** - реплікація даних між брокерами
- 🔧 **Гнучкість** - підтримка різних патернів обміну повідомленнями

## Основні концепції

### 📝 Topic (Топік)
**Топік** - це логічна категорія або канал, куди публікуються повідомлення.

```
Приклад топіків:
├── user-events          # Події користувачів
├── order-processing     # Обробка замовлень  
├── payment-notifications # Сповіщення про платежі
└── system-logs          # Системні логи
```

**Характеристики топіків:**
- Топіки мають унікальні імена в кластері
- Повідомлення в топіку впорядковані за часом
- Топіки можуть мати налаштування retention (час зберігання)

### 🗂️ Partition (Партиція)
**Партиція** - це фізичний розділ топіка для паралелізації та масштабування.

```
Topic: user-events
├── Partition 0: [msg1, msg4, msg7, ...]
├── Partition 1: [msg2, msg5, msg8, ...]
└── Partition 2: [msg3, msg6, msg9, ...]
```

**Важливі особливості партицій:**
- Повідомлення з однаковим ключем завжди потраплять в одну партицію
- Порядок гарантується тільки в межах однієї партиції
- Кількість партицій визначає максимальний паралелізм споживання

### 📤 Producer (Продюсер)
**Продюсер** - це клієнт, який публікує (відправляє) повідомлення до топіків.

**Основні функції:**
- Відправка повідомлень до топіків
- Вибір партиції (автоматично або за ключем)
- Батчинг повідомлень для продуктивності
- Обробка помилок та повторні спроби

### 📥 Consumer (Споживач)
**Споживач** - це клієнт, який читає повідомлення з топіків.

**Основні функції:**
- Читання повідомлень з топіків
- Відстеження позиції читання (offset)
- Групова робота з іншими споживачами

### 👥 Consumer Group (Група споживачів)
**Група споживачів** - це набір споживачів, які спільно обробляють повідомлення з топіка.

```
Topic з 3 партиціями:
├── Partition 0 → Consumer A (group-1)
├── Partition 1 → Consumer B (group-1)  
└── Partition 2 → Consumer C (group-1)
```

**Переваги груп:**
- Автоматичний розподіл партицій між споживачами
- Відмовостійкість - при падінні одного споживача, його партиції перерозподіляються
- Масштабування - можна додавати/видаляти споживачів

### 📍 Offset (Зміщення)
**Offset** - це унікальний ідентифікатор позиції повідомлення в партиції.

```
Partition 0:
[msg1:0] [msg2:1] [msg3:2] [msg4:3] [msg5:4]
                           ↑
                    Current offset: 3
```

**Типи offset'ів:**
- **Current offset** - поточна позиція споживача
- **Committed offset** - останній збережений offset
- **Latest offset** - останнє повідомлення в партиції

### 🖥️ Broker (Брокер)
**Брокер** - це сервер Kafka, який зберігає дані та обслуговує клієнтів.

**Функції брокера:**
- Зберігання партицій топіків
- Обробка запитів від продюсерів та споживачів
- Реплікація даних між брокерами
- Координація роботи кластера

### 🔄 Replication (Реплікація)
**Реплікація** - це копіювання даних між брокерами для забезпечення відмовостійкості.

```
Topic: orders (replication-factor: 3)
├── Broker 1: Partition 0 (Leader), Partition 1 (Follower)
├── Broker 2: Partition 0 (Follower), Partition 1 (Leader)  
└── Broker 3: Partition 0 (Follower), Partition 1 (Follower)
```

## Архітектура системи

### Наш кластер Kafka
```
┌─────────────────────────────────────────────────────────┐
│                    Kafka Cluster                        │
├─────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │   Broker 1  │  │   Broker 2  │  │   Broker 3  │     │
│  │ Port: 9091  │  │ Port: 9092  │  │ Port: 9093  │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
├─────────────────────────────────────────────────────────┤
│                   ZooKeeper                             │
│                  Port: 2181                             │
└─────────────────────────────────────────────────────────┘
```

### Компоненти нашого проекту
```
kafkagolang/
├── cmd/
│   ├── producer/     # Додаток для відправки повідомлень
│   └── consumer/     # Додаток для читання повідомлень
├── internal/
│   ├── kafka/        # Логіка роботи з Kafka
│   └── handler/      # Обробка повідомлень
├── compose.yml       # Docker конфігурація кластера
└── go.mod           # Go залежності
```

## Як запустити проект

### 1. Запуск Kafka кластера
```bash
# Запуск всіх сервісів (Kafka + ZooKeeper + UI)
docker-compose up -d

# Перевірка статусу
docker-compose ps
```

### 2. Збірка та запуск Producer
```bash
# Збірка
go build -o bin/producer cmd/producer/main.go

# Запуск (відправить 100 повідомлень)
./bin/producer
```

### 3. Запуск Consumer
```bash
# Збірка
go build -o bin/consumer cmd/consumer/main.go

# Запуск (3 споживача в одній групі)
./bin/consumer
```

### 4. Моніторинг через Kafka UI
Відкрийте http://localhost:8081 для перегляду:
- Топіків та партицій
- Повідомлень
- Груп споживачів
- Метрик кластера

## Детальний розбір коду

### Producer (cmd/producer/main.go)

```go
// Генеруємо унікальні UUID ключі для розподілу повідомлень по партиціях
keys := generateUUIDString()

for i := 0; i < 100; i++ {
    msg := fmt.Sprintf("Kafka message %d", i)
    
    // Використовуємо модуло для циклічного вибору ключа
    // Ключ визначає, в яку партицію потрапить повідомлення
    key := keys[i%numberOfKeys]
    
    // Відправляємо повідомлення
    if err := p.Produce(msg, topic, key, time.Now()); err != nil {
        logrus.Errorf("failed to produce message: %v", err)
    }
}
```

**Чому використовуємо ключі?**
- Повідомлення з однаковим ключем завжди потраплять в одну партицію
- Це гарантує порядок обробки для пов'язаних повідомлень
- Дозволяє рівномірно розподіляти навантаження

### Consumer (cmd/consumer/main.go)

```go
// Створюємо 3 consumer'и в одній групі
c1, err := kafka.NewConsumer(h, address, topic, consumerGroup, 1)
c2, err := kafka.NewConsumer(h, address, topic, consumerGroup, 2)
c3, err := kafka.NewConsumer(h, address, topic, consumerGroup, 3)

// Запускаємо в окремих горутинах
go c1.Start()
go c2.Start()
go c3.Start()
```

**Чому 3 споживача?**
- Кожен споживач обробляє свою частину партицій
- При падінні одного, інші продовжують роботу
- Можна масштабувати продуктивність

### Kafka Producer (internal/kafka/producer.go)

```go
func (p *Producer) Produce(message string, topic string, key string, t time.Time) error {
    kafkaMsg := &kafka.Message{
        TopicPartition: kafka.TopicPartition{
            Topic:     &topic,
            Partition: kafka.PartitionAny, // Автовибір партиції
        },
        Value:     []byte(message), // Тіло повідомлення
        Key:       []byte(key),     // Ключ для партиціонування
        Timestamp: t,               // Час створення
    }
    
    // Асинхронна відправка з очікуванням результату
    kafkaChan := make(chan kafka.Event)
    if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
        return fmt.Errorf("failed to produce message: %w", err)
    }
    
    // Обробка результату
    e := <-kafkaChan
    switch ev := e.(type) {
    case *kafka.Message:
        return nil // Успіх
    case kafka.Error:
        return fmt.Errorf("kafka error: %w", ev)
    }
}
```

### Kafka Consumer (internal/kafka/consumer.go)

```go
func (c *Consumer) Start() {
    for {
        if c.stop {
            break
        }

        // Блокуючий виклик - чекаємо нові повідомлення
        kafkaMsg, err := c.consumer.ReadMessage(noTimeout)
        if err != nil {
            logrus.Error(err)
        }

        if kafkaMsg == nil {
            continue
        }

        // Обробляємо повідомлення
        if err = c.handler.HandleMessage(kafkaMsg.Value, kafkaMsg.TopicPartition, c.consumerNumber); err != nil {
            logrus.Error(err)
            continue
        }

        // Зберігаємо offset для можливості відновлення
        if _, err = c.consumer.StoreMessage(kafkaMsg); err != nil {
            logrus.Error(err)
            continue
        }
    }
}
```

**Важливі налаштування Consumer:**
- `enable.auto.offset.store: false` - ручне керування offset'ами
- `enable.auto.commit: true` - автоматичний коміт кожні 5 секунд
- `auto.offset.reset: earliest` - починати з найраніших повідомлень
## Блок-схеми взаємодії

### 1. Архітектура Kafka кластера

```
                    ┌─────────────────────────────────────┐
                    │           ZooKeeper                 │
                    │      (Координація кластера)         │
                    └─────────────────┬───────────────────┘
                                      │
        ┌─────────────────────────────┼─────────────────────────────┐
        │                             │                             │
   ┌────▼────┐                   ┌────▼────┐                   ┌────▼────┐
   │Broker 1 │                   │Broker 2 │                   │Broker 3 │
   │:9091    │                   │:9092    │                   │:9093    │
   └─────────┘                   └─────────┘                   └─────────┘
        │                             │                             │
        └─────────────────────────────┼─────────────────────────────┘
                                      │
                              ┌───────▼───────┐
                              │  Kafka Topic  │
                              │   my-topic    │
                              └───────────────┘
```

### 2. Структура топіка з партиціями

```
Topic: my-topic
├── Partition 0: [msg1:0] [msg4:1] [msg7:2] [msg10:3] ...
├── Partition 1: [msg2:0] [msg5:1] [msg8:2] [msg11:3] ...
└── Partition 2: [msg3:0] [msg6:1] [msg9:2] [msg12:3] ...
                    ↑        ↑        ↑         ↑
                 offset:0  offset:1  offset:2  offset:3
```

**Розподіл повідомлень по партиціях:**
- Повідомлення з ключем "key-A" → завжди Partition 0
- Повідомлення з ключем "key-B" → завжди Partition 1  
- Повідомлення з ключем "key-C" → завжди Partition 2

### 3. Потік даних Producer → Kafka

```
┌─────────────┐    1. Produce     ┌─────────────┐
│  Producer   │ ─────────────────▶│   Broker    │
│             │                   │             │
│ - message   │    2. Ack/Error   │ - stores    │
│ - topic     │ ◀─────────────────│ - replicates│
│ - key       │                   │ - responds  │
│ - timestamp │                   │             │
└─────────────┘                   └─────────────┘

Кроки:
1. Producer відправляє повідомлення
2. Broker визначає партицію за ключем (hash(key) % partitions)
3. Broker зберігає повідомлення в лог партиції
4. Broker реплікує на інші брокери (якщо replication > 1)
5. Broker відправляє підтвердження Producer'у
```

### 4. Потік даних Kafka → Consumer

```
┌─────────────┐    1. Poll        ┌─────────────┐
│   Broker    │ ◀─────────────────│  Consumer   │
│             │                   │             │
│ - reads     │    2. Messages    │ - processes │
│ - tracks    │ ─────────────────▶│ - commits   │
│ - commits   │                   │ - stores    │
│             │    3. Commit      │   offset    │
│             │ ◀─────────────────│             │
└─────────────┘                   └─────────────┘

Кроки:
1. Consumer запитує нові повідомлення (poll)
2. Broker повертає повідомлення з поточного offset'у
3. Consumer обробляє повідомлення
4. Consumer зберігає offset локально
5. Consumer періодично комітить offset в Kafka
```

### 5. Consumer Group балансування

```
Consumer Group: my-consumer-group

Topic: my-topic (3 партиції)
├── Partition 0 ──────────────────▶ Consumer 1
├── Partition 1 ──────────────────▶ Consumer 2
└── Partition 2 ──────────────────▶ Consumer 3

При падінні Consumer 2:
├── Partition 0 ──────────────────▶ Consumer 1
├── Partition 1 ──────────────────▶ Consumer 3 (rebalance)
└── Partition 2 ──────────────────▶ Consumer 3

При додаванні Consumer 4:
├── Partition 0 ──────────────────▶ Consumer 1
├── Partition 1 ──────────────────▶ Consumer 2
├── Partition 2 ──────────────────▶ Consumer 3
└── (Consumer 4 чекає або отримає нові партиції)
```

### 6. Життєвий цикл повідомлення

```
┌─────────────┐
│  Producer   │
│  створює    │
│ повідомлення│
└──────┬──────┘
       │ 1. produce(msg, topic, key)
       ▼
┌─────────────┐
│   Kafka     │
│   Broker    │ ──┐ 2. зберігає в партицію
│             │   │    за hash(key)
└──────┬──────┘   │
       │          │ 3. реплікує на інші брокери
       │          │
       │ 4. ack   │
       ▼          │
┌─────────────┐   │
│  Producer   │   │
│ отримує     │   │
│підтвердження│   │
└─────────────┘   │
                  │
       ┌──────────┘
       │ 5. повідомлення доступне для читання
       ▼
┌─────────────┐
│  Consumer   │
│   читає     │ ──┐ 6. poll() запит
│повідомлення │   │
└──────┬──────┘   │
       │          │ 7. обробляє повідомлення
       │          │
       │ 8. commit offset
       ▼          │
┌─────────────┐   │
│   Kafka     │   │
│   зберігає  │ ◀─┘
│   offset    │
└─────────────┘
```

### 7. Реплікація та відмовостійкість

```
Topic: my-topic (replication-factor: 3)

Partition 0:
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Broker 1   │────▶│  Broker 2   │────▶│  Broker 3   │
│  (Leader)   │     │ (Follower)  │     │ (Follower)  │
│ [msg1][msg2]│     │ [msg1][msg2]│     │ [msg1][msg2]│
└─────────────┘     └─────────────┘     └─────────────┘

При падінні Leader (Broker 1):
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Broker 1   │     │  Broker 2   │────▶│  Broker 3   │
│   (DOWN)    │     │ (New Leader)│     │ (Follower)  │
│      X      │     │ [msg1][msg2]│     │ [msg1][msg2]│
└─────────────┘     └─────────────┘     └─────────────┘

Kafka автоматично:
1. Виявляє падіння Leader'а
2. Обирає нового Leader'а з Follower'ів
3. Перенаправляє клієнтів на нового Leader'а
4. Продовжує роботу без втрати даних
```

## Практичні сценарії використання

### 1. Event Sourcing
```go
// Зберігання всіх змін стану як події
events := []string{
    "UserCreated",
    "UserEmailChanged", 
    "UserDeactivated",
}

for _, event := range events {
    producer.Produce(event, "user-events", userID, time.Now())
}
```

### 2. Microservices Communication
```go
// Сервіс замовлень публікує подію
producer.Produce(`{
    "orderId": "12345",
    "status": "created",
    "amount": 100.50
}`, "order-events", orderID, time.Now())

// Сервіс платежів, складу, доставки підписуються на події
```

### 3. Real-time Analytics
```go
// Збір метрик в реальному часі
producer.Produce(`{
    "userId": "user123",
    "action": "page_view",
    "page": "/products",
    "timestamp": "2024-01-15T10:30:00Z"
}`, "user-analytics", userID, time.Now())
```

### 4. Log Aggregation
```go
// Централізований збір логів з різних сервісів
producer.Produce(`{
    "service": "auth-service",
    "level": "ERROR", 
    "message": "Failed login attempt",
    "timestamp": "2024-01-15T10:30:00Z"
}`, "application-logs", serviceID, time.Now())
```

## Моніторинг та налагодження

### Корисні команди для роботи з Kafka

```bash
# Створення топіка
docker exec kafka1 kafka-topics --create \
  --topic test-topic \
  --partitions 3 \
  --replication-factor 2 \
  --bootstrap-server localhost:9091

# Перегляд топіків
docker exec kafka1 kafka-topics --list \
  --bootstrap-server localhost:9091

# Деталі топіка
docker exec kafka1 kafka-topics --describe \
  --topic my-topic \
  --bootstrap-server localhost:9091

# Перегляд повідомлень
docker exec kafka1 kafka-console-consumer \
  --topic my-topic \
  --from-beginning \
  --bootstrap-server localhost:9091

# Перегляд груп споживачів
docker exec kafka1 kafka-consumer-groups --list \
  --bootstrap-server localhost:9091

# Деталі групи споживачів
docker exec kafka1 kafka-consumer-groups --describe \
  --group my-consumer-group \
  --bootstrap-server localhost:9091
```

### Метрики для моніторингу

**Producer метрики:**
- `record-send-rate` - швидкість відправки повідомлень
- `record-error-rate` - частота помилок
- `request-latency-avg` - середня затримка запитів

**Consumer метрики:**
- `records-consumed-rate` - швидкість споживання
- `records-lag-max` - максимальне відставання
- `commit-rate` - частота комітів offset'ів

**Broker метрики:**
- `bytes-in-per-sec` - вхідний трафік
- `bytes-out-per-sec` - вихідний трафік
- `under-replicated-partitions` - нерепліковані партиції

## Найкращі практики

### 1. Вибір кількості партицій
```
Рекомендації:
- Початкова кількість: 2-3 партиції на брокер
- Максимальна пропускність: 1 партиція ≈ 10MB/s
- Споживачі: не більше споживачів ніж партицій в групі
```

### 2. Налаштування retention
```yaml
# В compose.yml або конфігурації топіка
KAFKA_LOG_RETENTION_HOURS: 168  # 7 днів
KAFKA_LOG_RETENTION_BYTES: 1073741824  # 1GB
KAFKA_LOG_SEGMENT_BYTES: 104857600  # 100MB
```

### 3. Обробка помилок
```go
// Retry логіка для Producer
func (p *Producer) ProduceWithRetry(msg, topic, key string, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        if err := p.Produce(msg, topic, key, time.Now()); err != nil {
            if i == maxRetries-1 {
                return err
            }
            time.Sleep(time.Duration(i+1) * time.Second)
            continue
        }
        return nil
    }
    return nil
}

// Dead Letter Queue для Consumer
func (h *Handler) HandleMessage(message []byte, topic kafka.TopicPartition, consumerNumber int) error {
    if err := h.processMessage(message); err != nil {
        // Відправляємо в DLQ для подальшого аналізу
        return h.sendToDLQ(message, err)
    }
    return nil
}
```

### 4. Безпека
```go
// SSL/SASL конфігурація
cfg := &kafka.ConfigMap{
    "bootstrap.servers": "localhost:9093",
    "security.protocol": "SASL_SSL",
    "sasl.mechanism":    "PLAIN",
    "sasl.username":     "username",
    "sasl.password":     "password",
}
```

## Заключення

Apache Kafka - це потужна платформа для побудови розподілених систем реального часу. Основні переваги:

✅ **Масштабованість** - легко додавати брокери та споживачів  
✅ **Надійність** - реплікація та відмовостійкість  
✅ **Продуктивність** - мільйони повідомлень на секунду  
✅ **Гнучкість** - підтримка різних патернів використання  

Цей проект демонструє базові концепції роботи з Kafka на Go, включаючи:
- Налаштування кластера з 3 брокерів
- Створення Producer'а для відправки повідомлень
- Створення Consumer Group з 3 споживачами
- Обробку повідомлень та управління offset'ами
- Моніторинг через Kafka UI

Для production використання рекомендується додати:
- Метрики та алерти
- Логування та трейсинг  
- Обробку помилок та retry логіку
- Безпеку (SSL/SASL)
- Backup та disaster recovery

---

**Корисні посилання:**
- [Офіційна документація Kafka](https://kafka.apache.org/documentation/)
- [Confluent Kafka Go Client](https://github.com/confluentinc/confluent-kafka-go)
- [Kafka UI](https://github.com/provectus/kafka-ui)