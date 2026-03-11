package kafkaProduser

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	Producer     *kafka.Producer
	ProducerAddr string
	Topic        string
}

func NewKafkaProducer(addr string, topic string) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": addr})
	if err != nil {
		return nil, err
	}

	producer := &KafkaProducer{
		Producer:     p,
		ProducerAddr: addr,
		Topic:        topic,
	}

	go func() {
		for e := range producer.Producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Ошибка доставки: %v\n", ev.TopicPartition.Error)
				} else {
					fmt.Printf("Сообщение доставлено: %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	return producer, nil
}

func (p *KafkaProducer) Produce(msg string, reqId string, textId string) error {
	err := p.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.Topic, Partition: kafka.PartitionAny},
		Value:          []byte(msg),
		Headers: []kafka.Header{
			{
				Key:   "reqId",
				Value: []byte(reqId),
			},
			{
				Key:   "textId",
				Value: []byte(textId),
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	p.Producer.Flush(15 * 1000)
	return nil
}
