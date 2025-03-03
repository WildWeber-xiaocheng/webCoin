package svc

import (
	"jobcenter/internal/config"
	"jobcenter/internal/database"
)

type ServiceContext struct {
	Config      config.Config
	MongoClient *database.MongoClient
	KafkaClient *database.KafkaClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	kafkaClient := database.NewKafkaClient(c.Kafka)
	kafkaClient.StartWrite()
	return &ServiceContext{
		Config:      c,
		MongoClient: database.ConnectMongo(c.Mongo),
		KafkaClient: kafkaClient,
	}
}
