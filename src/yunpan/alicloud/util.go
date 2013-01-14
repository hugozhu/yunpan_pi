package alicloud

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func GetAccessToken() string {
	bytes, err := ioutil.ReadFile(filepath.Join(os.Getenv("PWD"), "token"))
	if err != nil || len(bytes) < 32 {
		panic(filepath.Join(os.Getenv("PWD"), "token") + " is not valid OAuth access token")
	}
	token := string(bytes)[0:32]
	return token
}

var ERROR_PATTERN = []byte("error_description")

func http_call(req *http.Request) ([]byte, error) {
	if req.Method == "POST" {
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	output, err2 := ioutil.ReadAll(resp.Body)
	if output != nil {
		if bytes.Contains(output, ERROR_PATTERN) {
			return nil, NewApiError(output)
		}
	}
	return output, err2
}

func debug(v ...interface{}) {
	log.Println("[debug]", v)
}

func md5_hash(filepath string) string {
	return hex.EncodeToString(_md5(filepath))
}

func md5_bytes(bytes []byte, n int) string {
	h := md5.New()
	h.Write(bytes[0:n])
	return hex.EncodeToString(h.Sum(nil))
}

func checksum_bytes(bytes []byte, n int) string {
	a := 0
	b := 0
	contentSize := n
	sp := 0
	ep := n
	for {
		if sp >= ep {
			break
		}
		a = (a + int(bytes[sp])) & 0xffff
		b = (b + contentSize*int(bytes[sp])) & 0xffff
		contentSize--
		sp++
	}
	return fmt.Sprintf("%d", b<<16|a)
}

func _md5(filepath string) []byte {
	h := md5.New()
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		h.Write(buf[0:n])
	}
	return h.Sum(nil)
}
