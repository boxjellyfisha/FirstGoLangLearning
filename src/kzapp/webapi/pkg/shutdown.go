package pkg

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ShutdownManager 管理應用程序的優雅關閉
type ShutdownManager struct {
	shutdownHooks []func() error
	mu            sync.Mutex
}

// NewShutdownManager 創建新的關閉管理器
func NewShutdownManager() *ShutdownManager {
	return &ShutdownManager{
		shutdownHooks: make([]func() error, 0),
	}
}

// AddShutdownHook 添加關閉鉤子
func (sm *ShutdownManager) AddShutdownHook(hook func() error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.shutdownHooks = append(sm.shutdownHooks, hook)
}

// WaitForShutdown 等待關閉信號並執行優雅關閉
func (sm *ShutdownManager) WaitForShutdown() {
	// 創建信號通道
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 等待信號
	<-quit
	log.Println("收到關閉信號，開始優雅關閉...")

	// 執行關閉鉤子
	sm.executeShutdownHooks()
}

// executeShutdownHooks 執行所有關閉鉤子
func (sm *ShutdownManager) executeShutdownHooks() {
	sm.mu.Lock()
	hooks := make([]func() error, len(sm.shutdownHooks))
	copy(hooks, sm.shutdownHooks)
	sm.mu.Unlock()

	// 並行執行所有鉤子
	var wg sync.WaitGroup
	for i, hook := range hooks {
		wg.Add(1)
		go func(index int, h func() error) {
			defer wg.Done()
			if err := h(); err != nil {
				log.Printf("關閉鉤子 %d 執行失敗: %v", index, err)
			} else {
				log.Printf("關閉鉤子 %d 執行成功", index)
			}
		}(i, hook)
	}

	// 等待所有鉤子完成，最多等待10秒
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("所有關閉鉤子執行完成")
	case <-time.After(10 * time.Second):
		log.Println("關閉超時，強制退出")
	}
}

// GracefulShutdown 提供一個便捷的函數來設置優雅關閉
func GracefulShutdown(server *http.Server, shutdownTimeout time.Duration) {
	sm := NewShutdownManager()

	// 添加服務器關閉鉤子
	sm.AddShutdownHook(func() error {
		log.Println("正在關閉HTTP服務器...")
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return server.Shutdown(ctx)
	})

	// 等待關閉信號
	sm.WaitForShutdown()
}
