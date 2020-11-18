package model

import "errors"

type Uuid string

func NewUuid(uuid string) (Uuid, error) {
	if len(uuid) == 0 || len(uuid) >= 37 {
		return "", errors.New("invalid Domain UUID is specified")
	}

	domainUuid := (Uuid)(uuid)
	return domainUuid, nil
}

func (d Uuid) String() string {
	return string(d)
}
