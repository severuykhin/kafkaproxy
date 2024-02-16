package kafkaproxy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/severuykhin/goerrors"
	"github.com/severuykhin/snfiber"
)

type kafkaproxy struct {
	port   string
	logger logger
	writer *kafka.Writer
}

func New(brokers []string, opts ...optFunc) *kafkaproxy {

	config := kafkaProxyConfig{
		Port:         defaultPort,
		Logger:       defaultLogger,
		BatchTimeout: defaultWriterBatchTimeout,
		BatchSize:    defaultWriterBatchSize,
		WriteTimeout: defaultWriterWriteTimeout,
	}

	for _, optFunc := range opts {
		optFunc(&config)
	}

	k := kafkaproxy{
		port:   config.Port,
		logger: config.Logger,
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      brokers,
			BatchTimeout: config.BatchTimeout,
			BatchSize:    config.BatchSize,
			WriteTimeout: config.WriteTimeout,
			// MaxAttempts: 1, // @TODO
		}),
	}

	return &k
}

func (k *kafkaproxy) Run(ctx context.Context) error {

	router := snfiber.NewRouter()
	router.Post("/topics/:topic_name", k.PushMessages)

	server := snfiber.NewServer(router, snfiber.WithLogger(k.logger), snfiber.WithMetricsRoute())
	defer server.Shutdown()

	errCh := make(chan error)
	go func() {
		errCh <- server.Listen(":" + k.port)
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}

func (k *kafkaproxy) PushMessages(request *snfiber.Request) (interface{}, error) {

	topicName := request.Params("topic_name", "")

	if topicName == "" {
		return nil, goerrors.NewBadRequestErr().WithMessage("topic_name must not be empty")
	}

	var requestData PushMessagesRequest
	err := request.BodyParser(&requestData)
	if err != nil {
		return nil, err
	}

	var messages []kafka.Message
	for _, record := range requestData.Records {
		messageValue, err := json.Marshal(record.Value)
		if err != nil {
			return nil, goerrors.NewBadRequestErr().WithMessage(err.Error())
		}

		var messageHeaders []kafka.Header

		for headerKey, value := range record.Headers {
			headerValue := fmt.Sprintf("%v", value)
			messageHeaders = append(messageHeaders, kafka.Header{
				Key:   headerKey,
				Value: []byte(headerValue),
			})
		}

		kafkaMessage := kafka.Message{
			Topic:   topicName,
			Value:   messageValue,
			Key:     []byte(record.Key),
			Headers: messageHeaders,
		}

		messages = append(messages, kafkaMessage)
	}

	err = k.writer.WriteMessages(request.Context(), messages...)
	if err != nil {
		if err == kafka.UnknownTopicOrPartition {
			return nil, goerrors.NewNotFoundErr().WithMessage(err.Error()).WithContext("topic", topicName)
		}
		return nil, goerrors.NewInternalErr().WithMessage(err.Error()).WithContext("topic", topicName)
	}

	return nil, nil
}
