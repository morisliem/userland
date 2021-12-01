package kafka_job

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"userland/store/broker"
	"userland/store/postgres"

	"github.com/rs/zerolog/log"
)

func LoginLog(ctx context.Context, kafka broker.BrokerInterface, ls postgres.LogStoreInterface, endChan <-chan int) {
	c := kafka.GetConsumer()
	err := c.SubscribeTopics([]string{broker.TopicName}, nil)
	if err != nil {
		fmt.Println("error from helper")
	}

	for {
		select {
		case <-endChan:
			return
		default:
			msg, err := c.ReadMessage(-1)

			if err != nil {
				continue
			}

			job := broker.LoginLog{}
			err = json.Unmarshal(msg.Value, &job)
			if err != nil {
				continue
			}

			logData := postgres.LoginLog{
				Username:   job.Username,
				Ip_address: job.Ip_address,
				Created_at: time.Now(),
			}

			err = ls.SetLoginLog(logData)
			if err != nil {
				log.Error().Err(err).Msg(err.Error())
			}

		}
	}

}
