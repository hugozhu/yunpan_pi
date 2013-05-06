package main

import (
	"fs"
	"mylog"
	"os"
	"path/filepath"
	"strings"
	"yunpan/alicloud"
)

var c = &alicloud.Client{
	AccessToken:     alicloud.GetAccessToken(), //oauth token, get yours from ali clould drive api doc: http://api.yunpan.alibaba.com/
	BaseApiURL:      "http://api.yunpan.alibaba.com/api",
	LocalBaseDir:    filepath.Join(os.Getenv("PWD"), "local_backup"), //the folder you want to sync with cloud disk
	RemoteBaseDir:   "working/cosocket/lua",                          //the sub folder you want to sync with local disk
	RemoteBaseDirId: 3094091,                                         //the sub folder Id, can be set to 0
}

var log = mylog.New(os.Stdout)

var IgnoreFileExtensions = []string{".tmp", ".swp"}

func init() {
	log.DebugEnabled = true

	if c.RemoteBaseDirId > 0 {
		return
	}
	if filepath.Clean(c.RemoteBaseDir) == "." {
		c.RemoteBaseDir = ""
		return
	}
	parts := strings.Split(c.RemoteBaseDir, "/")
	var dirId int64 = 0
	for _, part := range parts {
		fileList, err := c.FolderList(dirId)
		panic_if_error(err)
		found := false
		for _, d := range fileList.Dirs {
			if d.Name == part {
				dirId = d.Id
				found = true
				break
			}
		}
		if !found {
			d, err := c.MakeFolder(dirId, part)
			panic_if_error(err)
			dirId = d.Id
		}
	}
	c.RemoteBaseDirId = dirId

	if c.RemoteBaseDirId == 0 && c.RemoteBaseDir != "" {
		panic("fatal error")
	}
}

func panic_if_error(err error) {
	if err != nil {
		panic(err)
	}
}

func accept_filter(s os.FileInfo) bool {
	ext := filepath.Ext(s.Name())
	if !s.IsDir() {
		for _, p := range IgnoreFileExtensions {
			if ext == p {
				return false
			}
		}
	}
	return true
}

func SyncFolder(dirId int64, dirPath string, dirModTime int64) {
	ok, err := fs.Exists(dirPath)
	panic_if_error(err)
	if !ok {
		os.Mkdir(dirPath, 0755)
	}

	fileList, err := c.FolderList(dirId)
	panic_if_error(err)

	cloudFiles := make(map[string]*alicloud.File)
	cloudFolders := make(map[string]*alicloud.Folder)

	for _, remoteFile := range fileList.Files {
		localFilePath := filepath.Join(dirPath, remoteFile.GetFullName())
		modTime2 := remoteFile.ModifyTime / 1000
		fileInfo1, err := os.Stat(localFilePath)
		var modTime1 int64
		if err != nil && !os.IsNotExist(err) {
			panic_if_error(err)
		}
		if fileInfo1 != nil {
			modTime1 = fileInfo1.ModTime().Unix()
		}

		log.Debugf("+ %s local timestamp: %d remote: %d", localFilePath, modTime1, modTime2)

		if modTime2 > modTime1 {
			fileInfo2, err := c.FileInfo(remoteFile.Id, "", 3)
			panic_if_error(err)

			log.Info("Download:", localFilePath)
			fileInfo2.ModifyTime = remoteFile.ModifyTime
			c.DownloadFile(fileInfo2, localFilePath)
		} else if modTime1 > modTime2 {
			fileInfo2, err := c.FileInfo(remoteFile.Id, "", 3)
			panic_if_error(err)
			fileInfo2, err = c.ModifyFile(remoteFile.Id, dirId, localFilePath, fileInfo2)
			panic_if_error(err)
			upload_file(localFilePath, dirId, fileInfo2)

		}
		cloudFiles[localFilePath] = remoteFile
	}

	for _, d := range fileList.Dirs {
		p := filepath.Join(dirPath, d.Name)
		SyncFolder(d.Id, p, d.ModifyTime)
		cloudFolders[p] = d
	}

	//upload new files from local to cloud
	myFiles, myFolders, _ := fs.ListFiles(dirPath, accept_filter)

	for _, f := range myFiles {
		localFilePath := filepath.Join(dirPath, f.Name())
		if cloudFiles[localFilePath] == nil {
			//upload to cloud
			fileInfo, err := c.CreateFile(dirId, localFilePath)
			panic_if_error(err)
			upload_file(localFilePath, dirId, fileInfo)
		}
	}

	for _, f := range myFolders {
		localFilePath := filepath.Join(dirPath, f.Name())
		if cloudFolders[localFilePath] == nil {
			//upload to cloud
			newFolder, err := c.MakeFolder(dirId, f.Name())
			panic_if_error(err)
			upload_folder(localFilePath, newFolder.Id)
		}
	}

	if dirModTime > 0 {
		err := fs.ChangeModTime(dirPath, dirModTime/1000)
		panic_if_error(err)
	}
}

func upload_folder(localDirPath string, dirId int64) {
	myFiles, myFolders, err := fs.ListFiles(localDirPath, accept_filter)
	panic_if_error(err)
	for _, f := range myFiles {
		localFilePath := filepath.Join(localDirPath, f.Name())
		//upload to cloud
		fileInfo, err := c.CreateFile(dirId, localFilePath)
		panic_if_error(err)
		upload_file(localFilePath, dirId, fileInfo)
	}

	for _, f := range myFolders {
		localFilePath := filepath.Join(localDirPath, f.Name())
		//upload to cloud
		newFolder, err := c.MakeFolder(dirId, f.Name())
		panic_if_error(err)
		upload_folder(localFilePath, newFolder.Id)
	}
}

func upload_file(localFilePath string, dirId int64, fileInfo *alicloud.FileInfo) {
	var offset int64
	for _, chunk := range fileInfo.Chunks {
		r, e := c.UploadChunk(chunk.Id, localFilePath, offset, chunk.Size)
		if !r || e != nil {
			log.Error(e)
			panic(e)
		}
		offset += chunk.Size
	}
	fileInfo, err := c.CommitUpload(fileInfo.Id, fileInfo.UpdateVersion)
	panic_if_error(err)
	//update local file's last modified time, so we don't need sync next time
	newModTime := fileInfo.ModifyTime / 1000
	fs.ChangeModTime(localFilePath, newModTime)
	log.Info("Upload:", fileInfo.GetFullName(), fileInfo.Version, newModTime, fileInfo.ModifyTime/1000)
}

func main() {
	log.Info("Sync manager started at: ", c.RemoteBaseDirId)
	SyncFolder(c.RemoteBaseDirId, c.LocalBaseDir, 0)
	log.Info("Done")

}
