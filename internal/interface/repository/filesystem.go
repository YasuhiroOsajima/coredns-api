package repository

type IFilesystem interface {
	LoadTextFile(fileName string) (string, error)
	WriteTextFile(name, fileInfo string) error
	DeleteFile(fileName string) error
	GetFilenameList(directory string) ([]string, error)
}
