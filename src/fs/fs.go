package fs

import (
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
