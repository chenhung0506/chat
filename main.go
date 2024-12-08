package main

import (
	"chat/internal/controller"
	"chat/internal/redis"
	"chat/internal/websocket"
	"log"
	"net/http"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "172.17.0.1:6379"
	}

	elasticsearchAddr := os.Getenv("ELASTICSEARCH_ADDR")
	if elasticsearchAddr == "" {
		elasticsearchAddr = "http://139.162.2.175:9200"
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3002"
	}

	// 初始化 Redis 客戶端
	redisClient := redis.NewClient(redisAddr, "")

	// 註冊處理器
	http.HandleFunc("/websocket/", websocket.WebSocketHandler(redisClient)) // 注意結尾的斜線，允許匹配動態路徑

	// 註冊 /elk/autocomplete 並支持 CORS
	http.HandleFunc("/elk/autocomplete", func(w http.ResponseWriter, r *http.Request) {
		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // 注意：生产环境中替换为特定域名
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 处理 OPTIONS 请求
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 调用实际的处理逻辑
		controller.AutoCompleteHandler(elasticsearchAddr)(w, r)
	})

	log.Printf("Server running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
