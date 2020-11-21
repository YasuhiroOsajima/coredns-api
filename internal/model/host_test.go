package model

import "testing"

func TestGetFQDN(t *testing.T) {
	targetHostname := "hogeserver1"
	targetDomain := "hogehoge.hoge"
	fqdn := GetFQDN(targetHostname, targetDomain)
	if fqdn != targetHostname+"."+targetDomain {
		t.Error("FQDN is missmatched")
	}
}
