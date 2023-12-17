package kafkaproxy

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/severuykhin/snfiber"
)

type kafkaproxy struct {
	Port              string
	Logger            logger
	KafkaWriterConfig *kafkaWriterConfig
}

func New(brokers []string, opts ...optFunc) *kafkaproxy {
	k := kafkaproxy{
		Port:   defaultPort,
		Logger: defaultLogger,
		KafkaWriterConfig: &kafkaWriterConfig{
			Brokers:      brokers,
			BatchSize:    defaultWriterBatchSize,
			BatchTimeout: defaultWriterBatchTimeout,
			WriteTimeout: defaultWriterWriteTimeout,
		},
	}

	for _, optFunc := range opts {
		optFunc(&k)
	}

	return &k
}

func (k *kafkaproxy) Run(ctx context.Context) error {
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      k.KafkaWriterConfig.Brokers,
		BatchTimeout: k.KafkaWriterConfig.BatchTimeout,
		BatchSize:    k.KafkaWriterConfig.BatchSize,
		WriteTimeout: k.KafkaWriterConfig.WriteTimeout,
		// MaxAttempts: 1, // @TODO
	})

	controller := controller{
		logger: k.Logger,
		writer: kafkaWriter,
	}

	router := snfiber.NewRouter()
	router.Post("/topics/:topic_name", controller.PushMessages)

	server := snfiber.NewServer(router, snfiber.WithLogger(k.Logger), snfiber.WithMetricsRoute())

	return server.Listen(":" + k.Port)
}
