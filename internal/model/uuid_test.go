package model

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewUuid(t *testing.T) {
	u, _ := uuid.NewRandom()
	newUuid, err := NewUuid(u.String())
	if err != nil {
		t.Error(err)
	}
	if newUuid.String() != u.String() {
		t.Error("domainName is missmatched")
	}
}
