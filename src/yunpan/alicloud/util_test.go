package alicloud

import (
	"net/http"
	"testing"
)

func TestHttpCall_Get(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://www.baidu.com", nil)
	bytes, err := http_call(req)
	if err != nil {
		t.Log("output:" + string(bytes))
		t.Error(err)
	}
}

func TestHttpCall_with_error(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://api.yunpan.alibaba.com/api/upload/commit", nil)
	_, err := http_call(req)
	if err == nil {
		t.Error("Error is expected")
	} else {
		t.Log("Got expected error:", err.(ApiError).ErrorDescription)
	}
}
