package model

import (
	"fmt"
	"testing"
)

func TestGetInfoCoreDNSConf(t *testing.T) {
	name := "hogehoge.hoge"
	domainFileInfo := `# DomainUUID: 3e8fc6b1-0a93-4c57-9f2b-95d3e66a66e0
172.21.1.1  hogeserver1.hogehoge.hoge  # 5b9ea8eb-5ce5-422a-9d70-37d25fa896ae
172.21.1.2  hogeserver2.hogehoge.hoge  # f0c5edcd-3b18-4c26-a8e1-3f3495504dd6
`
	domain, err := NewDomain(name, domainFileInfo)
	if err != nil {
		t.Error(err)
	}

	var list []*Domain
	list = append(list, domain)
	conf := NewCoreDNSConf(list)
	confInfo, err := conf.GetFileInfo()
	if err != nil {
		t.Error(err)
	}

	expect := `hogehoge.hoge. {
    hosts hogehoge.hoge
    reload 10s 5s
    log
}

. {
    forward . 8.8.8.8
}
`

	fmt.Print(confInfo)
	if confInfo != expect {
		t.Error(confInfo)
	}
}

func TestAddCoreDNSConf(t *testing.T) {
	name := "hogehoge.hoge"
	domainFileInfo := `# DomainUUID: 3e8fc6b1-0a93-4c57-9f2b-95d3e66a66e0
172.21.1.1  hogeserver1.hogehoge.hoge  # 5b9ea8eb-5ce5-422a-9d70-37d25fa896ae
172.21.1.2  hogeserver2.hogehoge.hoge  # f0c5edcd-3b18-4c26-a8e1-3f3495504dd6
`
	domain, err := NewDomain(name, domainFileInfo)
	if err != nil {
		t.Error(err)
	}

	var list []*Domain
	list = append(list, domain)
	conf := NewCoreDNSConf(list)

	addDomName := "fugafuga.fuga"
	addDomain, _ := NewOriginalDomain(addDomName)

	conf.Add(addDomain)
	addedDomain, _ := conf.GetByName(addDomain.Name)
	if addedDomain != addDomain {
		t.Error(addedDomain)
	}

	addedDomain, _ = conf.GetByUuid(addDomain.Uuid)
	if addedDomain != addDomain {
		t.Error(addedDomain)
	}

	addedDomainList := conf.GetAll()
	if len(addedDomainList) > 1 && (addedDomainList[0] != addDomain && addedDomainList[1] != addDomain) {
		t.Error(addedDomainList[0])
	}
}
