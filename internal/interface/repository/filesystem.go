package repository

type IFilesystem interface {
	LoadTextFile(fileName string) (string, error)
}
