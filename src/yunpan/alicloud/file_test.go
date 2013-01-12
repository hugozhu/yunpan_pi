package alicloud

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
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
		buf[i] = byte(i)
	}
	file.Write(buf)

	return file.Name(), nil
}

func TestFileUpload(t *testing.T) {
	file := filepath.Join(os.Getenv("PWD"), client.LocalBaseDir, fmt.Sprintf("testfile_%d.txt", time.Now().Unix()))
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
			t.Error(e)
		}
	}

	if succ {
		client.CommitUpload(fileInfo.Id, fileInfo.UpdateVersion)
	}
}
