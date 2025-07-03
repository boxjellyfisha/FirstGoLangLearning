package user

import (
	"kzapp/webapi/pkg"
	
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// 測試 login 函數
func TestLogin(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    UserResponse
		expectedStatus int
		shouldSucceed  bool
	}{
		{
			name:           "正常登入",
			requestBody:    UserResponse{Name: "testuser"},
			expectedStatus: http.StatusOK,
			shouldSucceed:  true,
		},
		{
			name:           "空名稱",
			requestBody:    UserResponse{Name: ""},
			expectedStatus: http.StatusBadRequest,
			shouldSucceed:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := pkg.CreateTestRequest("POST", "/login", string(jsonBody))
			rr := pkg.ExecuteRequest(login, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.shouldSucceed {
				var response UserResponse
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Errorf("無法解析響應 JSON: %v", err)
				}

				if response.Name != tt.requestBody.Name {
					t.Errorf("handler returned wrong name: got %v want %v", response.Name, tt.requestBody.Name)
				}

				if response.Message == nil || *response.Message != "Hello, Web API with Go!" {
					t.Errorf("handler returned wrong message: got %v want %v", response.Message, "Hello, Web API with Go!")
				}

				// 檢查是否設置了 cookie
				if cookies := rr.Result().Cookies(); len(cookies) == 0 {
					t.Errorf("沒有設置 session cookie")
				}
			}

			doLoggingout()
		})
	}
}

// 測試 login 函數的錯誤情況
func TestLoginInvalidJSON(t *testing.T) {
	req := pkg.CreateTestRequest("POST", "/login", `{"invalid": "json"`)
	rr := pkg.ExecuteRequest(login, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

// 測試 logout 函數
func TestLogout(t *testing.T) {
	rr := doLoggingout()

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// 檢查是否設置了 cookie
	if cookies := rr.Result().Cookies(); len(cookies) == 0 {
		t.Errorf("沒有設置 session cookie")
	}
}

func doLoggingout() *httptest.ResponseRecorder {
	req := pkg.CreateTestRequest("POST", "/logout", "")
	rr := pkg.ExecuteRequest(logout, req)
	return rr
}

// 測試 userSecretData 函數 - 未認證情況
func TestUserSecretDataUnauthenticated(t *testing.T) {
	req := pkg.CreateTestRequest("GET", "/info", "")
	rr := pkg.ExecuteRequest(userSecretData, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
	}
}

// 測試 userSecretData 函數 - 認證情況
func TestUserSecretDataAuthenticated(t *testing.T) {
	// 先登入
	loginReq := pkg.CreateTestRequest("POST", "/login", `{"name":"testuser"}`)
	loginRR := pkg.ExecuteRequest(login, loginReq)

	// 創建帶有 session cookie 的請求
	req := pkg.CreateTestRequest("GET", "/info", "")
	for _, cookie := range loginRR.Result().Cookies() {
		req.AddCookie(cookie)
	}

	rr := pkg.ExecuteRequest(userSecretData, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "The cake is a lie!") {
		t.Errorf("響應中沒有包含預期的秘密訊息")
	}

	doLoggingout()
}
