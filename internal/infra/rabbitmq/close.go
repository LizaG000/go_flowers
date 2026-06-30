package rabbitmq

func (rcl *RabbitClient) Close() {
	if rcl.sendChan != nil {
		_ = rcl.sendChan.Close()
	}

	if rcl.recChan != nil {
		_ = rcl.recChan.Close()
	}

	if rcl.sendConn != nil {
		_ = rcl.sendConn.Close()
	}

	if rcl.recConn != nil {
		_ = rcl.recConn.Close()
	}
}
