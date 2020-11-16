package model

import (
	"errors"
	"regexp"
)

type Domain struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

func NewDomain(uuid, name string) (*Domain, error) {

	if len(uuid) == 0 || len(uuid) >= 37 {
		return nil, errors.New("Invalid Domain UUID is specified")
	}

	nameMatcher := regexp.MustCompile("^[0-9a-zA-Z_-]+$").MatchString
	if len(name) == 0 || !nameMatcher(name) {
		return nil, errors.New("Invalid Domain name is specified")
	}

	return &Domain{uuid, name}, nil
}
