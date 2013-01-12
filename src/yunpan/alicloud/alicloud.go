package alicloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

const DEFAULT_CHUNK_SIZE = 2097143

type Client struct {
	AccessToken   string
	BaseApiURL    string
	LocalBaseDir  string
	RemoteBaseDir string
}

type FileList struct {
	ErrorCode int
	HasError  bool
	Files     []*File
	Dirs      []*Dir
}

type Operationable struct {
	Operation  int
	ResultCode int
	Suc        bool
}

type Dir struct {
	Name       string
	Id         int64
	ModifyTime int32
	Operation  int
	ResultCode int
	Suc        bool
}

type File struct {
	Id         int64
	Size       int
	FileName   string
	Version    int64
	Extension  string
	Md5        string
	ResultCode int
	Suc        bool
	DirId      int64
	ModifyTime int64
}

type FileInfo struct {
	ChangedBy     int64    `json:"changedBy"`
	Chunks        []*Chunk `json:"chunks"`
	DirId         int64    `json:"dirId"`
	Direct        bool     `json:"direct"`
	Extension     string   `json:"extension"`
	FileAttribute int32
	FileName      string `json:"fileName"`
	FullName      string `json:"fullName"`
	Id            int64  `json:"id"`
	Md5           string `json:"md5"`
	ModifyTime    int64  `json:"modifyTime"`
	Operation     int    `json:"operation"`
	PlatformInfo  int    `json:"platformInfo"`
	ResultCode    int    `json:"resultCode"`
	Size          int64  `json:"size"`
	Suc           bool   `json:"suc"`
	UpdateVersion int64  `json:"updateVersion"`
	Version       int64  `json:"version"`
}

type Chunk struct {
	CheckSum   string `json:"checkSum"`
	GenerNext  bool   `json:"generNext"`
	GenerPre   bool   `json:"generPre"`
	Id         int64
	Index      int64  `json:"index"`
	Key        string `json:"key"`
	Md5        string `json:"md5"`
	NeedUpload bool   `json:"needUpload"`
	Next       int64  `json:"next"`
	Operation  int    `json:"operation"`
	Pre        int64  `json:"pre"`
	Size       int64  `json:"size"`
}

type ApiError struct {
	ErrorCode        int64  `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorURI         string `json:"error_uri"`
}

func (e ApiError) Error() string {
	return fmt.Sprintf("api error: [%d] %s", e.ErrorCode, e.ErrorDescription)
}

func NewApiError(s []byte) ApiError {
	var e ApiError
	err := json.Unmarshal(s, &e)
	if err != nil {
		return ApiError{ErrorCode: 0, ErrorDescription: "Json response parse error"}
	}
	return e
}

func (c *Client) GetCall(path string, params *url.Values) ([]byte, error) {
	params.Set("accessToken", c.AccessToken)
	queryString := params.Encode()
	if len(queryString) > 0 {
		queryString = "?" + queryString
	}
	req, err := http.NewRequest("GET", c.BaseApiURL+path+queryString, nil)
	if err != nil {
		panic(err)
	}
	return http_call(req)
}

func (c *Client) PostCall(path string, params *url.Values) ([]byte, error) {
	params.Set("accessToken", c.AccessToken)
	req, err := http.NewRequest("POST", c.BaseApiURL+path, strings.NewReader(params.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return http_call(req)
}

func (c *Client) UploadCall(path string,
	params *url.Values,
	uploadFieldName string,
	file string,
	reader io.Reader) ([]byte, error) {
	params.Set("accessToken", c.AccessToken)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	for k, _ := range *params {
		w.WriteField(k, params.Get(k))
	}
	wr, _ := w.CreateFormFile(uploadFieldName, filepath.Base(file))
	io.Copy(wr, reader)
	w.Close()

	req, err := http.NewRequest("POST", c.BaseApiURL+path, buf)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	return http_call(req)
}
