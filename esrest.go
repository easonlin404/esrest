package esrest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"reflect"
	"time"
)

type Builder struct {
	Url       string
	Method    string
	Path      string
	Headers   map[string]string
	Querys    map[string]string
	DebugMode bool

	logger    *log.Logger
	timeout   time.Duration
	basicAuth auth
	bodyByte  []byte
}


type auth  struct{ username, password string }

const DefaultContentType = "application/json"

func New() *Builder {
	return &Builder{
		Headers: make(map[string]string),
		Querys:  make(map[string]string),
		logger:  log.New(os.Stdout, "", log.LstdFlags),
		timeout: time.Duration(20 * time.Second),
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

func (b *Builder) Head(url string) *Builder {
	b.Url = url
	b.Method = http.MethodHead
	return b
}

func (b *Builder) Header(key, value string) *Builder {
	b.Headers[key] = value
	return b
}

func (b *Builder) Query(key, value string) *Builder {
	b.Querys[key] = value
	return b
}

func (b *Builder) Body(v interface{}) *Builder {
	rv := reflect.ValueOf(v)
	//fmt.Printf("%+v\n",rv)
	//fmt.Println(rv.Kind())

	switch rv.Kind() {
	case reflect.String:
		b.bodyByte = []byte(rv.String())
	case reflect.Slice:
		slice, _ := rv.Interface().([]byte)
		b.bodyByte = slice
	case reflect.Struct, reflect.Ptr:
		byte, _ := json.Marshal(v)
		b.bodyByte = byte
	}
	return b
}

func (b *Builder) Do() (*http.Response, error) {
	if err := b.valid(); err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: b.timeout}

	request := b.newRequest()

	if b.DebugMode {
		dump, _ := httputil.DumpRequest(request, true)
		b.logger.Println(string(dump))
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if b.DebugMode {
		dump, _ := httputil.DumpResponse(resp, true)
		b.logger.Println(string(dump))
	}

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

	//Set Default Content-Type Header
	if len(req.Header.Get("Content-Type")) == 0 {
		req.Header.Set("Content-Type", DefaultContentType)
	}

	//Set Header
	for k, v := range b.Headers {
		req.Header.Set(k, v)
	}

	//Set Query
	q := req.URL.Query()
	for k, v := range b.Querys {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	if b.basicAuth != (auth{}) {
		req.SetBasicAuth(b.basicAuth.username, b.basicAuth.password)
	}

	return req
}

func (b *Builder) DoJson(v interface{}) (*http.Response, error) {
	resp, err := b.Do()
	if err != nil {
		return resp, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	json.Unmarshal(body, v)

	return resp, nil
}

func (b *Builder) valid() error {
	if len(b.Url) == 0 {
		return errors.New("url is empty")
	}
	if b.logger == nil {
		return errors.New("logger is empty")
	}
	return nil
}

func (b *Builder) Debug(debug bool) *Builder {
	b.DebugMode = debug
	return b
}

func (b *Builder) Logger(log *log.Logger) *Builder {
	b.logger = log
	return b
}

func (b *Builder) Timeout(timeout time.Duration) *Builder {
	b.timeout = timeout
	return b
}

func (b *Builder) BasicAuth(username, password string) *Builder {
	b.basicAuth = auth{username:username,password:password}
	return b
}
