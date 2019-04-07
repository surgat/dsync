package main

import (
	"os"
)

type pathInfo struct {
	fullPath string
	path     string
	fileInfo os.FileInfo
	readed   bool
}

func loadPathInfos(io IoInterface) map[string]*pathInfo {

	storage := make(map[string]*pathInfo)

	var visit = func(path string, fileInfo os.FileInfo, err error) error {
		info := new(pathInfo)
		info.path = path
		info.fileInfo = fileInfo
		info.readed = false

		storage[info.path] = info

		return nil
	}

	io.Walk(visit)

	return storage
}

type PathInfoCompare struct {
	withHash bool
	fromIo   IoInterface
	toIo     IoInterface
}

func (service *PathInfoCompare) compare(from *pathInfo, to *pathInfo) bool {
	isEqualsSize := from.fileInfo.Size() == to.fileInfo.Size()

	if isEqualsSize {
		fromHash := service.fromIo.Hash(from.path)
		return fromHash == service.toIo.Hash(to.path)
	}

	return isEqualsSize
}
