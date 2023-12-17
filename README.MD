Library for convinient access to Kafka via the rest interface

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

@TODO

- transfer the bussiness logic for generating and sending messages to usecase