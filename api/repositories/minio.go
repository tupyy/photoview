package repositories

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	s3Preffix = "s3://"
	pathSep   = "/"
)

type MinioReader struct {
	client *minio.Client
}

func NewMinioReader(url, accessKey, secretKey string) (*MinioReader, error) {
	minioClient, err := minio.New(url, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &MinioReader{minioClient}, nil
}

func (m *MinioReader) Stat(name string) (os.FileInfo, error) {
	if name == "" {
		return nil, os.ErrNotExist
	}

	bucket, prefix, filename := m.split(name)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if filename == "" {
		if prefix == "" {
			return m.statBucket(ctx, bucket)
		}
		return m.statPrefix(ctx, bucket, prefix)
	}

	return m.statFile(ctx, bucket, prefix, filename)
}

func (m *MinioReader) Lstat(name string) (os.FileInfo, error) {
	return m.Stat(name)
}

func (m *MinioReader) Open(filepath string) (File, error) {
	bucket, prefix, filename := m.split(filepath)

	filename = path.Join(prefix, filename)

	obj, err := m.client.GetObject(context.Background(), bucket, filename, minio.GetObjectOptions{})
	if err != nil {
		minioErr, ok := err.(minio.ErrorResponse)
		if ok && minioErr.StatusCode == 404 {
			return nil, os.ErrNotExist
		}
		return nil, err
	}

	return obj, nil
}

func (m *MinioReader) IsAbs(path string) bool {
	return true
}

func (m *MinioReader) ReadDir(path string) ([]os.DirEntry, error) {
	if path == "" {
		return []os.DirEntry{}, os.ErrNotExist
	}

	bucket, prefix := m.getBucket(path)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	found, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return []os.DirEntry{}, os.ErrInvalid
	}
	if !found {
		return []os.DirEntry{}, os.ErrNotExist
	}

	if prefix != "" && !strings.HasSuffix(prefix, pathSep) {
		prefix = fmt.Sprintf("%s/", prefix)
	}

	objectCh := m.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix: prefix,
	})

	entries := make([]os.DirEntry, 0, len(objectCh))
	for object := range objectCh {
		if object.Err != nil {
			return []os.DirEntry{}, object.Err
		}
		entries = append(entries, os.DirEntry(MinioDirEntry{object}))
	}

	return entries, nil
}

func (m *MinioReader) getBucket(path string) (bucket, prefix string) {
	path = strings.TrimPrefix(path, pathSep)

	if strings.Contains(path, pathSep) {
		parts := strings.Split(path, pathSep)
		return parts[0], strings.Join(parts[1:], pathSep)
	}

	return path, ""
}

func (m *MinioReader) statFile(ctx context.Context, bucket, prefix, filename string) (os.FileInfo, error) {
	if prefix != "" {
		filename = path.Join(prefix, filename)
	}
	objInfo, err := m.client.StatObject(ctx, bucket, filename, minio.StatObjectOptions{})
	if err != nil {
		minioErr, ok := err.(minio.ErrorResponse)
		if ok && minioErr.StatusCode == 404 {
			return nil, os.ErrNotExist
		}
		return nil, err
	}

	return FileInfo{objInfo}, nil
}

func (m *MinioReader) statPrefix(ctx context.Context, bucket, prefix string) (os.FileInfo, error) {
	objectCh := m.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix: prefix,
	})

	found := false
	for range objectCh {
		found = true
		break
	}

	if !found {
		return nil, os.ErrNotExist
	}

	return FolderInfo(prefix), nil
}

func (m *MinioReader) statBucket(ctx context.Context, bucket string) (os.FileInfo, error) {
	found, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, os.ErrNotExist
	}
	return FolderInfo(bucket), nil
}

// split splits the path from "/bucket/prefix1/prefix2/filename" into parts
func (m *MinioReader) split(name string) (bucket, prefix, filename string) {
	name = strings.TrimPrefix(name, pathSep)

	if !strings.Contains(name, pathSep) {
		return name, "", ""
	}

	parts := strings.Split(name, pathSep)
	if len(parts) == 1 {
		bucket = parts[0]
		return
	}

	if len(parts) == 2 {
		if strings.Contains(parts[1], ".") { // weak implementation. we consider file to have a dot in the name
			return parts[0], "", parts[1]
		} else {
			return parts[0], parts[1], ""
		}
	}

	bucket = parts[0]
	prefix = strings.Join(parts[1:len(parts)-1], pathSep)
	filename = path.Base(name)

	return
}

// MinioDirEntry implements os.DirEntry
type MinioDirEntry struct {
	obj minio.ObjectInfo
}

func (d MinioDirEntry) Name() string {
	if strings.Contains(d.obj.Key, pathSep) {
		return path.Base(d.obj.Key)
	}
	return d.obj.Key
}

func (d MinioDirEntry) IsDir() bool {
	return strings.HasSuffix(d.obj.Key, pathSep)
}

func (d MinioDirEntry) Type() os.FileMode {
	if d.IsDir() {
		return os.ModeDir
	}
	return fs.ModePerm
}

func (d MinioDirEntry) Info() (os.FileInfo, error) {
	return os.FileInfo(FileInfo{d.obj}), nil
}

func (d MinioDirEntry) Size() int64 {
	return d.Size()
}

// FileInfo implements os.FileInfo
type FileInfo struct {
	obj minio.ObjectInfo
}

func (f FileInfo) Name() string {
	if strings.HasSuffix(f.obj.Key, pathSep) {
		return strings.TrimSuffix(f.obj.Key, pathSep)
	}
	return f.obj.Key
}

func (f FileInfo) Size() int64 {
	return f.obj.Size
}

func (f FileInfo) Mode() os.FileMode {
	return os.FileMode(0777)
}

func (f FileInfo) ModTime() time.Time {
	return f.obj.LastModified
}

func (f FileInfo) IsDir() bool {
	return strings.HasSuffix(f.obj.Key, pathSep)
}

func (f FileInfo) Sys() any {
	return nil
}

type FolderInfo string

func (b FolderInfo) Name() string {
	return string(b)
}

func (b FolderInfo) Size() int64 {
	return int64(4096)
}

func (b FolderInfo) Mode() os.FileMode {
	return os.FileMode(0777)
}

func (b FolderInfo) ModTime() time.Time {
	return time.Now()
}

func (b FolderInfo) IsDir() bool {
	return true
}

func (b FolderInfo) Sys() any {
	return nil
}
