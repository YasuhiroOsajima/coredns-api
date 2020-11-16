package repository

type FilesystemRepository struct {
	filesystem IFilesystem
}

func NewFileRepository(fs IFilesystem) *FilesystemRepository {
	return &FilesystemRepository{fs}
}
