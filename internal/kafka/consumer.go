package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"wildberries-order-service/internal/config"
	"wildberries-order-service/internal/models"
)

// Consumer представляет Kafka consumer для заказов
type Consumer struct {
	config        *config.KafkaConfig
	orderService  models.OrderService
	consumerGroup sarama.ConsumerGroup
	ready         chan bool
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

// ConsumerGroupHandler реализует sarama.ConsumerGroupHandler
type ConsumerGroupHandler struct {
	consumer *Consumer
}

// NewConsumer создает новый Kafka consumer
func NewConsumer(cfg *config.KafkaConfig, orderService models.OrderService) (*Consumer, error) {
	if !cfg.Enabled {
		log.Println("📵 Kafka consumer disabled")
		return nil, nil
	}

	// Конфигурация Sarama
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Session.Timeout = 10 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	config.Consumer.MaxProcessingTime = 2 * time.Minute
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	// Создаем consumer group
	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	consumer := &Consumer{
		config:        cfg,
		orderService:  orderService,
		consumerGroup: consumerGroup,
		ready:         make(chan bool),
		ctx:           ctx,
		cancel:        cancel,
	}

	log.Printf("📨 Kafka consumer created for brokers: %v, topic: %s, group: %s", 
		cfg.Brokers, cfg.Topic, cfg.GroupID)

	return consumer, nil
}

// Start запускает consumer
func (c *Consumer) Start() error {
	if c == nil {
		log.Println("📵 Kafka consumer is disabled, skipping start")
		return nil
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-c.ctx.Done():
				log.Println("🛑 Kafka consumer context cancelled")
				return
			default:
				handler := &ConsumerGroupHandler{consumer: c}
				
				// Consume запускается в бесконечном цикле, пока контекст не отменен
				err := c.consumerGroup.Consume(c.ctx, []string{c.config.Topic}, handler)
				if err != nil {
					log.Printf("❌ Error from consumer: %v", err)
					// Небольшая пауза перед повторной попыткой
					select {
					case <-c.ctx.Done():
						return
					case <-time.After(5 * time.Second):
						continue
					}
				}
			}
		}
	}()

	// Ожидаем готовности consumer
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-c.ctx.Done():
				return
			case err := <-c.consumerGroup.Errors():
				if err != nil {
					log.Printf("❌ Consumer group error: %v", err)
				}
			case <-c.ready:
				log.Println("✅ Kafka consumer is ready and consuming messages")
				return
			}
		}
	}()

	return nil
}

// Stop останавливает consumer
func (c *Consumer) Stop() {
	if c == nil {
		return
	}

	log.Println("🛑 Stopping Kafka consumer...")
	
	c.cancel()
	c.wg.Wait()
	
	if err := c.consumerGroup.Close(); err != nil {
		log.Printf("❌ Error closing consumer group: %v", err)
	}
	
	log.Println("✅ Kafka consumer stopped")
}

// Setup реализует sarama.ConsumerGroupHandler
func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	close(h.consumer.ready)
	return nil
}

// Cleanup реализует sarama.ConsumerGroupHandler
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim реализует sarama.ConsumerGroupHandler
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Обрабатываем сообщения
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}
			
			if err := h.processMessage(message); err != nil {
				log.Printf("❌ Error processing message: %v", err)
				// Не помечаем сообщение как обработанное при ошибке
				continue
			}
			
			// Помечаем сообщение как обработанное
			session.MarkMessage(message, "")
			
		case <-session.Context().Done():
			return nil
		}
	}
}

// processMessage обрабатывает одно сообщение
func (h *ConsumerGroupHandler) processMessage(message *sarama.ConsumerMessage) error {
	log.Printf("📨 Received message: partition=%d, offset=%d, key=%s", 
		message.Partition, message.Offset, string(message.Key))

	// Парсим JSON сообщение
	var order models.Order
	if err := json.Unmarshal(message.Value, &order); err != nil {
		log.Printf("❌ Failed to unmarshal message: %v", err)
		return err
	}

	// Валидируем заказ
	if err := h.consumer.orderService.ValidateOrder(&order); err != nil {
		log.Printf("❌ Invalid order data: %v", err)
		// Возвращаем nil, чтобы сообщение было помечено как обработанное
		// (невалидные сообщения не должны блокировать очередь)
		return nil
	}

	// Сохраняем заказ через сервис
	if err := h.consumer.orderService.CreateOrder(&order); err != nil {
		log.Printf("❌ Failed to create order: %v", err)
		return err
	}

	log.Printf("✅ Successfully processed order: %s", order.OrderUID)
	return nil
}

// IsReady возвращает true, если consumer готов к работе
func (c *Consumer) IsReady() bool {
	if c == nil {
		return false
	}
	
	select {
	case <-c.ready:
		return true
	default:
		return false
	}
}