package lib

import (
	"io/ioutil"
	"net/http"
	"strings"
    
    "crypto/md5"
    "crypto/rand"
    "encoding/base64"
    "encoding/hex"
    "io"
)

///向nsq服务器推送一条消息
func HttpDo_NSQ(faction, furl, fdata string) string {
	client := &http.Client{}
	req, err := http.NewRequest(faction, furl, strings.NewReader(fdata))
	if err != nil {
		return "接口错误"
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "发送失败"
	}
	return string(body)
}


//生成32位md5字串
func GetMd5String(s string) string {
    h := md5.New()
    h.Write([]byte(s))
    return hex.EncodeToString(h.Sum(nil))
}

//生成Guid字串
func GetGuid() string {
    b := make([]byte, 48)

    if _, err := io.ReadFull(rand.Reader, b); err != nil {
        return ""
    }
    return GetMd5String(base64.URLEncoding.EncodeToString(b))
}
