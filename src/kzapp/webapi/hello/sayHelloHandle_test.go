package hello

import (
	"kzapp/webapi/pkg"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func runTest(task func(target GreetingHandler)) {
	target := GreetingHandler{}
	task(target)
}

// 測試 greet 函數
func TestGreet(t *testing.T) { runTest(func(target GreetingHandler) {
		tests := []struct {
			name           string
			queryParam     string
			expectedName   string
			expectedStatus int
		}{
			{
				name:           "正常查詢參數",
				queryParam:     "?name=Alice",
				expectedName:   "Alice",
				expectedStatus: http.StatusOK,
			},
			{
				name:           "空查詢參數",
				queryParam:     "",
				expectedName:   "",
				expectedStatus: http.StatusOK,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := pkg.CreateTestRequest("GET", "/greet"+tt.queryParam, "")
				rr := pkg.ExecuteRequest(target.greet, req)

				if status := rr.Code; status != tt.expectedStatus {
					t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
				}

				var response GreetingResponse
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Errorf("無法解析響應 JSON: %v", err)
				}

				if response.Name != tt.expectedName {
					t.Errorf("handler returned wrong name: got %v want %v", response.Name, tt.expectedName)
				}

				if response.Message == nil || *response.Message != "Welcome Back!" {
					t.Errorf("handler returned wrong message: got %v want %v", response.Message, "Welcome Back!")
				}
			})
		}
	})
}

// 基準測試
func BenchmarkGreet(b *testing.B) { runTest(func(target GreetingHandler){
		req := pkg.CreateTestRequest("GET", "/greet?name=test", "")

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			rr := httptest.NewRecorder()
			target.greet(rr, req)
		}
	})
}
