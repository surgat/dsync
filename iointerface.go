package main

import (
	"io"
	"os"
)

type IoInterface interface {
	Walk(walkFn WalkFunc)
	Hash(path string) string
	Remove(path string)
	Mkdir(path string)
	Rename(fromPath string, toPath string)
	FileReader(path string) io.Reader
	FileWriter(path string) io.Writer
}

type WalkFunc func(path string, info os.FileInfo, err error) error
