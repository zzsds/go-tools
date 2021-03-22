package aliyun

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"path"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/zzsds/kratos-tools/pkg/storage"
)

type AliYun struct {
	opts       Options
	client     *oss.Client
	bucket     *oss.Bucket
	class, acl oss.Option // 根据配置文件进行指定存储类型
}

func NewAliYun(opts ...Option) storage.Storage {
	ali := AliYun{
		opts: newOptions(opts...),
	}
	if err := ali.Init(); err != nil {
		log.Fatal(err)
	}

	return &ali
}

func (a *AliYun) Init() error {
	client, err := oss.New(a.opts.Endpoint, a.opts.AccessKey, a.opts.SecretKey, oss.Timeout(10, 120))
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(a.opts.Bucket) // 获取存储空间
	if err != nil {
		return err
	}
	a.client = client
	a.bucket = bucket

	var class oss.Option
	switch a.opts.ClassType { // 根据配置文件进行指定存储类型
	case "Standard": // 指定存储类型为标准存储
		class = oss.ObjectStorageClass(oss.StorageStandard)
	case "IA": // 指定存储类型为很少访问存储
		class = oss.ObjectStorageClass(oss.StorageIA)
	case "Archive": // 指定存储类型为归档存储。
		class = oss.ObjectStorageClass(oss.StorageArchive)
	case "ColdArchive": // 指定存储类型为归档存储。
		class = oss.ObjectStorageClass(oss.StorageColdArchive)
	default: // 无匹配结果就是标准存储
		class = oss.ObjectStorageClass(oss.StorageStandard)
	}

	var objectAcl oss.Option
	switch a.opts.AclType { // 根据配置文件进行指定访问权限
	case "private": // 指定访问权限为私有读写
		objectAcl = oss.ObjectACL(oss.ACLPrivate) // 指定访问权限为公共读
	case "public-read":
		objectAcl = oss.ObjectACL(oss.ACLPublicRead) // 指定访问权限为公共读
	case "public-read-write":
		objectAcl = oss.ObjectACL(oss.ACLPublicReadWrite) // 指定访问权限为公共读写
	case "default":
		objectAcl = oss.ObjectACL(oss.ACLDefault) // 指定访问权限为公共读
	default:
		objectAcl = oss.ObjectACL(oss.ACLPrivate) // 默认为访问权限为公共读
	}

	a.class = class
	a.acl = objectAcl
	return nil
}

// Upload 上传文件
func (a *AliYun) Upload(file *multipart.FileHeader, work ...string) (*storage.Resource, error) {
	f, openError := file.Open() // 读取文件
	if openError != nil {
		return nil, errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	name := storage.ResetName(file.Filename)
	p := storage.SetPath(work...)
	resource := &storage.Resource{
		Path: p,
		Name: name,
		File: fmt.Sprintf("%s/%s", p, name),
		Ext:  path.Ext(file.Filename),
		Size: file.Size,
		Host: a.opts.Path,
	}
	// 获取文件类型
	putErr := a.bucket.PutObject(resource.File, f, a.class, oss.ContentType(file.Header.Get("content-type")), a.acl) // 上传
	if putErr != nil {
		return nil, errors.New("function bucket.PutObject() Filed, err:" + putErr.Error())
	}

	url, _ := url.Parse(a.opts.Path)
	url.Path = resource.File
	resource.Url = url
	return resource, nil
}

// Delete 删除文件
func (a *AliYun) Delete(key string) error {
	// 删除单个文件。objectName表示删除OSS文件时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	// 如需删除文件夹，请将objectName设置为对应的文件夹名称。如果文件夹非空，则需要将文件夹下的所有object删除后才能删除该文件夹。
	if err := a.bucket.DeleteObject(key); err != nil {
		return errors.New("function bucket.DeleteObject() Filed, err:" + err.Error())
	}
	return nil
}

func (a *AliYun) String() string {
	return "aliyun"
}
