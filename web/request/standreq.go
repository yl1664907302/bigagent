package request

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PostStand struct {
	p Post
}

func (p *PostStand) do() {
	// 指定请求的URL
	url := "http://example.com/post"

	// 创建请求体
	payload := bytes.NewBufferString("key1=value1&key2=value2")

	// 创建一个HTTP客户端
	client := &http.Client{}

	// 发送POST请求
	resp, err := client.Post(url, "application/x-www-form-urlencoded", payload)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close() // 确保在函数返回时关闭响应体

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 打印响应内容
	fmt.Println(string(body))
}
