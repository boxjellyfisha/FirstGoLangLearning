package webapi

import (
	"log"
	"net/http"
	"time"

	"kzapp/webapi/pkg"

	"github.com/gorilla/mux"
)

func RunServer() {
	// Create the router
	router := mux.NewRouter()

	handlers := CreateHandler()

	// Handle incoming request
	for _, handler := range handlers {
		handler.InitService(router)
	}

	// 創建HTTP服務器
	srv := &http.Server{
		Addr:    ":80",
		Handler: router,
	}

	// 創建關閉管理器
	shutdownManager := pkg.NewShutdownManager()

	// 添加數據庫關閉鉤子
	shutdownManager.AddShutdownHook(func() error {
		log.Println("正在關閉數據庫連接...")
		return Shutdown()
	})

	// 在goroutine中啟動服務器
	go func() {
		log.Printf("服務器正在啟動，監聽端口 80...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服務器啟動失敗: %v", err)
		}
	}()

	// 使用關閉管理器處理優雅關閉
	pkg.GracefulShutdown(srv, 5*time.Second)
}
