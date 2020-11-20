package model

type InvalidParameterGiven struct {
	err string
}

func NewInvalidParameterGiven(text string) error {
	return &InvalidParameterGiven{err: text}
}

func (e *InvalidParameterGiven) Error() string {
	return e.err
}

type ServerSideError struct {
	err string
}

func NewServerSideError(text string) error {
	return &ServerSideError{err: text}
}

func (e *ServerSideError) Error() string {
	return e.err
}
