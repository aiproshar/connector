package publisher

import (
	"backEnd/config"
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
	kafkaGo "github.com/segmentio/kafka-go"
	"sync"
	"time"
)

type KafkaManagerInterface interface {
	InitProducer() error
	CloseProducer()
	Publish(topic string, message string) error
	Consume(topic string, consumerGroup string) ([]string, error)
	CommitOffset(topic string, consumerGroup string, partition int, offset int, metadata string) error
}

type kafkaManager struct {
	Producer *kafka.Producer
}

// Implementing Singleton
var singletonKafkaManager *kafkaManager
var onceKmManager sync.Once

func GetKafkaManager() *kafkaManager {
	onceKmManager.Do(func() {
		log.Debug().Msg("Initializing Kafka Manager")
		singletonKafkaManager = &kafkaManager{}
		err := singletonKafkaManager.InitProducer()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to init Kafka Producer")
		} else {
			log.Debug().Msg("Kafka Producer Initialized")
		}
	})
	return singletonKafkaManager
}

func (km *kafkaManager) InitProducer() error {
	var p *kafka.Producer = nil
	var err error
	if config.SaslSsl == true {
		configMap := kafka.ConfigMap{
			"bootstrap.servers":                   config.KafkaBroker,
			"security.protocol":                   config.SASL_SSL,
			"enable.ssl.certificate.verification": config.TRUE,
			"sasl.mechanism":                      config.SASL_MECHANISM_AOUTHBEARER,
			"sasl.oauthbearer.method":             config.SASL_OAUTHBEARER_METHOD_OIDC,
			"sasl.oauthbearer.client.id":          config.SaslOauthbearerClientId,
			"sasl.oauthbearer.client.secret":      config.SaslOauthbearerClientSecret,
			"sasl.oauthbearer.token.endpoint.url": config.SaslOauthbearerTokenEndpointUri,
			"delivery.timeout.ms":                 2000,
		}
		log.Trace().
			Bool(config.CONFIDENTIAL, true).
			Interface("kafkaConfigmap", configMap).Msg("")

		p, err = kafka.NewProducer(&configMap)

		if err != nil {
			log.Error().Err(err).Msg("Kafka Producer construct error")
			return err
		}

		go func(eventsChan chan kafka.Event) {
			for ev := range eventsChan {
				oart, ok := ev.(kafka.OAuthBearerTokenRefresh)
				if !ok {
					log.Error().Err(err).
						Str("oart token", oart.String()).
						Msg("refresh token event error")
					continue
				}
				handleOAuthBearerTokenRefreshEvent(p, oart)
			}
		}(p.Events())

	} else {
		configMap := kafka.ConfigMap{
			"bootstrap.servers":   config.KafkaBroker,
			"delivery.timeout.ms": 2000,
		}
		p, err = kafka.NewProducer(&configMap)

		log.Trace().
			Bool(config.CONFIDENTIAL, true).
			Interface("kafkaConfigmap", configMap).Msg("")
		if err != nil {
			log.Error().Err(err).Msg("Kafka Producer construct error")
			return err
		}
	}

	km.Producer = p

	return nil
}

func (km *kafkaManager) CloseProducer() {
	log.Debug().Msg("closing Kafka Producer")
	if km.Producer != nil {
		km.Producer.Close()
	}
	log.Debug().Msg("Kafka Producer is closed")
}

func (km *kafkaManager) Publish(topic string, payload []byte) error {
	deliveryChan := make(chan kafka.Event)
	err := km.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          payload,
	}, deliveryChan)

	if err != nil {
		log.Error().Err(err).Msg("kafka publish error")
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		log.Error().Err(m.TopicPartition.Error).Msg("kafka delivery failed")
	} else {
		if log.Trace().Enabled() {
			log.Debug().
				Str("topic", *m.TopicPartition.Topic).
				Int32("partition", m.TopicPartition.Partition).
				Int("offset", int(m.TopicPartition.Offset)).
				Msg("kafka publish successful")
		} else {
			log.Debug().Msg("kafka publish successful")
		}

	}
	close(deliveryChan)

	return m.TopicPartition.Error
}

func (km *kafkaManager) Consume(topic string, consumerGroup string) ([]string, error) {

	var messages []string

	r := kafkaGo.NewReader(kafkaGo.ReaderConfig{
		Brokers:  []string{config.KafkaBroker},
		GroupID:  consumerGroup,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	d := time.Now().Add(10 * time.Second)
	ctx, _ := context.WithDeadline(context.Background(), d)

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			break
		}
		messages = append(messages, string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Error().Err(err).Msg("failed to close reader")

	}

	return messages, nil

}

func (km *kafkaManager) CommitOffset(topic string, consumerGroup string, partition int, offset int, metadata string) error {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.KafkaBroker,
		"group.id":          consumerGroup,
	})

	res, err := c.CommitOffsets([]kafka.TopicPartition{{
		Topic:     &topic,
		Partition: int32(partition),
		Metadata:  &metadata,
		Offset:    kafka.Offset(offset),
	}})
	if err != nil {
		log.Error().Err(err).Msg("Failed to commit offset")
		return err
	}

	log.Debug().
		Str("topic", *res[0].Topic).
		Int64("offset", int64(res[0].Offset)).
		Msg("Offset committed successfully")
	return nil
}
