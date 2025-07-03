# firstServer API 單元測試指南

## 📋 測試概述

這個測試套件為 `firstServer.go` 中的所有 API 端點提供了完整的單元測試覆蓋。

## 🚀 運行測試

### 方法 1：使用測試腳本

```bash
./run_tests.sh
```

### 方法 2：直接使用 go test

```bash
# 執行所有測試
go test -v ./webapi/*/

# 執行特定測試
go test -v -run TestGreet ./webapi/*/

# 執行基準測試
go test -bench=. ./webapi/*/

# 查看測試覆蓋率
go test -cover ./webapi/*/

# 生成覆蓋率報告
go test -coverprofile=coverage.out ./webapi/*/
go tool cover -html=coverage.out
```

## 📊 測試覆蓋範圍

### API 端點測試

- ✅ `/greet` - 問候 API
- ✅ `/add` - 加法 API
- ✅ `/square/{num}` - 平方 API
- ✅ `/login` - 登入 API
- ✅ `/logout` - 登出 API
- ✅ `/info` - 用戶資訊 API

### 中間件測試

- ✅ `Logging()` - 日誌中間件
- ✅ `Method()` - HTTP 方法驗證中間件
- ✅ `Chain()` - 中間件鏈接函數

### 輔助函數測試

- ✅ `encodeResponse()` - JSON 響應編碼

## 🧪 測試類型

### 1. 單元測試

- **正常情況測試**：驗證 API 在正確輸入下的行為
- **錯誤情況測試**：驗證 API 在錯誤輸入下的錯誤處理
- **邊界條件測試**：測試極端值和邊界情況

### 2. 基準測試

- 測量 API 的效能表現
- 幫助識別效能瓶頸

### 3. 整合測試

- 測試中間件與處理函數的整合
- 測試路由設置

## 📝 測試案例說明

### TestGreet

- 測試帶查詢參數的請求
- 測試空查詢參數的請求
- 驗證響應格式和內容

### TestAdd

- 測試正常數字加法
- 測試負數加法
- 測試零值加法
- 測試無效 JSON 輸入

### TestSquare

- 測試正常數字平方
- 測試零的平方
- 測試負數平方
- 測試無效數字輸入

### TestLogin/TestLogout

- 測試正常登入流程
- 測試空名稱登入
- 測試無效 JSON 登入
- 測試登出功能
- 驗證 session cookie 設置

### TestUserSecretData

- 測試未認證用戶訪問
- 測試已認證用戶訪問
- 驗證認證狀態檢查

## 🔧 測試輔助函數

### pkg.CreateTestRequest()

創建標準化的測試 HTTP 請求

### pkg.ExecuteRequest()

執行 HTTP 處理函數並返回響應記錄器

### stringPtr()

創建字串指標的輔助函數

## 📈 測試最佳實踐

1. **測試隔離**：每個測試都是獨立的，不依賴其他測試
2. **清晰命名**：測試函數名稱清楚描述測試內容
3. **完整覆蓋**：包含正常情況、錯誤情況和邊界條件
4. **可讀性**：使用表格驅動測試提高可讀性
5. **效能測試**：包含基準測試監控效能

## 🐛 除錯技巧

### 查看詳細測試輸出

```bash
go test -v ./webapi
```

### 運行特定測試

```bash
go test -v -run TestGreet ./webapi
```

### 查看測試覆蓋率

```bash
go test -coverprofile=coverage.out ./webapi
go tool cover -html=coverage.out -o coverage.html
```

### 並行測試

```bash
go test -parallel 4 ./webapi
```

## 📚 相關資源

- [Go Testing Package](https://golang.org/pkg/testing/)
- [httptest Package](https://golang.org/pkg/net/http/httptest/)
- [Gorilla Mux Testing](https://github.com/gorilla/mux#testing-handlers)
- [Go Test Coverage](https://blog.golang.org/cover)
