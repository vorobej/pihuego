package request

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// POST request to bridge IP
func POST(url string, data io.Reader) ([]byte, error) {
	fmt.Printf("POST <%s>\n", url)
	resp, err := http.Post("http://"+url, "application/json", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp == nil {
		return nil, fmt.Errorf("response is nil")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code is not OK")
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

// PUT Sends request to bridge IP
func PUT(url string, data []byte) ([]byte, error) {
	fmt.Printf("PUT <%s> data<%s>\n", url, string(data))
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, "http://"+url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code is not OK")
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

// GET send get request
func GET(url string) ([]byte, error) {
	fmt.Printf("GET <%s>\n", url)
	resp, err := http.Get("http://" + url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code is not OK")
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}
