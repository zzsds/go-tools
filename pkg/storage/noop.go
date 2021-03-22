package storage

import (
	"mime/multipart"
)

type noopStorage struct{}

func NewNoop() Storage {
	noop := noopStorage{}
	return &noop
}

func (n *noopStorage) Init() error {
	return nil
}

func (n *noopStorage) Upload(*multipart.FileHeader, ...string) (*Resource, error) {
	return &Resource{}, nil
}

func (n *noopStorage) Delete(key string) error {
	return nil
}

func (n *noopStorage) String() string {
	return "noop"
}
