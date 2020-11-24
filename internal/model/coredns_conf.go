package model

import (
	"bytes"
	"log"
	"os"
	"strings"
	"sync"
	"text/template"
)

// # cat coredns.conf
// hogehoge.hoge. {
//     hosts /var/lib/coredns/hosts/hogehoge.hoge
//     reload 10s 5s
//     log
// }
//
// fugafuga.fuga. {
//     hosts /var/lib/coredns/hosts/fugafuga.fuga
//     reload 10s 5s
//     log
// }
//
// . {
//     forward . 8.8.8.8
// }

var hostsDir = os.Getenv("HOSTS_DIR")

func GetHostsDir() string {
	return hostsDir
}

type CoreDNSConf struct {
	sync.Mutex
	locked int

	Cache map[DomainName]*Domain

	forward  string
	ConfPath string
}

func NewCoreDNSConf(allDomainInfo []*Domain) *CoreDNSConf {
	confPath := os.Getenv("CONF_PATH")
	forward := `
. {
    forward . 8.8.8.8
}
`

	cache := map[DomainName]*Domain{}
	for _, dom := range allDomainInfo {
		cache[dom.Name] = dom
	}
	return &CoreDNSConf{locked: 0, Cache: cache, forward: forward, ConfPath: confPath}
}

func (d *CoreDNSConf) Add(domain *Domain) {
	d.Cache[domain.Name] = domain
}

func (d *CoreDNSConf) GetByName(domainName DomainName) (*Domain, error) {
	domain := d.Cache[domainName]
	if domain == nil {
		return nil, NewInvalidParameterGiven("target domain chache is not found. domain: " + domainName.String())
	}

	return domain, nil
}

func (d *CoreDNSConf) GetByUuid(domainUuid Uuid, requestTenantUuid Uuid) (*Domain, error) {
	for _, domain := range d.Cache {
		if domain.Uuid == domainUuid {
			for _, t := range domain.Tenants {
				if t == requestTenantUuid {
					return domain, nil
				} else {
					return nil, NewDomainPermissionError()
				}
			}
		}
	}
	return nil, NewDomainNotFoundError()

}

func (d *CoreDNSConf) GetAll() []*Domain {
	var domains []*Domain
	for _, domain := range d.Cache {
		domains = append(domains, domain)
	}
	return domains
}

func (d *CoreDNSConf) GetTenantAll(requestTenantUuid Uuid) []*Domain {
	var domains []*Domain
	for _, domain := range d.Cache {
		for _, tenantUuid := range domain.Tenants {
			if requestTenantUuid == tenantUuid {
				domains = append(domains, domain)
				continue
			}
		}
	}
	return domains
}

func (d *CoreDNSConf) Delete(domain *Domain) {
	delete(d.Cache, domain.Name)
}

func (d *CoreDNSConf) GetFileInfo() (string, error) {
	conf := ""

	domainBottomTemplate := `    hosts {{ .DomainFilePath }}
    reload {{ .ReloadInterval }} {{ .ReloadJitter }}
    log
}
`
	tmpl := template.Must(template.New("").Parse(domainBottomTemplate))

	for domName, dom := range d.Cache {
		domainInfoTop := strings.TrimSpace(domName.String()) + `. {
`

		var out bytes.Buffer
		err := tmpl.Execute(&out, dom)
		if err != nil {
			log.Print(err)
			return "", err
		}
		domainInfoBottom := out.String()
		conf = conf + domainInfoTop + domainInfoBottom
	}

	conf = conf + d.forward
	return conf, nil
}

func (d *CoreDNSConf) IsLocked() bool {
	if d.locked == 0 {
		return false
	} else {
		return true
	}
}

func (d *CoreDNSConf) SetLocke() {
	d.Lock()
	d.locked += 1
}

func (d *CoreDNSConf) UnSetLocke() {
	d.Unlock()
	d.locked -= 1
}
