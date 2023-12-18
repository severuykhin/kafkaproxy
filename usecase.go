package kafkaproxy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/severuykhin/goerrors"
)

type usecase struct {
	logger logger
	writer *kafka.Writer
}

func (u *usecase) PushMessages(ctx context.Context, topicName string, records []Record) error {

	var messages []kafka.Message
	for _, record := range records {
		messageValue, err := json.Marshal(record.Value)
		if err != nil {
			return goerrors.NewBadRequestErr().WithMessage(err.Error())
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

	err := u.writer.WriteMessages(ctx, messages...)
	if err != nil {
		return goerrors.NewInternalErr().WithMessage(err.Error())
	}

	return nil
}
