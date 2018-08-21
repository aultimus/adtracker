package adtracker

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cocoonlife/testify/assert"
)

func doRequest(
	handler func(w http.ResponseWriter, r *http.Request, vars map[string]string),
	r *http.Request, vars map[string]string) (int, []byte) {
	w := httptest.NewRecorder()
	handler(w, r, vars)
	resp := w.Result()

	var body []byte
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		body, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}
	return resp.StatusCode, body
}

func parseResponse(b []byte) (int, error) {
	resp := Resp{}
	err := json.Unmarshal(b, &resp)
	return resp.Value, err
}

func TestHandlers(t *testing.T) {
	a := assert.New(t)

	at := App{NewBasicStore()}
	vars := map[string]string{"id": "some_id"}

	// GET on non-existing key results in 404
	req := httptest.NewRequest(http.MethodGet, "http://foo", nil)
	code, _ := doRequest(at.getHandlerVars, req, vars)
	a.Equal(http.StatusNotFound, code)

	// increment non-existent key is valid
	req = httptest.NewRequest(http.MethodPut, "http://foo", nil)
	code, _ = doRequest(at.incHandlerVars, req, vars)
	a.Equal(http.StatusOK, code)

	// GET on incremented key returns 1
	req = httptest.NewRequest(http.MethodGet, "http://foo", nil)
	code, respBytes := doRequest(at.getHandlerVars, req, vars)
	a.Equal(http.StatusOK, code)
	v, err := parseResponse(respBytes)
	a.NoError(err)
	a.Equal(1, v)

	// increment existing key is valid
	req = httptest.NewRequest(http.MethodPut, "http://foo", nil)
	code, _ = doRequest(at.incHandlerVars, req, vars)
	a.Equal(http.StatusOK, code)

	// GET on already incremented key returns 2
	req = httptest.NewRequest(http.MethodGet, "http://foo", nil)
	code, respBytes = doRequest(at.getHandlerVars, req, vars)
	a.Equal(http.StatusOK, code)
	v, err = parseResponse(respBytes)
	a.NoError(err)
	a.Equal(2, v)
}

func TestInvalidMethod(t *testing.T) {
	// TODO: Check which method types are invalid for each handler
}
