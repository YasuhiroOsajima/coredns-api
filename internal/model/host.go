package model

import (
	"errors"
	"regexp"
	"strings"
)

type Host struct {
	Uuid    string `json:"uuid"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

func NewHost(uuid, name, address string) (*Host, error) {

	if len(uuid) == 0 || len(uuid) >= 37 {
		return nil, errors.New("Invalid Host UUID is specified")
	}

	nameMatcher := regexp.MustCompile("^[0-9a-zA-Z_-]+$").MatchString
	if len(name) == 0 || !nameMatcher(name) {
		return nil, errors.New("Invalid Host name is specified")
	}

	ipMatcher := regexp.MustCompile("^[0-9.]+$").MatchString
	if len(address) < 7 || len(address) > 15 || strings.Count(address, ".") != 3 || !ipMatcher(address) {
		return nil, errors.New("Invalid Host name is specified")
	}

	return &Host{uuid, name, address}, nil
}
