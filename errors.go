package grupttor

type DefaultError struct {
	Message string
}

func (de *DefaultError) Error() string {
	return de.Message
}

type InterruptorWrongStateError struct {
	*DefaultError
}

func CreateInterruptorWrongStateError(message string) InterruptorWrongStateError {
	return InterruptorWrongStateError{
		&DefaultError{
			Message: message,
		},
	}
}
