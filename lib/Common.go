package lib

import (
	"io/ioutil"
	"net/http"
	"strings"
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
