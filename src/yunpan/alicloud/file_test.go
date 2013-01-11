package alicloud

import (
	"os"
	"path/filepath"
	"testing"
)

func create_file(name string, size int) (string, error) {
	file, err := os.Create(name)
	file.Name()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = 0
	}
	file.Write(buf)

	return file.Name(), nil
}

func TestFileUpload(t *testing.T) {
	file := filepath.Join(os.Getenv("PWD"), client.LocalBaseDir, "testfile 123.txt")
	create_file(file, 1024)
	t.Log(file)
	fileInfo, err := client.CreateFile(0, file)
	if err != nil {
		t.Error("got api err " + err.Error())
	}
	defer func() {
		client.FileRemove(fileInfo.Id)
	}()
	succ := true
	for _, chunk := range fileInfo.Chunks {
		r, e := client.UploadChunk(chunk.Id, chunk.Size, file, 0, chunk.Size)
		if !r || e != nil {
			succ = false
			break
		}
	}

	if succ {
		client.CommitUpload(fileInfo.Id, fileInfo.UpdateVersion)
	} else {
		t.Error("Failed to upload file")
	}
}
