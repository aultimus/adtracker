package adtracker

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cocoonlife/testify/assert"
)

func doRequest(handler http.HandlerFunc, r *http.Request) (int, []byte) {
	w := httptest.NewRecorder()
	handler(w, r)
	resp := w.Result()

	var body []byte
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		body, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}
	return resp.StatusCode, body
}

func parseResponse(b []byte) (int, error) {
	resp := AdCountResp{}
	err := json.Unmarshal(b, &resp)
	return resp.Value, err
}

func TestHandlers(t *testing.T) {
	a := assert.New(t)

	at := AdTracker{NewBasicStore()}

	// GET on non-existing key results in 404
	b, err := json.Marshal(AdCountReq{"some_id"})
	a.NoError(err)
	req := httptest.NewRequest(http.MethodGet, "http://foo", bytes.NewReader(b))
	code, _ := doRequest(at.adCountHandler, req)
	a.Equal(http.StatusNotFound, code)

	// increment non-existent key is valid
	req = httptest.NewRequest(http.MethodPut, "http://foo", bytes.NewReader(b))
	code, _ = doRequest(at.trackHandler, req)
	a.Equal(http.StatusOK, code)

	// GET on incremented key returns 1
	req = httptest.NewRequest(http.MethodGet, "http://foo", bytes.NewReader(b))
	code, respBytes := doRequest(at.adCountHandler, req)
	a.Equal(http.StatusOK, code)
	v, err := parseResponse(respBytes)
	a.NoError(err)
	a.Equal(1, v)

	// increment existing key is valid
	req = httptest.NewRequest(http.MethodPut, "http://foo", bytes.NewReader(b))
	code, _ = doRequest(at.trackHandler, req)
	a.Equal(http.StatusOK, code)

	// GET on already incremented key returns 2
	req = httptest.NewRequest(http.MethodGet, "http://foo", bytes.NewReader(b))
	code, respBytes = doRequest(at.adCountHandler, req)
	a.Equal(http.StatusOK, code)
	v, err = parseResponse(respBytes)
	a.NoError(err)
	a.Equal(2, v)
}

func TestEmptyID(t *testing.T) {
	a := assert.New(t)

	at := AdTracker{NewBasicStore()}

	// empty id is a valid identifier - perhaps should not be
	b, err := json.Marshal(AdCountReq{""})
	a.NoError(err)
	req := httptest.NewRequest(http.MethodGet, "http://foo", bytes.NewReader(b))
	code, _ := doRequest(at.adCountHandler, req)
	a.Equal(http.StatusNotFound, code)
}

func TestInvalidReq(t *testing.T) {
	a := assert.New(t)

	at := AdTracker{NewBasicStore()}

	// empty id is a valid identifier - perhaps should not be
	b := []byte("{'foo': 'bar'}")
	req := httptest.NewRequest(http.MethodGet, "http://foo", bytes.NewReader(b))
	code, _ := doRequest(at.adCountHandler, req)
	a.Equal(http.StatusBadRequest, code)

	req = httptest.NewRequest(http.MethodPut, "http://foo", bytes.NewReader(b))
	code, _ = doRequest(at.trackHandler, req)
	a.Equal(http.StatusBadRequest, code)
}

func TestInvalidMethod(t *testing.T) {
	// TODO: Check which method types are invalid for each handler
}
