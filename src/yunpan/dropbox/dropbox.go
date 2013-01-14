package dropbox

type Client struct {
	AccessToken     string
	BaseApiURL      string
	LocalBaseDir    string
	RemoteBaseDir   string
	RemoteBaseDirId int64
}
