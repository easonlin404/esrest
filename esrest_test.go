package esrest

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
	"io/ioutil"
)

type H struct {
	Data string `json:"data"`
	Args struct {
		Param1 string `json:"param1"`
	}
	Headers struct {
		TestHader string `json:"testHader"`
	}
}

func TestGet(t *testing.T) {
	h := &H{}
	r, _ := New().
		Get("http://httpbin.org/get").
		Query("Param1","value").
		DoJson(h)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, "value", h.Args.Param1)
	fmt.Printf("%+v", h)
}

func TestPost(t *testing.T) {
	h := &H{}
	b:= string(`{"message":"ok"}`)
	r, _ := New().
		Post("http://httpbin.org/post").
		Body(b).
		DoJson(h)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, b, h.Data)

}

func TestPut(t *testing.T) {
	r, _ := New().
		Put("http://httpbin.org/put").
		Do()
	assert.Equal(t, 200, r.StatusCode)
}

func TestDelete(t *testing.T) {
	r, _ := New().
		Delete("http://httpbin.org/delete").
		Do()
	assert.Equal(t, 200, r.StatusCode)
}

func TestHeader(t *testing.T) {
	h := &H{}
	r, _ := New().
		Get("http://httpbin.org/headers").
		Header("TestHader", "value").
		DoJson(h)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, "value", h.Headers.TestHader)
}

func TestValidOk(t *testing.T) {
	r, err := New().Get("http://httpbin.org/get").Do()
	assert.Equal(t, nil, err)
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(body))
}

func TestValidFail(t *testing.T) {
	_, err := New().Do()
	e := errors.New("url is empty")
	assert.Equal(t, err, e)
}

func TestDummyUrl(t *testing.T) {
	_, err := New().Get("http://123").Do()
	assert.Error(t, err)
}

func TestAsJsonSuccess(t *testing.T) {
	type Json struct {
		Message string `json:"message"`
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json := []byte(`{"message":"hi"}`)
		w.Write(json)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	json := &Json{}

	New().Get(server.URL).DoJson(json)
	assert.Equal(t, json.Message, "hi")
}

func TestAsJsonFail(t *testing.T) {
	_, err := New().Get("dummy").DoJson(nil)
	assert.Error(t, err)
}
