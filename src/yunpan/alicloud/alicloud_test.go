package alicloud

var client = &Client{
	AccessToken:   GetAccessToken(),
	BaseApiURL:    "http://api.yunpan.alibaba.com/api",
	LocalBaseDir:  "local_backup",
	RemoteBaseDir: "raspberry_pi",
}
