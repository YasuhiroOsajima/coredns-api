package usecase

// error status with HTTP 500
type IsNotLockedError struct {
	err string
}

func NewIsNotLockedError() error {
	return &IsNotLockedError{err: "this function have to be inside of Mutex lock"}
}

func (e *IsNotLockedError) Error() string {
	return e.err
}

// error status with HTTP 400
type HostDuplicatedError struct {
	err string
}

func NewHostDuplicatedError(param, value string) error {
	return &HostDuplicatedError{err: "specified host parameter is already assigned in the domain. '" + param + ": " + value + "'"}
}

func (e *HostDuplicatedError) Error() string {
	return e.err
}
