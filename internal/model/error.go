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

type DomainNotFoundError struct {
	err string
}

func NewDomainNotFoundError() error {
	return &DomainNotFoundError{err: "target domain is not found in CoreDNS"}
}

func (e *DomainNotFoundError) Error() string {
	return e.err
}

type HostNotFoundError struct {
	err string
}

func NewHostNotFoundError() error {
	return &HostNotFoundError{err: "target host is not found in CoreDNS"}
}

func (e *HostNotFoundError) Error() string {
	return e.err
}

type DomainPermissionError struct {
	err string
}

func NewDomainPermissionError() error {
	return &DomainPermissionError{err: "specified tenant does not have permission"}
}

func (e *DomainPermissionError) Error() string {
	return e.err
}
