package handlers

// DefaultError implements errors interface to consist error handling
type DefaultError struct {
	Message string
}

// implemented error interface
func (de *DefaultError) Error() string {
	return de.Message
}

// UnableToShutdownHTTPServer wrap the http shutdown error
type UnableToShutdownHTTPServer struct {
	*DefaultError
}

// CreateUnableToShutdownHTTPServer factory to create http shutdown error
func CreateUnableToShutdownHTTPServer(message string) UnableToShutdownHTTPServer {
	return UnableToShutdownHTTPServer{
		&DefaultError{
			Message: message,
		},
	}
}

// UnableToShutdownHTTPServer wrap the http shutdown error
type UnableToShutdownAmqpChannelError  struct {
	*DefaultError
}

// CreateUnableToShutdownAmqpChannelError factory to create amqp channel cancel error
func CreateUnableToShutdownAmqpChannelError(message string) UnableToShutdownAmqpChannelError  {
	return UnableToShutdownAmqpChannelError {
		&DefaultError{
			Message: message,
		},
	}
}
