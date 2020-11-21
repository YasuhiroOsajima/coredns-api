package model

import "regexp"

type DomainName string

func NewDomainName(name string) (DomainName, error) {
	nameMatcher := regexp.MustCompile("^[0-9a-zA-Z.-]+$").MatchString
	if len(name) == 0 || !nameMatcher(name) {
		return "", NewInvalidParameterGiven("invalid Domain name is specified. name: " + name)
	}

	return (DomainName)(name), nil
}

func (d DomainName) String() string {
	return string(d)
}
