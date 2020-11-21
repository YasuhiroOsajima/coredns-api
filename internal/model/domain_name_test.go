package model

import "testing"

func TestNewDomainName(t *testing.T) {
	name := "hogehoge.hoge"
	domainName, err := NewDomainName(name)
	if err != nil {
		t.Error(err)
	}
	if domainName.String() != name {
		t.Error("domainName is missmatched")
	}
}
