package rabbitmq

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (rcl *RabbitClient) channel(isRec, recreate bool) (*amqp.Channel, error) {
	if recreate {
		if isRec {
			rcl.recChan = nil
		} else {
			rcl.sendChan = nil
		}
	}

	// если нет connection — сбрасываем channel
	if isRec && rcl.recConn == nil {
		rcl.recChan = nil
	}
	if !isRec && rcl.sendConn == nil {
		rcl.sendChan = nil
	}

	// если канал уже есть — возвращаем
	if isRec && rcl.recChan != nil {
		return rcl.recChan, nil
	}
	if !isRec && rcl.sendChan != nil {
		return rcl.sendChan, nil
	}

	// пытаемся установить connection
	for {
		conn, err := rcl.connect(isRec, recreate)
		if err == nil && conn != nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// создаём channel
	var err error
	if isRec {
		rcl.recChan, err = rcl.recConn.Channel()
	} else {
		rcl.sendChan, err = rcl.sendConn.Channel()
	}

	if err != nil {
		log.Println("--- could not create channel ---")
		time.Sleep(1 * time.Second)
		return nil, err
	}

	if isRec {
		return rcl.recChan, nil
	}
	return rcl.sendChan, nil
}
