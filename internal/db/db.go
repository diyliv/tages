package db

import (
	"io/ioutil"

	"github.com/diyliv/tages/pkg/file"
)

type Database interface {
	Store(name string, fb *file.FileBuffer) error
}

type Storage struct {
}

func NewDb() *Storage {
	return &Storage{}
}

func (s *Storage) Store(name string, file *file.FileBuffer) error {
	if err := ioutil.WriteFile(name, file.Buffer.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}
