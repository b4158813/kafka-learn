package server

import (
	"context"
	"my_test/db"
	"os"
	"os/signal"
	"sync"

	"github.com/IBM/sarama"
)

type ExampleConsumerGroupHandler struct{}

func (ExampleConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ExampleConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h ExampleConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		db.L.Infof("Message key:%s value:%s topic:%q partition:%d offset:%d\n", msg.Key, msg.Value, msg.Topic, msg.Partition, msg.Offset)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func StartConsumerServer(ctx context.Context, groupId string) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup([]string{db.C.Kafka.Url}, groupId, config)
	if err != nil {
		db.L.Panicf("Error creating consumer group client: %v\n", err)
	}
	db.L.Infof("start consumer group %s done!", groupId)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side re-balance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(context.Background(), []string{db.C.Kafka.Topic1}, &ExampleConsumerGroupHandler{}); err != nil {
				db.L.Panicf("Error from consumer: %v\n", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
		}
	}(ctx)

	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	select {
	case <-signals:
		db.L.Infoln("Interrupt is detected")
	case <-ctx.Done():
		db.L.Errorln("closing consumer")
	}

	wg.Wait()
	if err = client.Close(); err != nil {
		db.L.Panicf("Error closing client: %v\n", err)
	}
}
