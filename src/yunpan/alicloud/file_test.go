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
	random_filename := fmt.Sprintf("testfile_%d.txt", time.Now().Unix())
	file := filepath.Join(client.LocalBaseDir, random_filename)
	downloaded_file := filepath.Join(client.LocalBaseDir, "downloaded_"+random_filename)
	create_file(file, DEFAULT_CHUNK_SIZE*0+1024)
	t.Log("prepare file:", file)
	fileInfo, err := client.CreateFile(0, file)
	if err != nil {
		t.Fatal("remote api error " + err.Error())
	}
	defer func() {
		os.Remove(file)
		os.Remove(downloaded_file)
		client.FileRemove(fileInfo.Id)
	}()

	// if len(fileInfo.Chunks) != 3 {
	// 	t.Fatal("the number of chunks is not as expected")
	// }

	var offset int64
	for _, chunk := range fileInfo.Chunks {
		r, e := client.UploadChunk(chunk.Id, file, offset, chunk.Size)
		if !r || e != nil {
			t.Fatal(e)
		}
		offset += chunk.Size
	}

	var fileInfo2 *FileInfo
	fileInfo2, err = client.CommitUpload(fileInfo.Id, fileInfo.UpdateVersion)
	if err != nil {
		t.Fatal("Failed to commit upload", err)
	}

	if fileInfo.Size != fileInfo2.Size {
		t.Error("size is not same")
	}

	fileInfo3, err3 := client.FileInfo(fileInfo.Id, fileInfo2.FullName, 3)
	if err3 != nil {
		t.Fatal(err3)
	}

	if fileInfo.Size != fileInfo3.Size {
		t.Error("size is not same")
	}

	// if fileInfo.ModifyTime != fileInfo2.ModifyTime {
	// 	t.Error("modify time is not same", fileInfo.ModifyTime, fileInfo2.ModifyTime)
	// }

	err = client.DownloadFile(fileInfo, downloaded_file)
	if err != nil {
		t.Fatal(err)
	}

	if md5_hash(file) != md5_hash(downloaded_file) {
		t.Fatal("Downloaded file's md5 is not same as uploaded one")
	}

}
