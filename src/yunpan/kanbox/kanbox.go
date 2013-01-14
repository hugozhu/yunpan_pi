package kanbox

type Client struct {
	AccessToken     string
	BaseApiURL      string
	LocalBaseDir    string
	RemoteBaseDir   string
	RemoteBaseDirId int64
}
