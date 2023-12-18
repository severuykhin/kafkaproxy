package kafkaproxy

import (
	"github.com/severuykhin/goerrors"
	"github.com/severuykhin/snfiber"
)

type controller struct {
	logger  logger
	usecase usecase
}

func (h *controller) PushMessages(request *snfiber.Request) (interface{}, error) {

	topicName := request.Params("topic_name", "")

	if topicName == "" {
		return nil, goerrors.NewBadRequestErr().WithMessage("topic_name must not be empty")
	}

	var requestParams PushMessagesRequest
	err := request.BodyParser(&requestParams)
	if err != nil {
		return nil, err
	}

	if err := h.usecase.PushMessages(request.Context(), topicName, requestParams.Records); err != nil {
		return nil, err
	}

	return "ok", nil
}
