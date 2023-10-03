package db

import (
	"fmt"
)

type SystemConfig struct {
	Kafka KafkaConfig `json:"kafka_config"`
}

type KafkaConfig struct {
	Url    string `json:"url"`
	Topic1 string `json:"topic_1"`
}

var C SystemConfig

func InitConfig(external bool) {
	// 这里需要注意，如果内网则使用docker-compose的hostnam和9092
	// 如果外网则需要docker-compose指定EXTERNAL的端口例如9094
	Url := "kafka:9092"
	if external {
		Url = "localhost:9094"
	}
	kafka_config := KafkaConfig{
		Url:    Url,
		Topic1: "mytest",
	}
	C.Kafka = kafka_config
	fmt.Printf("Init config done: %+v\n", C)
}
