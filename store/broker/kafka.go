package broker

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const TopicName = "login_succeed"

type LoginLog struct {
	Username   string `json:"username"`
	Ip_address string `json:"ip_address"`
}

type Broker struct {
	producer *kafka.Producer
	consumer *kafka.Consumer
}

type BrokerInterface interface {
	SendLog(topic string, logData LoginLog) error
	GetConsumer() *kafka.Consumer
}

func (b Broker) SendLog(topic string, ll LoginLog) error {
	msg, err := json.Marshal(ll)
	if err != nil {
		return err
	}

	return b.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(msg),
	}, nil)
}

func (b Broker) GetConsumer() *kafka.Consumer {
	return b.consumer
}

func NewBroker(consumerConfig *kafka.ConfigMap, producerConfig *kafka.ConfigMap) (BrokerInterface, error) {
	c, err := kafka.NewConsumer(consumerConfig)
	if err != nil {
		return nil, err
	}

	p, err := kafka.NewProducer(producerConfig)
	if err != nil {
		return nil, err
	}

	broker := Broker{
		producer: p,
		consumer: c,
	}

	return broker, nil
}
