package tt

// import (
// 	"encoding/json"
// 	"fmt"

// 	"github.com/confluentinc/confluent-kafka-go/kafka"
// )

// type LoginSuccessData struct {
// 	Ip_address string `json:"Ip_address"`
// 	Username   string `json:"Username"`
// 	Userid     string `json:"User_id"`
// }

// func Consumer() {
// 	// terminate := false

// 	var value LoginSuccessData
// 	c, err := kafka.NewConsumer(&kafka.ConfigMap{
// 		"bootstrap.servers": "kafka:9092",
// 		"group.id":          "success_login",
// 		"auto.offset.reset": "latest",
// 	})
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer c.Close()
// 	err = c.SubscribeTopics([]string{"success_login"}, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	for {
// 		msg, err := c.ReadMessage(-1)
// 		if err == nil {
// 			err = json.Unmarshal(msg.Value, &value)
// 			if err != nil {
// 				fmt.Println(err)
// 			}

// 			fmt.Println(value)

// 		} else {
// 			// The client will automatically try to recover from all errors.
// 			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
// 		}
// 	}
// }
