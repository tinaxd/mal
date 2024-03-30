package main

type MalError struct {
	message MalValue
	value   MalValue
}

func NewError(message string) *MalError {
	return &MalError{message: NewString(message), value: NewString(message)}
}

func NewErrorFromError(err error) *MalError {
	return &MalError{message: NewString(err.Error()), value: NewString(err.Error())}
}

func NewErrorFromValue(message MalValue) *MalError {
	return &MalError{message: message, value: message}
}

func (e *MalError) Error() string {
	return PrStr(e.message, false)
}

func (e *MalError) Message() MalValue {
	return e.message
}

func (e *MalError) Value() MalValue {
	return e.value
}
