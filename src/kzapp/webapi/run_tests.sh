#!/bin/bash

echo "🚀 開始執行 firstServer API 單元測試..."
echo "=========================================="

# 執行所有測試
echo "📋 執行所有測試..."
go test -v ./webapi/*/

echo ""
echo "📊 執行基準測試..."
go test -bench=. ./webapi/*/

echo ""
echo "📈 執行測試覆蓋率..."
go test -cover ./webapi/*/

echo ""
echo "✅ 測試完成！" 