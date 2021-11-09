package tt

// import (
// 	"encoding/json"
// 	"fmt"

// 	"github.com/confluentinc/confluent-kafka-go/kafka"
// )

// type LoginLog struct {
// 	Username   string
// 	User_id    string
// 	Ip_address string
// }

// func Producer(ll LoginLog) {
// 	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka:9092"})
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	value, err := json.Marshal(ll)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	defer p.Close()

// 	// Delivery report handler for produced messages
// 	go func() {
// 		for e := range p.Events() {
// 			switch ev := e.(type) {
// 			case *kafka.Message:
// 				if ev.TopicPartition.Error != nil {
// 					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
// 				}
// 			}
// 		}
// 	}()

// 	// Produce messages to topic (asynchronously)
// 	topic := "success_login"

// 	p.Produce(&kafka.Message{
// 		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
// 		Value:          []byte(value),
// 	}, nil)

// 	p.Flush(15 * 50)
// }
