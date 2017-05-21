package esrest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"io"
)

type Builder struct {
	Url     string
	Method  string
	Path    string
	Headers map[string]string

	bodyByte []byte
}

const DefaultContentType = "application/json"

func New() *Builder {
	return &Builder{
		Headers: make(map[string]string),
	}
}

func (b *Builder) Get(url string) *Builder {
	b.Url = url
	b.Method = http.MethodGet
	return b
}

func (b *Builder) Post(url string) *Builder {
	b.Url = url
	b.Method = http.MethodPost
	return b
}

func (b *Builder) Put(url string) *Builder {
	b.Url = url
	b.Method = http.MethodPut
	return b
}

func (b *Builder) Delete(url string) *Builder {
	b.Url = url
	b.Method = http.MethodDelete
	return b
}

func (b *Builder) Header(key, value string) *Builder {
	b.Headers[key] = value
	return b
}

func (b *Builder) Body(body []byte) *Builder {
	b.bodyByte = body
	return b
}

func (b *Builder) Do() (*http.Response, error) {
	if err := b.valid(); err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: time.Duration(20 * time.Second)}
	resp, err := client.Do(b.newRequest())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return resp, nil
}

func (b *Builder) newRequest() *http.Request {
	var reader io.Reader
	if len(b.bodyByte) > 0 {
		reader = bytes.NewBuffer(b.bodyByte)
	}

	req, _ := http.NewRequest(b.Method, b.Url, reader)

	fmt.Println(req.Header.Get("Content-Type"))
	if len(req.Header.Get("Content-Type")) == 0 {
		req.Header.Set("Content-Type", DefaultContentType)
	}

	for k, v := range b.Headers {
		req.Header.Set(k, v)
	}

	return req
}

func (b *Builder) AsJson(v interface{}) (*http.Response, error) {
	resp, err := b.Do()
	if err != nil {
		return resp, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &v)
	return resp, nil
}

func (b *Builder) valid() error {
	if len(b.Url) == 0 {
		return errors.New("url is empty")
	}
	return nil
}
