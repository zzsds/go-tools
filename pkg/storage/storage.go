package storage

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"mime/multipart"
	"net/url"
	"path"
	"strings"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
	// DefaultStorage Storage = new(noopStorage)
)

// Storage is a data storage interface
type Storage interface {
	Upload(*multipart.FileHeader, ...string) (*Resource, error)
	Delete(string) error
	// String returns the name of the implementation.
	String() string
}

type Resource struct {
	Url  *url.URL `json:"-"`
	Host string   `json:"host"`
	File string   `json:"file"`
	Path string   `json:"path"`
	Name string   `json:"name"`
	Size int64    `json:"size"`
	Ext  string   `json:"ext"`
}

func SetPath(path ...string) string {
	return strings.Join(append(path, time.Now().Format("20060102")), "/")
}

func ResetName(name string) string {
	// 读取文件后缀
	ext := path.Ext(name)
	// 读取文件名并加密
	name = strings.TrimSuffix(name, ext)
	h := md5.New()
	h.Write([]byte(name))
	name = hex.EncodeToString(h.Sum(nil))
	// 拼接新文件名
	return fmt.Sprintf("%s%s", name+"_"+time.Now().Format("20060102150405"), ext)
}
