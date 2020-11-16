package model

import (
	"coredns_api/internal/infrastructure"
	"errors"
	"strings"
)

type HostsFile struct {
	domain *Domain
	hosts  []*Host
}

func NewHostsFile(domainName string) (*HostsFile, error) {
	domainInfo, err := infrastructure.LoadHostsFile(domainName)
	if err != nil {
		return nil, err
	}

	var domain *Domain
	var hosts []*Host
	for _, line := range strings.Split(domainInfo, "\n") {
		// Expected hosts info is like this:
		//
		// ```
		// # DomainUUID: 3e8fc6b1-0a93-4c57-9f2b-95d3e66a66e0
		// 172.21.1.1  hogeserver1.hogehoge.hoge  # 5b9ea8eb-5ce5-422a-9d70-37d25fa896ae
		// 172.21.1.2  hogeserver2.hogehoge.hoge  # f0c5edcd-3b18-4c26-a8e1-3f3495504dd6
		// ````

		splitedLine := strings.Split(line, "#")
		hostInfo := splitedLine[0]
		commentInfo := splitedLine[len(splitedLine)-1]

		if strings.Contains(commentInfo, "DomainUUID:") {
			splitedComment := strings.Split(commentInfo, " ")
			domainID := splitedComment[len(splitedComment)-1]
			domain, err = NewDomain(domainID, domainName)
			if err != nil {
				return nil, err
			}
		}

		uuid := strings.Split(commentInfo, " ")[1]

		splitedHostInfo := strings.Split(hostInfo, " ")
		address := splitedHostInfo[0]
		hostname := splitedHostInfo[len(splitedHostInfo)-1]

		host, err := NewHost(uuid, hostname, address)
		if err != nil {
			return nil, err
		}

		hosts = append(hosts, host)
	}

	if domain == nil {
		return nil, errors.New("DomainUUID is not in hosts file.")
	}

	return &HostsFile{domain, hosts}, nil
}
