package fs

import (
	"io/ioutil"
	"os"
	"time"
)

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func ChangeModTime(path string, unix_time int64) error {
	mtime := time.Unix(unix_time, 0)
	return os.Chtimes(path, mtime, mtime)
}

func ListFiles(path string, accept_filter func(os.FileInfo) bool) ([]os.FileInfo, []os.FileInfo, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, nil, err
	}
	dirList := make([]os.FileInfo, 0)
	fileList := make([]os.FileInfo, 0)
	for _, f := range files {
		if !accept_filter(f) {
			continue
		}
		if f.IsDir() {
			dirList = append(dirList, f)
		} else {
			fileList = append(fileList, f)
		}
	}
	return fileList, dirList, err
}
