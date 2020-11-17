package model

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
	"text/template"
)

type Domain struct {
	Uuid  string
	Name  string
	Hosts []*Host
}

func NewDefaultDomain(uuid, name string) (*Domain, error) {

	if len(uuid) == 0 || len(uuid) >= 37 {
		return nil, errors.New("invalid Domain UUID is specified")
	}

	nameMatcher := regexp.MustCompile("^[0-9a-zA-Z_-]+$").MatchString
	if len(name) == 0 || !nameMatcher(name) {
		return nil, errors.New("invalid Domain name is specified")
	}

	var hosts []*Host

	return &Domain{uuid, name, hosts}, nil
}

func NewDomain(name, fileInfo string) (*Domain, error) {
	var domain *Domain
	var hosts []*Host
	var err error

	for _, line := range strings.Split(fileInfo, "\n") {
		// Expected hosts info is like this:
		//
		// ```
		// # DomainUUID: 3e8fc6b1-0a93-4c57-9f2b-95d3e66a66e0
		// 172.21.1.1  hogeserver1.hogehoge.hoge  # 5b9ea8eb-5ce5-422a-9d70-37d25fa896ae
		// 172.21.1.2  hogeserver2.hogehoge.hoge  # f0c5edcd-3b18-4c26-a8e1-3f3495504dd6
		// ````

		splitLine := strings.Split(line, "#")
		hostInfo := splitLine[0]
		commentInfo := splitLine[len(splitLine)-1]

		splitComment := strings.Fields(commentInfo)

		if strings.Contains(commentInfo, "DomainUUID:") {
			domainID := splitComment[len(splitComment)-1]
			domain, err = NewDefaultDomain(domainID, name)
			if err != nil {
				return nil, err
			}
		} else {
			hostId := splitComment[1]
			splitHost := strings.Fields(hostInfo)
			address := splitHost[0]
			hostName := splitHost[1]

			host, err := NewHost(hostName, address, hostId)
			if err != nil {
				return nil, err
			}

			hosts = append(hosts, host)
		}
	}

	if domain == nil {
		return nil, errors.New("domainUUID is not in hosts file")
	}

	domain.Hosts = hosts
	return domain, nil
}

func (d *Domain) GetFileInfo() (string, error) {
	fileInfo := `# DomainUUID: {{ .Uuid }}`
	tmpl := template.Must(template.New("").Parse(fileInfo))

	var out bytes.Buffer
	err := tmpl.Execute(&out, d)
	if err != nil {
		return "", err
	}
	result := out.String()

	for _, h := range d.Hosts {
		i, err := h.GetHostInfo()
		if err != nil {
			return "", err
		}
		result = result + "¥n" + i
	}

	return result, nil
}
