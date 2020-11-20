package model

import (
	"bytes"
	"regexp"
	"strings"
	"text/template"

	"github.com/google/uuid"
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
	nameMatcher := regexp.MustCompile("^[0-9a-zA-Z._-]+$").MatchString
	if len(name) == 0 || !nameMatcher(name) {
		mes := "invalid Host name is specified with uuid: '" + uuid.String() + "', nam: '" + name + "', address: '" + address + "'"
		return nil, NewInvalidParameterGiven(mes)
	}

	ipMatcher := regexp.MustCompile("^[0-9.]+$").MatchString
	if len(address) < 7 || len(address) > 15 || strings.Count(address, ".") != 3 || !ipMatcher(address) {
		mes := "invalid IP address is specified with uuid: '" + uuid.String() + "', nam: '" + name + "', address: '" + address + "'"
		return nil, NewInvalidParameterGiven(mes)
	}

	return &Host{uuid, name, address}, nil
}

func (h *Host) GetHostInfo() (string, error) {
	hostInfo := `{{ .Address }}  {{ .Name }}  # {{ .Uuid }}
`
	tmpl := template.Must(template.New("").Parse(hostInfo))

	var out bytes.Buffer
	err := tmpl.Execute(&out, h)
	if err != nil {
		return "", err
	}
	result := out.String()

	return result, nil
}
