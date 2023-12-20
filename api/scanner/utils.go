package scanner

import (
	"os"
	"path/filepath"

	"github.com/photoview/photoview/api/repositories"
	"github.com/pkg/errors"
)

// IsDirSymlink checks that the given path is a symlink and resolves to a
// directory.
func IsDirSymlink(path string) (bool, error) {
	isDirSymlink := false

	fileInfo, err := repositories.GetDataRepository().Lstat(path)
	if err != nil {
		return false, errors.Wrapf(err, "could not stat %s", path)
	}

	//Resolve symlinks
	if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		resolvedPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			return false, errors.Wrapf(err, "Cannot resolve linktarget of %s, ignoring it", path)
		}

		resolvedFile, err := os.Stat(resolvedPath)
		if err != nil {
			return false, errors.Wrapf(err, "Cannot get fileinfo of linktarget %s of symlink %s, ignoring it", resolvedPath, path)
		}
		isDirSymlink = resolvedFile.IsDir()

		return isDirSymlink, nil
	}

	return false, nil
}
