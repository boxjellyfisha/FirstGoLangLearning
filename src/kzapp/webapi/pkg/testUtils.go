package pkg

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

// 測試輔助函數：創建測試請求
func CreateTestRequest(method, path, body string) *http.Request {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

// 測試輔助函數：執行請求並返回響應
func ExecuteRequest(handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

// 輔助函數：創建字串指標
func StringPtr(s string) *string {
	return &s
}
