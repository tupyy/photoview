package repositories

import (
	"os"
	"path"
)

type FileSystemReader string

func NewFileSystemReader() *FileSystemReader {
	return new(FileSystemReader)
}

func (fr *FileSystemReader) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (fr *FileSystemReader) ReadDir(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

func (fr *FileSystemReader) Open(path string) (File, error) {
	return os.Open(path)
}

func (fr *FileSystemReader) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(name)
}

func (fr *FileSystemReader) IsAbs(p string) bool {
	return path.IsAbs(p)
}
