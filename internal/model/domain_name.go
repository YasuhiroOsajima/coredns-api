package model

import (
	"errors"
	"regexp"
)

type DomainName string

func NewDomainName(name string) (DomainName, error) {
	nameMatcher := regexp.MustCompile("^[0-9a-zA-Z_-]+$").MatchString
	if len(name) == 0 || !nameMatcher(name) {
		return "", errors.New("invalid Domain name is specified")
	}

	return (DomainName)(name), nil
}

func (d DomainName) String() string {
	return string(d)
}
