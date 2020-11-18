package model

import (
	"bytes"
	"errors"
	"github.com/google/uuid"
	"regexp"
	"strings"
	"text/template"
)

type Host struct {
	Uuid    Uuid
	Name    string
	Address string
}

func NewOriginalHost(name, address string) (*Host, error) {
	u, _ := uuid.NewRandom()
	hostUuid, err := NewUuid(u.String())
	if err != nil {
		return nil, err
	}

	host, err := NewHost(hostUuid, name, address)
	if err != nil {
		return nil, err
	}

	return host, nil
}

func NewHost(uuid Uuid, name, address string) (*Host, error) {
	nameMatcher := regexp.MustCompile("^[0-9a-zA-Z_-]+$").MatchString
	if len(name) == 0 || !nameMatcher(name) {
		return nil, errors.New("invalid Host name is specified")
	}

	ipMatcher := regexp.MustCompile("^[0-9.]+$").MatchString
	if len(address) < 7 || len(address) > 15 || strings.Count(address, ".") != 3 || !ipMatcher(address) {
		return nil, errors.New("invalid Host name is specified")
	}

	return &Host{uuid, name, address}, nil
}

func (h *Host) GetHostInfo() (string, error) {
	hostInfo := `{{ .Address }}  {{ .Name }}  # {{ .Uuid }}`
	tmpl := template.Must(template.New("").Parse(hostInfo))

	var out bytes.Buffer
	err := tmpl.Execute(&out, h)
	if err != nil {
		return "", err
	}
	result := out.String()

	return result, nil
}
