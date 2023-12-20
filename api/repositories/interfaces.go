package repositories

import (
	"os"

	"github.com/photoview/photoview/api/utils"
)

type RepositoryReader interface {
	Stat(name string) (os.FileInfo, error)
	Lstat(name string) (os.FileInfo, error)
	ReadDir(path string) ([]os.DirEntry, error)
	Open(path string) (File, error)
	IsAbs(path string) bool
}

type File interface {
	Read(b []byte) (n int, err error)
	Close() error
}

func GetDataRepository() RepositoryReader {
	if utils.EnvMinio.GetBool() {
		m, err := NewMinioReader("localhost:9000", utils.EnvMinioAccessKey.GetValue(), utils.EnvMinioSecretKey.GetValue())
		if err != nil {
			panic(err)
		}
		return m
	}
	return NewFileSystemReader()
}
