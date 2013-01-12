package alicloud

import (
	"os"
	"path/filepath"
)

var client = &Client{
	AccessToken:   GetAccessToken(),
	BaseApiURL:    "http://api.yunpan.alibaba.com/api",
	LocalBaseDir:  filepath.Join(os.Getenv("PWD"), "local_backup"),
	RemoteBaseDir: "raspberry_pi",
}
