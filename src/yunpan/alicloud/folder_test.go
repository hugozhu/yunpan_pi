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

func TestFolderOperations(t *testing.T) {
	folder, err := client.MakeFolder(0, "test folder")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(folder.Name)
	folder2, err2 := client.MakeFolder(folder.Id, "test folder 2")
	if err2 != nil {
		t.Fatal(err2)
	}

	folder3, err3 := client.RenameFolder(folder2.Id, "test folder 3")
	if err3 != nil {
		t.Fatal(err3)
	}

	if folder3.Name != "test folder 3" {
		t.Fatal("failed to rename")
	}

	folder4, err4 := client.MakeFolder(0, "test folder 4")
	if err4 != nil {
		t.Fatal(err4)
	}

	folder5, err5 := client.MoveFolder(folder4.Id, folder3.Id)
	if err5 != nil || !folder5.Suc {
		t.Fatal(err5)
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
