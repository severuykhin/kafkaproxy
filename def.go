package kafkaproxy

type Record struct {
	Value    any            `json:"value"`
	Parition int            `json:"parition"`
	Key      string         `json:"key"`
	Headers  map[string]any `json:"headers"`
}

type PushMessagesRequest struct {
	Records []Record
}
