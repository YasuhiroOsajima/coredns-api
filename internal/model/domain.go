package model

import (
	"bytes"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/google/uuid"
)

func GetHostsFilePath(domainName DomainName) string {
	return filepath.Join(GetHostsDir(), domainName.String())
}

type Domain struct {
	Uuid           Uuid
	Name           DomainName
	Tenants        []Uuid
	Hosts          []*Host
	DomainFilePath string
	ReloadInterval string
	ReloadJitter   string
}

func NewOriginalDomain(name string, tenantList []string) (*Domain, error) {
	u, _ := uuid.NewRandom()
	domainUuid, err := NewUuid(u.String())
	if err != nil {
		return nil, err
	}

	var tenantUuidList []Uuid
	for _, t := range tenantList {
		tenantUuid, err := NewUuid(t)
		if err != nil {
			return nil, err
		}
		tenantUuidList = append(tenantUuidList, tenantUuid)
	}

	domain, err := NewEmptyDomain(domainUuid, name)
	if err != nil {
		return nil, err
	}
	domain.Tenants = tenantUuidList
	return domain, nil
}

func NewEmptyDomain(uuid Uuid, name string) (*Domain, error) {
	domainName, err := NewDomainName(name)
	if err != nil {
		return nil, err
	}

	hostsPath := GetHostsFilePath(domainName)
	var hosts []*Host
	var tenants []Uuid

	domain := &Domain{
		Uuid:           uuid,
		Name:           domainName,
		Tenants:        tenants,
		Hosts:          hosts,
		DomainFilePath: hostsPath,
		ReloadInterval: "10s",
		ReloadJitter:   "5s"}

	return domain, nil
}

func NewDomain(name, fileInfo string) (*Domain, error) {
	var domain *Domain
	var hosts []*Host
	var tenants []Uuid
	inTenats := false
	var err error

	for _, line := range strings.Split(fileInfo, "\n") {
		// Expected hosts info is like this:
		//
		// ```
		// # DomainUUID: 3e8fc6b1-0a93-4c57-9f2b-95d3e66a66e0
		// # Tenats:
		// #   - df397e50-8006-450e-b18b-5c5bd940baff
		// #   - 02c03bd4-fe2e-45f2-85b6-b535af15215d
		// 172.21.1.1  hogeserver1.hogehoge.hoge  # 5b9ea8eb-5ce5-422a-9d70-37d25fa896ae
		// 172.21.1.2  hogeserver2.hogehoge.hoge  # f0c5edcd-3b18-4c26-a8e1-3f3495504dd6
		// ````

		splitLine := strings.Split(line, "#")
		hostInfo := splitLine[0]
		commentInfo := splitLine[len(splitLine)-1]

		splitComment := strings.Fields(commentInfo)

		if strings.Contains(commentInfo, "DomainUUID:") {
			domainID := splitComment[len(splitComment)-1]
			dID := Uuid(domainID)
			domain, err = NewEmptyDomain(dID, name)
			if err != nil {
				return nil, err
			}
		} else if strings.Contains(commentInfo, "Tenats:") {
			inTenats = true
		} else if inTenats && strings.Contains(commentInfo, " - ") && strings.HasPrefix(line, "#") {
			tenantId := splitComment[len(splitComment)-1]
			tenantUuid, err := NewUuid(tenantId)
			if err != nil {
				return nil, NewServerSideError(err.Error())
			}
			tenants = append(tenants, tenantUuid)
		} else if strings.Contains(line, "-") && strings.Contains(line, ".") && strings.Contains(line, "#") {
			inTenats = false

			hostId := splitComment[0]
			splitHost := strings.Fields(hostInfo)
			address := splitHost[0]
			hostName := splitHost[1]

			hUuid, err := NewUuid(hostId)
			if err != nil {
				return nil, err
			}

			host, err := NewHost(hUuid, hostName, address)
			if err != nil {
				return nil, err
			}

			hosts = append(hosts, host)
		}
	}

	if domain == nil {
		return nil, NewServerSideError("domainUUID is not in hosts file info for " + name)
	}

	domain.Hosts = hosts
	domain.Tenants = tenants
	return domain, nil
}

func (d *Domain) GetFileInfo() (string, error) {
	fileInfo := `# DomainUUID: {{ .Uuid }}
`
	tmpl := template.Must(template.New("").Parse(fileInfo))

	var out bytes.Buffer
	err := tmpl.Execute(&out, d)
	if err != nil {
		return "", err
	}
	result := out.String()

	result += `# Tenats:
`

	for _, tUuid := range d.Tenants {
		result += `#   - ` + tUuid.String() + `
`
	}

	for _, h := range d.Hosts {
		i, err := h.GetHostInfo()
		if err != nil {
			return "", err
		}
		result = result + i
	}

	return result, nil
}

func (d *Domain) UpdateTenants(requestTenantUuid Uuid, tenantUuidList []Uuid) error {
	accessible := false
	for _, t := range d.Tenants {
		if t == requestTenantUuid {
			accessible = true
		}
	}
	if !accessible {
		return NewDomainPermissionError()
	}

	d.Tenants = tenantUuidList
	return nil
}
