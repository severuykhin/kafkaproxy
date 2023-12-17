package kafkaproxy

type Record struct {
	Value    interface{}            `json:"value"`
	Parition int                    `json:"parition"`
	Key      []byte                 `json:"key"`
	Headers  map[string]interface{} `json:"headers"`
}

type PushMessagesRequest struct {
	Records []Record
}
