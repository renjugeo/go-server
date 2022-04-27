package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHealthCheck(t *testing.T) {
	srv := handler{}
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	srv.handleHealthCheck()(w, req)
	resp := w.Result()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("error reading response body %v", err)
	}
	if string(respData) != "OK" {
		t.Errorf("expected 'OK', received %v", string(respData))
	}
}
