rabbitmq:
- ack and nack properly (when we have an error we need to return the message back to the queue)

documentation
tests
when we get a cancel we need to defer it. check where this is relevant

organize the utils
- functions in core.url should be moved to the encoding_decoding