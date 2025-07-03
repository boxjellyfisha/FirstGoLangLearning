package count

import (
	"kzapp/webapi/pkg"
	
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 測試 add 函數
func TestAdd(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    Plusable
		expectedSum    int
		expectedStatus int
	}{
		{
			name:           "正常加法",
			requestBody:    Plusable{Num1: 3, Num2: 5},
			expectedSum:    8,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "負數加法",
			requestBody:    Plusable{Num1: -2, Num2: 3},
			expectedSum:    1,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "零值加法",
			requestBody:    Plusable{Num1: 0, Num2: 0},
			expectedSum:    0,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := pkg.CreateTestRequest("POST", "/add", string(jsonBody))
			rr := pkg.ExecuteRequest(add, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			var response Sum
			if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
				t.Errorf("無法解析響應 JSON: %v", err)
			}

			if response.Total != tt.expectedSum {
				t.Errorf("handler returned wrong sum: got %v want %v", response.Total, tt.expectedSum)
			}
		})
	}
}

// 測試 add 函數的錯誤情況
func TestAddInvalidJSON(t *testing.T) {
	req := pkg.CreateTestRequest("POST", "/add", `{"invalid": "json"`)
	rr := pkg.ExecuteRequest(add, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

// 測試 square 函數
func TestSquare(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedSquare int
		expectedStatus int
	}{
		{
			name:           "正常平方",
			path:           "/square/9",
			expectedSquare: 81,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "零的平方",
			path:           "/square/0",
			expectedSquare: 0,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "負數平方",
			path:           "/square/-3",
			expectedSquare: 9,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := pkg.CreateTestRequest("GET", tt.path, "")
			rr := pkg.ExecuteRequest(square, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			var response Sum
			if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
				t.Errorf("無法解析響應 JSON: %v", err)
			}

			if response.Total != tt.expectedSquare {
				t.Errorf("handler returned wrong square: got %v want %v", response.Total, tt.expectedSquare)
			}
		})
	}
}

// 測試 square 函數的錯誤情況
func TestSquareInvalidNumber(t *testing.T) {
	req := pkg.CreateTestRequest("GET", "/square/abc", "")
	rr := pkg.ExecuteRequest(square, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func BenchmarkAdd(b *testing.B) {
	jsonBody, _ := json.Marshal(Plusable{Num1: 3, Num2: 5})
	req := pkg.CreateTestRequest("POST", "/add", string(jsonBody))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		add(rr, req)
	}
}

func BenchmarkSquare(b *testing.B) {
	req := pkg.CreateTestRequest("GET", "/square/9", "")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		square(rr, req)
	}
}
