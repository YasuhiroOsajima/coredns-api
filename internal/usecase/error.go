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

// error status with HTTP 404
type RecordNotFoundError struct {
	err string
}

func NewRecordNotFoundError() error {
	return &RecordNotFoundError{err: "target domain is not found in CoreDNS"}
}

func (e *RecordNotFoundError) Error() string {
	return e.err
}

// error status with HTTP 404
type HostNotFoundError struct {
	err string
}

func NewHostNotFoundError() error {
	return &HostNotFoundError{err: "target host is not found in CoreDNS"}
}

func (e *HostNotFoundError) Error() string {
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
