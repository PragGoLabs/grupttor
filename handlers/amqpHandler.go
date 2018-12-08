package handlers

import (
	"github.com/PragGoLabs/grupttor"
	"github.com/streadway/amqp"
)

type AmqpHandler struct {
	channel     *amqp.Channel
	consumerTag string
}

func NewAmqpHandler(channel *amqp.Channel, consumerTag string) AmqpHandler {
	return AmqpHandler{
		channel:     channel,
		consumerTag: consumerTag,
	}
}

// HandleInterrupt handler interrupt signal
func (ah AmqpHandler) HandleInterrupt(interruptor *grupttor.Grupttor) error {
	// close channel and wait until consuming end
	err := ah.channel.Cancel(ah.consumerTag, false)
	if err != nil {
		return CreateUnableToShutdownAmqpChannelError(err.Error())
	}

	return nil
}

// HandleStop handle stop signal
func (ah AmqpHandler) HandleStop(interruptor *grupttor.Grupttor) error {
	// do nothing
	return nil
}
