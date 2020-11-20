package usecase

// error status with HTTP 500
type IsNotLocked struct {
	err string
}

func NewIsNotLocked() error {
	return &IsNotLocked{err: "this function have to be inside of Mutex lock"}
}

func (e *IsNotLocked) Error() string {
	return e.err
}
