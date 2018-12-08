package grupttor

// DefaultError implements errors interface to consist error handling
type DefaultError struct {
	Message string
}

// implemented error interface
func (de *DefaultError) Error() string {
	return de.Message
}

// InterruptorWrongStateError wrap error of wrong internal interruptor state
type InterruptorWrongStateError struct {
	*DefaultError
}

// CreateInterruptorWrongStateError factory to create wrong state error, accepts error string
func CreateInterruptorWrongStateError(message string) InterruptorWrongStateError {
	return InterruptorWrongStateError{
		&DefaultError{
			Message: message,
		},
	}
}

type HandlerAlreadyDefinedError struct {
	*DefaultError
}

func CreateHandlerAlreadyDefinedError(message string) HandlerAlreadyDefinedError {
	return HandlerAlreadyDefinedError{
		&DefaultError{
			Message: message,
		},
	}
}
