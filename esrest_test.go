package esrest

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	n1 := New()
	n2 := New()

	assert.Equal(t, n1.Headers != nil, true)
	assert.Equal(t, n1, n2)
	assert.Equal(t, n1 != nil, true)
}

func TestGet(t *testing.T) {
	b := New().Get("https://www.google.com")
	b.Body([]byte("message=ok"))
	assert.Equal(t, b.Method, http.MethodGet)
}

func TestPost(t *testing.T) {
	b := New().Post("https://www.google.com")
	assert.Equal(t, b.Method, http.MethodPost)
	b.Body([]byte("hi"))
	_,err:=b.Do()

	assert.Equal(t,nil,err)
}

func TestSetPut(t *testing.T) {
	b := New().Put("https://www.google.com")
	assert.Equal(t, b.Method, http.MethodPut)
}

func TestSetDelete(t *testing.T) {
	b := New().Delete("https://www.google.com")
	assert.Equal(t, b.Method, http.MethodDelete)
}

func TestSetHeader(t *testing.T) {
	b := New().Header("key", "value")
	value := b.Headers["key"]
	assert.Equal(t, "value", value)

	b.Get("https://www.google.com").Do()

}

func TestValidOk(t *testing.T) {
	_, err := New().Get("https://www.google.com").Do()
	assert.Equal(t, nil, err)
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

	New().Get(server.URL).AsJson(json)
	assert.Equal(t, json.Message, "hi")
}

func TestAsJsonFail(t *testing.T) {
	_,err := New().Get("dummy").AsJson(nil)
	assert.Error(t,err)
}