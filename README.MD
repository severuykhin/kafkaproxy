Library for convinient access to Kafka via the rest interface

## Usage
```
func main() {

	k := kafkaproxy.New(
		[]string{"localhost:9092"},
		kafkaproxy.WithPort("5436"), // default 80
		kafkaproxy.WithBatchSize(5), // default 10
		kafkaproxy.WithBatchTimeout(time.Millisecond*20), // default 100ms
		kafkaproxy.WithWriteTimeout(time.Second), // default 1s
	)

	ctx, cancel := context.WithCancel(context.TODO())

	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-interruptCh
		cancel()
	}()

	if err := k.Run(ctx); err != nil {
		fmt.Println("kafkaproxy finished with error: ", err)
	}

	fmt.Println("process work finished")
}
```

## Methods

### Push messages
Produce messages for a topic, optionally specifying keys, partitions or headers for the messages. Currently only json is available as a message format.
```
POST /topics/:topic_name
```
```
{
    "records": [
        {
            "value": {
                "name": "Igor",
                "age": 32
            },
            "partition": 0,
            "headers": {
                "key": "value",
                "version": "1",
                "some": 1
            }
        },
        {
            "value": {
                "name": "Alex",
                "age": 30
            },
            "partition": 0
        }
    ]
}
```
