package api

import (
	"encoding/json"
	"io"
	"net/http"
)

func HttpGet(url string, headers *http.Header) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header = *headers
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetJson(url string, result interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return err
	}
	return nil
}

func GetJsonWithHeader(url string, headers *http.Header, result interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	//fmt.Println(headers)
	headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0")
	req.Header = *headers
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return err
	}
	return nil
}
