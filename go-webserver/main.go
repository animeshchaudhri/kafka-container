package main

import (
    "context"
    "log"
    "github.com/gofiber/fiber/v2"
    "github.com/IBM/sarama"
    "github.com/go-redis/redis/v8"
)

var (
    kafkaProducer sarama.SyncProducer
    redisClient   *redis.Client
    ctx           = context.Background()
)

func main() {
    var err error

  
    kafkaProducer, err = sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
    if err != nil {
        log.Fatalf("Failed to start Sarama producer: %v", err)
    }
    defer kafkaProducer.Close()

  
    consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
    if err != nil {
        log.Fatalf("Failed to start Sarama consumer: %v", err)
    }
    defer consumer.Close()

    redisClient = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    
    app := fiber.New()

    app.Post("/send", handlePost)
    app.Get("/receive", handleGet)

    go consumeFromKafka(consumer)

    // Start the Fiber app
    log.Fatal(app.Listen(":8080"))
}

func handlePost(c *fiber.Ctx) error {
    data := c.Body()
    go produceToKafka(data)
    return c.SendStatus(fiber.StatusAccepted)
}

func produceToKafka(data []byte) {
    msg := &sarama.ProducerMessage{
        Topic: "my_topic",
        Value: sarama.ByteEncoder(data),
    }
    _, _, err := kafkaProducer.SendMessage(msg)
    if err != nil {
        log.Printf("Failed to send message to Kafka: %v", err)
    }
}

func handleGet(c *fiber.Ctx) error {
    data, err := redisClient.Get(ctx, "data_key").Result()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }
    go produceToKafka([]byte(data))
    return c.SendString(data)
}

func consumeFromKafka(consumer sarama.Consumer) {
    partitionConsumer, err := consumer.ConsumePartition("my_topic", 0, sarama.OffsetNewest)
    if err != nil {
        log.Fatalf("Failed to start Sarama partition consumer: %v", err)
    }
    defer partitionConsumer.Close()

    for message := range partitionConsumer.Messages() {
        go func(msg *sarama.ConsumerMessage) {
            err := redisClient.Set(ctx, "data_key", string(msg.Value), 0).Err()
            if err != nil {
                log.Printf("Failed to set value in Redis: %v", err)
            }
        }(message)
    }
}
