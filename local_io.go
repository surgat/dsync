package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type localIo struct {
	root string
}

func (this *localIo) getFullPath(partPath string) string {
	return this.root + partPath
}

func (this *localIo) Walk(walkFn WalkFunc) {
	var visit = func(fullPath string, fileInfo os.FileInfo, err error) error {
		path := strings.Replace(fullPath, this.root, "", -1)
		return walkFn(path, fileInfo, err)
	}
	filepath.Walk(this.root, visit)
}

func (this *localIo) Hash(path string) string {
	filePath := this.getFullPath(path)

	f, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(h.Sum(nil))
}

func (this *localIo) Remove(path string) {
	filePath := this.getFullPath(path)
	os.RemoveAll(filePath)
}

func (this *localIo) Mkdir(path string) {
	filePath := this.getFullPath(path)
	os.MkdirAll(filePath, 0755)
}

func (this *localIo) Rename(fromPath string, toPath string) {
	fromFullPath := this.getFullPath(fromPath)
	toFullPath := this.getFullPath(toPath)

	os.Rename(fromFullPath, toFullPath)
}

func (this *localIo) FileReader(path string) io.Reader {
	filePath := this.getFullPath(path)

	reader, err := os.OpenFile(filePath, os.O_RDONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}

	return reader
}

func (this *localIo) FileWriter(path string) io.Writer {
	filePath := this.getFullPath(path)

	writer, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}

	return writer
}

func CreateLocalIo(root string) IoInterface {
	io := new(localIo)
	io.root = root

	return io
}
