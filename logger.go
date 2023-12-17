package kafkaproxy

type logger interface {
	Error(message any, kayvals ...any)
}
