package model

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestNewDomain(t *testing.T) {
	name := "hogehoge.hoge"
	domainFileInfo := `# DomainUUID: 3e8fc6b1-0a93-4c57-9f2b-95d3e66a66e0
# Tenats:
#   - df397e50-8006-450e-b18b-5c5bd940baff
#   - 02c03bd4-fe2e-45f2-85b6-b535af15215d
172.21.1.1  hogeserver1.hogehoge.hoge  # 5b9ea8eb-5ce5-422a-9d70-37d25fa896ae
172.21.1.2  hogeserver2.hogehoge.hoge  # f0c5edcd-3b18-4c26-a8e1-3f3495504dd6
`
	domain, err := NewDomain(name, domainFileInfo)
	if err != nil {
		t.Error(err)
	}

	expect, err := NewDomainName(name)
	if err != nil {
		t.Error(err)
	}
	if domain.Name != expect {
		t.Error(errors.New("NewDomain is not successes"))
	}

	info, err := domain.GetFileInfo()
	if err != nil {
		t.Error(err)
	}

	fmt.Print(domainFileInfo)
	if info != domainFileInfo {
		t.Error("domainFileInfo is missmatched")
	}
}

func TestNewOriginalDomain(t *testing.T) {
	name := "hogehoge.hoge"
	tenant := []string{"5cdc62c5-a110-4d89-9cdd-5e19f1983f0f"}
	domain, err := NewOriginalDomain(name, tenant)
	if err != nil {
		t.Error(err)
	}

	expect, err := NewDomainName(name)
	if err != nil {
		t.Error(err)
	}
	if domain.Name != expect {
		t.Error(errors.New("NewDomain is not successes"))
	}
	info, err := domain.GetFileInfo()
	if err != nil {
		t.Error(err)
	}

	if !strings.HasPrefix(info, "# DomainUUID: ") {
		t.Error("domainFileInfo is missmatched")
	}
}
