package local

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/zzsds/kratos-tools/pkg/storage"
)

type Local struct {
	opts Options
}

func NewLocal(opts ...Option) storage.Storage {
	return &Local{opts: newOptions(opts...)}
}

// UploadFile 上传文件
func (l *Local) Upload(file *multipart.FileHeader, work ...string) (*storage.Resource, error) {
	name := storage.ResetName(file.Filename)
	p := storage.SetPath(work...)
	resource := &storage.Resource{
		Path: p,
		Name: name,
		File: fmt.Sprintf("%s/%s", p, name),
		Ext:  path.Ext(file.Filename),
		Size: file.Size,
		Host: l.opts.Path,
	}
	// 尝试创建此路径
	mkdirErr := os.MkdirAll(l.opts.Path, os.ModePerm)
	if mkdirErr != nil {
		return nil, errors.New("function os.MkdirAll() Filed, err:" + mkdirErr.Error())
	}

	f, openError := file.Open() // 读取文件
	if openError != nil {
		return nil, errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭

	out, createErr := os.Create(resource.File)
	if createErr != nil {

		return nil, errors.New("function os.Create() Filed, err:" + createErr.Error())
	}
	defer out.Close() // 创建文件 defer 关闭

	_, copyErr := io.Copy(out, f) // 传输（拷贝）文件
	if copyErr != nil {
		return nil, errors.New("function io.Copy() Filed, err:" + copyErr.Error())
	}
	return resource, nil
}

// DeleteFile 删除文件
func (l *Local) Delete(key string) error {
	p := l.opts.Path + "/" + key
	if strings.Contains(p, l.opts.Path) {
		if err := os.Remove(p); err != nil {
			return errors.New("本地文件删除失败, err:" + err.Error())
		}
	}
	return nil
}

func (l *Local) String() string {
	return "local"
}
