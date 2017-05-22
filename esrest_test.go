package esrest

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"log"
	"os"
	"time"
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
		Query("Param1", "value").
		DoJson(h)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, "value", h.Args.Param1)
}

func TestPost(t *testing.T) {
	h := &H{}
	b := string(`{"message":"ok"}`)
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
	_, err := New().Get("http://httpbin.org/get").Do()
	assert.Equal(t, nil, err)
}

func TestValidUrlFail(t *testing.T) {
	_, err := New().Do()
	e := errors.New("url is empty")
	assert.Equal(t, err, e)
}

func TestValidLoggerFail(t *testing.T) {
	_, err := New().Get("dummy").Logger(nil).Do()
	e := errors.New("logger is empty")
	assert.Equal(t,e,err)
}

func TestDummyUrl(t *testing.T) {
	_, err := New().Get("http://123").Do()
	assert.Error(t, err)
}

func TestDoJsonSuccess(t *testing.T) {
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

func TestDoJsonFail(t *testing.T) {
	_, err := New().Get("dummy").DoJson(nil)
	assert.Error(t, err)
}

func TestSendByteSliceBody(t *testing.T) {
	h := &H{}
	b := []byte(`{"message":"ok"}`) //send byte slice body
	r, _ := New().
		Post("http://httpbin.org/post").
		Body(b).
		DoJson(h)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, string(b), h.Data)

}

func TestSendStructBody(t *testing.T) {
	h := &H{}
	b := struct {
		Message string `json:"message"`
	}{"ok"}

	r, _ := New().
		Post("http://httpbin.org/post").
		Body(b).
		DoJson(h)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, `{"message":"ok"}`, h.Data)

}

func TestDebugMode(t *testing.T) {
	h := &H{}
	b := struct {
		Message string `json:"message"`
	}{"ok"}

	r, _ := New().
		Debug(true).
		Header("h","v").
		Post("http://httpbin.org/post").
		Body(b).
		DoJson(h)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, `{"message":"ok"}`, h.Data)

}

func TestLogger(t *testing.T) {
	h := &H{}
	r, _ := New().
		Logger(log.New(os.Stdout, "", log.LstdFlags)).
		Get("http://httpbin.org/get").
		Query("Param1", "value").
		DoJson(h)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, "value", h.Args.Param1)

}

func TestTimeout(t *testing.T) {
	r, _ := New().
		Timeout(time.Duration(10 * time.Second)).
		Get("http://httpbin.org/get").
		Query("Param1", "value").
		Do()
	assert.Equal(t, 200, r.StatusCode)


}
