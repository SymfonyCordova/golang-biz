package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
)

type HttpDClient struct {
}

func NewHttpClient() *HttpDClient {
	return new(HttpDClient)
}

func (hc *HttpDClient) PostJson(url string, byteJson []byte) ([]byte, error) {
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(byteJson))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		//Timeout: 3 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (hc *HttpDClient) Post(url string, params string) ([]byte, error) {
	contentType := "application/x-www-form-urlencoded"

	data := strings.NewReader(params)

	response, err := http.Post(url, contentType, data)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	//读取网页内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (hc *HttpDClient) Get(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

//golang 官方提供的 http 包里的 http client 可以通过一下两种方法设置超时
//设置1s超时
//cli := http.Client{Timeout: time.Second}

// 设置1s超时
//req := http.NewRequest(....)
//ctx, _  := context.WithTimeout(time.Second)
//req.WithContext(ctx)

//go httpclient 如何判断请求超时
//uri := "http://www.google.com/"
//client := http.Client{
//Timeout: 3 * time.Second,
//}

//resp, err := client.Get(uri)
//if err != nil {
//if strings.Contains(err.Error(), "Client.Timeout exceeded") {
//	fmt.Println("HTTP Get timeout")
//}
//}
//defer resp.Body.Close()

//req, err := http.NewRequest("GET", url, nil)
//if err != nil {
//// TODO:
//}
//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//defer cancel()
//isTimeout := func() bool {
//	select {
//	default:
//		return false
//	case <-ctx.Done():
//		// 是否超时
//		return errors.Is(ctx.Err(), context.DeadlineExceeded)
//	}
//}
//req = req.WithContext(ctx)
//resp, err := http.DefaultClient.Do(req)
//if err != nil {
//if isTimeout() {
//// TODO: 超时处理
//}
//return
//}
//defer resp.Body.Close()
