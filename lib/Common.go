package lib

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

///生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

///生成Guid字串
func GetGuid() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

///字符串编码
func EncodeStr(src string) (dst string) {
	data, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), simplifiedchinese.GBK.NewEncoder()))
	if err == nil {
		dst = string(data)
	}
	return
}

///字符串解码
func DecodeStr(src string) (dst string) {
	data, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), simplifiedchinese.GBK.NewDecoder()))
	if err == nil {
		dst = string(data)
	}
	return
}

///获取url参数
func GetUrlValue(req *http.Request) map[string]string {
	values := make(map[string]string)
	req.ParseForm()
	if len(req.Form) > 0 {
		for k, v := range req.Form {
			values[k] = v[0]
		}
	}

	return values
}


//实现 try catch
func Try(fun func(), handler func(interface{})) {
    defer func() {
        if err := recover(); err != nil {
            handler(err)
        }
    }()
    fun()
}

///demo
// func testTry() {
//     Try(func() {
//        panic("foo")
//     }, func(e interface{}) {
//        print(e)
//     })
// }
