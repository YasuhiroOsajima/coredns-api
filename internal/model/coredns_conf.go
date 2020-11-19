package model

import (
	"bytes"
	"errors"
	"log"
	"os"
	"strings"
	"sync"
	"text/template"
)

// # cat coredns.conf
// hogehoge.hoge. {
//     hosts /var/lib/coredns/hosts/hogehoge.hoge
//
//     reload 10s 5s
//     log
// }
//
// fugafuga.fuga. {
//     hosts /var/lib/coredns/hosts/fugafuga.fuga
//
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
	Cache map[DomainName]*Domain

	confPath string
	forward  string
}

func NewCoreDNSConf(allDomainInfo []*Domain) *CoreDNSConf {
	confPath := os.Getenv("CONF_PATH")
	forward := `
. {
    forward . 8.8.8.8
}`

	cache := map[DomainName]*Domain{}
	for _, dom := range allDomainInfo {
		cache[dom.Name] = dom
	}
	return &CoreDNSConf{Cache: cache, confPath: confPath, forward: forward}
}

func (d *CoreDNSConf) Add(domain *Domain) {
	d.Cache[domain.Name] = domain
}

func (d *CoreDNSConf) Get(domainName DomainName) (*Domain, error) {
	domain := d.Cache[domainName]
	if domain == nil {
		return nil, errors.New("target domain chache is not found")
	}

	return domain, nil
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
			log.Fatal(dom)
			return "", err
		}
		domainInfoBottom := out.String()
		conf = conf + domainInfoTop + domainInfoBottom
	}

	conf = conf + d.forward
	return conf, nil
}
