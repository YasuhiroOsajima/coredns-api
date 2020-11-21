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

func (f *Filesystem) LoadTextFile(filePath string) (string, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (f *Filesystem) WriteTextFile(filePath, fileInfo string) error {
	file, err := os.Create(filePath)
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

func (f *Filesystem) DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

func (f *Filesystem) GetFilenameList(directory string) ([]string, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var fileNameList []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileNameList = append(fileNameList, file.Name())
	}

	return fileNameList, nil
}
