package kafkaproxy

import (
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/severuykhin/goerrors"
	"github.com/severuykhin/snfiber"
)

type controller struct {
	logger logger
	writer *kafka.Writer
}

func (h *controller) PushMessages(request *snfiber.Request) (interface{}, error) {

	topicName := request.Params("topic_name", "")

	if topicName == "" {
		return nil, goerrors.NewBadRequestErr().WithMessage("topic_name must not be empty")
	}

	var requestParams PushMessagesRequest
	request.BodyParser(&requestParams)

	var messages []kafka.Message
	for _, record := range requestParams.Records {
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
			Key:     record.Key,
			Headers: messageHeaders,
		}

		messages = append(messages, kafkaMessage)
	}

	err := h.writer.WriteMessages(request.Context(), messages...)
	if err != nil {
		return nil, goerrors.NewInternalErr().WithMessage(err.Error())
	}

	return "ok", nil
}
