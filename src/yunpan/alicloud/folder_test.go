package alicloud

import (
	"testing"
)

func TestFolderList(t *testing.T) {
	_, err := client.FolderList(0)
	if err != nil {
		t.Error(err)
	}
}
