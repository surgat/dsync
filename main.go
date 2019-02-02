// dsync project main.go
package main

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"io"
	"log"
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

type pathInfo struct {
	fullPath string
	path     string
	fileInfo os.FileInfo
	hash     hash.Hash
}

func loadPathInfos(location string) map[string][]*pathInfo {

	storage := make(map[string][]*pathInfo)

	var visit = func(fullPath string, fileInfo os.FileInfo, err error) error {

		info := new(pathInfo)
		info.fullPath = fullPath
		info.path = strings.Replace(fullPath, location, "", -1)
		info.fileInfo = fileInfo

		if !fileInfo.IsDir() {

			f, err := os.Open(fullPath)
			if err != nil {
				return nil
				log.Fatal(err)
			}
			defer f.Close()

			h := md5.New()
			if _, err := io.Copy(h, f); err != nil {
				return nil
				log.Fatal(err)
			}

			md5String := hex.EncodeToString(h.Sum(nil))

			storage[md5String] = append(storage[md5String], info)
			//			storage[md5String][] = info
		}

		return nil
	}

	filepath.Walk(location, visit)

	return storage
}

func main() {
	flag.Parse()

	args := flag.Args()
	len := len(args)
	if len != 2 {
		panic("Not fount from and to")
	}

	from := args[0]

	fromPathInfos := loadPathInfos(from)

	to := args[1]
	toPathInfos := loadPathInfos(to)

	for h, fromPathInfo := range fromPathInfos {
		//		fmt.Println(h)
		if toPathInfo, ok := toPathInfos[h]; ok {
			fmt.Println(fromPathInfo.path + "   " + toPathInfo.path)
			if fromPathInfo.path == toPathInfo.path {
				fmt.Println("Ok")
			}
		}
	}

	//    fromFiles := map[string]*pathInfo
	//	err := filepath.Walk(from, loadPathInfos(&fromFiles))
	//	if err != nil {
	//		fmt.Printf("filepath.Walk() returned %v\n", err)
	//	}

}
