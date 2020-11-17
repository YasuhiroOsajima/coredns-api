package infrastructure

import (
	"io/ioutil"
	"os"

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

func (f *Filesystem) WriteTextFile(fileName, fileInfo string) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte(fileInfo))
	if err != nil {
		return err
	}

	return nil
}

func (f *Filesystem) DeleteFile(fileName string) error {
	return os.Remove(fileName)
}
