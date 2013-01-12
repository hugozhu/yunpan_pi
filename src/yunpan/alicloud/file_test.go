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
	create_file(file, DEFAULT_CHUNK_SIZE*2+1024)
	t.Log("prepare file:", file)
	fileInfo, err := client.CreateFile(0, file)
	if err != nil {
		t.Fatal("remote api error " + err.Error())
	}
	defer func() {
		client.FileRemove(fileInfo.Id)
	}()
	succ := true

	if len(fileInfo.Chunks) != 3 {
		t.Fatal("the number of chunks is not as expected")
	}

	var offset int64
	for _, chunk := range fileInfo.Chunks {
		r, e := client.UploadChunk(chunk.Id, file, offset, chunk.Size)
		if !r || e != nil {
			succ = false
			t.Fatal(e)
		}
		offset += chunk.Size
	}

	if succ {
		client.CommitUpload(fileInfo.Id, fileInfo.UpdateVersion)
	}
}
