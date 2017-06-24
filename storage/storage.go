package storage

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Backend interface {
	Store(int64, []byte) error
	Get(int64) ([]byte, error)
}

type FileBackend struct {
	BaseDir string
}

func (b *FileBackend) Init() {
	err := os.MkdirAll(b.BaseDir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("FileBackend ready at %s", b.BaseDir)
}

func (b *FileBackend) Store(id int64, data []byte) error {
	p := b.path(id)
	log.Printf("Writing %d bytes to %s", len(data), p)

	return ioutil.WriteFile(p, data, 0644)
}

func (b *FileBackend) Get(id int64) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(b.BaseDir, fmt.Sprintf("%d", id)))
}

func (b *FileBackend) path(id int64) string {
	return filepath.Join(b.BaseDir, fmt.Sprintf("%d", id))
}
