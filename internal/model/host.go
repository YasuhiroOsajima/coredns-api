package model

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
	"text/template"
)

type Host struct {
	Uuid    string
	Name    string
	Address string
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
