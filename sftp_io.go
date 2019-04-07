package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type sftpIo struct {
	client *sftp.Client
	root   string
}

func (this *sftpIo) getFullPath(partPath string) string {
	return this.root + partPath
}

func (this *sftpIo) Walk(walkFn WalkFunc) {
	walker := this.client.Walk(this.root)

	for walker.Step() {
		path := strings.Replace(walker.Path(), this.root, "", -1)

		walkFn(path, walker.Stat(), walker.Err())
	}
}

func (this *sftpIo) Hash(path string) string {
	filePath := this.getFullPath(path)

	f, err := this.client.Open(filePath)

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

func (this *sftpIo) Remove(path string) {
	filePath := this.getFullPath(path)
	this.client.Remove(filePath)
}

func (this *sftpIo) Mkdir(path string) {
	filePath := this.getFullPath(path)
	this.client.MkdirAll(filePath)
}

func (this *sftpIo) Rename(fromPath string, toPath string) {
	fromFullPath := this.getFullPath(fromPath)
	toFullPath := this.getFullPath(toPath)

	this.client.Rename(fromFullPath, toFullPath)
}

func (this *sftpIo) FileReader(path string) io.Reader {
	filePath := this.getFullPath(path)

	reader, err := this.client.OpenFile(filePath, os.O_RDONLY)
	if err != nil {
		log.Fatal(err)
	}

	return reader
}

func (this *sftpIo) FileWriter(path string) io.Writer {
	filePath := this.getFullPath(path)

	writer, err := this.client.OpenFile(filePath, os.O_WRONLY|os.O_CREATE)
	if err != nil {
		log.Fatal(err)
	}

	return writer
}

func CreateSftpIo(root string, addr string, config *ssh.ClientConfig) IoInterface {
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		log.Fatal("Failed to sftp: ", err)
	}

	io := new(sftpIo)
	io.root = root
	io.client = sftpClient

	return io
}
