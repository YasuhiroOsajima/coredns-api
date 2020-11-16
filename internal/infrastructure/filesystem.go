package infrastructure

import (
	"io/ioutil"

	"coredns_api/internal/interface/repository"
)

type Filesystem struct{}

func NewFilesystem() repository.IFilesystem {
	return &Filesystem{}
}

func (f *Filesystem) LoadTextFile(fileName string) (string, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
