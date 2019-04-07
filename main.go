// dsync project main.go
package main

import (
	"flag"
	"fmt"
	"io"
)

type operationType int

const (
	COPY operationType = 0 + iota
	RENAME
	REMOVE
	MKDIR
)

var operationsList = map[operationType]string{
	COPY:   "COPY",
	RENAME: "RENAME",
	REMOVE: "REMOVE",
	MKDIR:  "MKDIR",
}

type operation struct {
	srcPath  string
	destPath string
	otype    operationType
}

func (o operationType) String() string {
	return operationsList[o]
}

var CompareService *PathInfoCompare

func main() {
	flag.Parse()

	args := flag.Args()
	len := len(args)
	if len != 2 {
		panic("Not fount from and to")
	}
	from := args[0]
	fromIo := CreateLocalIo(from)
	fromPathInfos := loadPathInfos(fromIo)

	to := args[1]
	toIo := CreateLocalIo(to)
	toPathInfos := loadPathInfos(toIo)

	CompareService := new(PathInfoCompare)
	CompareService.fromIo = fromIo
	CompareService.toIo = toIo

	typeOperations := buildOperations(fromPathInfos, toPathInfos, CompareService)

	for _, operation := range typeOperations[REMOVE] {
		toIo.Remove(operation.destPath)
		fmt.Printf("%s -> %s : %s\n", operation.srcPath, operation.destPath, operation.otype.String())
	}

	for _, operation := range typeOperations[MKDIR] {
		toIo.Mkdir(operation.destPath)
		fmt.Printf("%s -> %s : %s\n", operation.srcPath, operation.destPath, operation.otype.String())
	}

	for _, operation := range typeOperations[RENAME] {
		toIo.Rename(operation.srcPath, operation.destPath)
		fmt.Printf("%s -> %s : %s\n", operation.srcPath, operation.destPath, operation.otype.String())
	}

	for _, operation := range typeOperations[COPY] {
		reader := fromIo.FileReader(operation.srcPath)
		writer := toIo.FileWriter(operation.destPath)

		io.Copy(writer, reader)
		fmt.Printf("%s -> %s : %s\n", operation.srcPath, operation.destPath, operation.otype.String())
	}
}

func buildOperations(fromPathInfos, toPathInfos map[string]*pathInfo, CompareService *PathInfoCompare) map[operationType][]operation {
	operations := make(map[operationType][]operation, 0)
	for k, _ := range operationsList {
		operations[k] = make([]operation, 0)
	}

	for fromPath, fromPathInfo := range fromPathInfos {
		if fromPathInfo.fileInfo.IsDir() {
			fromPathInfo.readed = true

			if toPathInfo, ok := toPathInfos[fromPath]; ok {
				if toPathInfo.fileInfo.IsDir() {
					toPathInfo.readed = true
					continue
				}

				op := new(operation)
				op.srcPath = ""
				op.destPath = fromPathInfo.path
				op.otype = REMOVE
				operations[op.otype] = append(operations[op.otype], *op)
			}

			op := new(operation)
			op.srcPath = ""
			op.destPath = fromPathInfo.path
			op.otype = MKDIR
			operations[op.otype] = append(operations[op.otype], *op)

			continue
		}

		if toPathInfo, ok := toPathInfos[fromPath]; ok {
			if CompareService.compare(fromPathInfo, toPathInfo) {
				fromPathInfo.readed = true
				toPathInfo.readed = true

				continue
			}
		}

		for toPath, toPathInfo := range toPathInfos {
			if toPathInfo.fileInfo.IsDir() {
				continue
			}

			if CompareService.compare(fromPathInfo, toPathInfo) {

				fromPathInfo.readed = true
				toPathInfo.readed = true

				op := new(operation)
				op.srcPath = toPath
				op.destPath = fromPath
				op.otype = RENAME

				operations[op.otype] = append(operations[op.otype], *op)
				continue
			}
		}

		if fromPathInfo.readed == false {
			fromPathInfo.readed = true

			op := new(operation)
			op.srcPath = fromPath
			op.destPath = fromPath
			op.otype = COPY

			operations[op.otype] = append(operations[op.otype], *op)
		}
	}

	for toPath, toPathInfo := range toPathInfos {
		if toPathInfo.readed == false {
			op := new(operation)
			op.srcPath = ""
			op.destPath = toPath
			op.otype = REMOVE

			operations[op.otype] = append(operations[op.otype], *op)
		}
	}

	return operations
}
