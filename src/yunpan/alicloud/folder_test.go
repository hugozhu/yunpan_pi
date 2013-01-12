package alicloud

import (
	"testing"
)

func TestFolderList(t *testing.T) {
	_, err := client.FolderList(0)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewFolder(t *testing.T) {
	folder, err := client.MakeFolder(0, "test folder")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(folder.Name)
	folder2, err2 := client.MakeFolder(folder.Id, "test folder 2")
	if err2 != nil {
		t.Fatal(err2)
	}

	folder2, err2 = client.RemoveFolder(folder2.Id)
	if err2 != nil {
		t.Fatal(err)
	}

	folder, err = client.RemoveFolder(folder.Id)
	if err != nil {
		t.Fatal(err)
	}
}
