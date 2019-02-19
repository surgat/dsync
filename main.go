// dsync project main.go
package main

import (
	//	"crypto/md5"
	//	"encoding/hex"
	"hash"
	//	"io"
	//	"log"
	"strings"

	//	"errors"

	"flag"
	"fmt"
	"os"
	"path/filepath"
)

//func init() {
//	flag.StringVar(&from, "gopher_type", defaultGopher, usage)
//	flag.StringVar(&gopherType, "g", defaultGopher, usage+" (shorthand)")
//
//	wordPtr := flag.String("word", "foo", "a string")
//}

type operationType int

const (
	COPY operationType = 0 + iota
	RENAME
	REMOVE
)

type operation struct {
	srcPath  string
	destPath string
	otype    operationType
}

func (o operationType) String() string {
	return [...]string{"COPY", "RENAME", "REMOVE"}[o]
}

type pathInfo struct {
	fullPath string
	path     string
	fileInfo os.FileInfo
	hash     hash.Hash
	readed   bool
}

func (p *pathInfo) compare(toP *pathInfo) bool {
	return p.fileInfo.Size() == toP.fileInfo.Size()
}

func loadPathInfos(location string) map[string]*pathInfo {

	storage := make(map[string]*pathInfo)

	var visit = func(fullPath string, fileInfo os.FileInfo, err error) error {

		info := new(pathInfo)
		info.fullPath = fullPath
		info.path = strings.Replace(fullPath, location, "", -1)
		info.fileInfo = fileInfo
		info.readed = false

		if !fileInfo.IsDir() {

			//			f, err := os.Open(fullPath)
			//			if err != nil {
			//				return nil
			//				log.Fatal(err)
			//			}
			//			defer f.Close()

			//			h := md5.New()
			//			if _, err := io.Copy(h, f); err != nil {
			//				return nil
			//				log.Fatal(err)
			//			}

			//			md5String := hex.EncodeToString(h.Sum(nil))
			//
			//			if _, ok := storage[md5String]; !ok {
			//				storage[md5String] = make(map[string]*pathInfo)
			//			}

			storage[info.path] = info
		}

		return nil
	}

	filepath.Walk(location, visit)

	return storage
}

func main() {
	flag.Parse()
	//	db, _ := buntdb.Open(":memory:")
	//
	//	err := db.Update(func(tx *buntdb.Tx) error {
	//		_, _, err := tx.Set("mykey", "myvalue", nil)
	//		return err
	//	})

	args := flag.Args()
	len := len(args)
	if len != 2 {
		panic("Not fount from and to")
	}
	to := args[1]
	toPathInfos := loadPathInfos(to)

	from := args[0]
	fromPathInfos := loadPathInfos(from)

	operations := make([]operation, 0)
	for fromPath, fromPathInfo := range fromPathInfos {
		if toPathInfo, ok := toPathInfos[fromPath]; ok {
			if fromPathInfo.compare(toPathInfo) {
				fromPathInfo.readed = true
				toPathInfo.readed = true

				continue
			}
		}

		for toPath, toPathInfo := range toPathInfos {
			if fromPathInfo.compare(toPathInfo) {

				fromPathInfo.readed = true
				toPathInfo.readed = true

				op := new(operation)
				op.srcPath = fromPath
				op.destPath = toPath
				op.otype = RENAME

				operations = append(operations, *op)
				continue
			}
		}

		if fromPathInfo.readed == false {
			fromPathInfo.readed = true

			op := new(operation)
			op.srcPath = fromPath
			op.destPath = fromPath
			op.otype = COPY

			operations = append(operations, *op)
		}
	}

	for toPath, toPathInfo := range toPathInfos {
		if toPathInfo.readed == false {
			op := new(operation)
			op.srcPath = ""
			op.destPath = toPath
			op.otype = REMOVE

			operations = append(operations, *op)
		}
	}

	for _, operation := range operations {
		fmt.Printf("%s -> %s : %s\n", operation.srcPath, operation.destPath, operation.otype.String())
	}
}
