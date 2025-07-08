package webapi

import (
	"testing"
	"time"
)

func TestShutdownManager(t *testing.T) {
	// 測試數據庫初始化
	userDao, err := ProvideUserDao()
	if err != nil {
		t.Fatalf("數據庫初始化失敗: %v", err)
	}

	// 測試數據庫操作
	if userDao == nil {
		t.Fatal("UserDao 不應為 nil")
	}

	// 測試關閉功能
	err = Shutdown()
	if err != nil {
		t.Fatalf("關閉數據庫時發生錯誤: %v", err)
	}
}

func TestGracefulShutdown(t *testing.T) {
	// 這個測試模擬優雅關閉的流程
	// 在實際環境中，這需要信號處理，所以這裡只是驗證基本功能

	// 初始化數據庫
	userDao, err := ProvideUserDao()
	if err != nil {
		t.Fatalf("數據庫初始化失敗: %v", err)
	}

	// 驗證數據庫連接正常
	if userDao == nil {
		t.Fatal("UserDao 不應為 nil")
	}

	// 模擬關閉流程
	go func() {
		time.Sleep(100 * time.Millisecond)
		err := Shutdown()
		if err != nil {
			t.Errorf("關閉時發生錯誤: %v", err)
		}
	}()

	// 等待關閉完成
	time.Sleep(200 * time.Millisecond)
}
