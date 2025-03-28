package main

import (
	"chat/internal/controller"
	"chat/internal/redis"
	"chat/internal/service"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	go func() {
		if err := http.ListenAndServe("0.0.0.0:6060", nil); err != nil {
			log.Fatalf("pprof 啟動失敗: %v", err)
		}
	}()

	log.Println("啟動 Receiver 1")

	redisAddr := service.GetEnvWithDefault("REDIS_ADDR", "localhost:6379")
	elasticsearchAddr := service.GetEnvWithDefault("ELASTICSEARCH_ADDR", "http://139.162.2.175:9200")
	port := service.GetEnvWithDefault("SERVER_PORT", "3002")

	redisClient := redis.NewClient(redisAddr, "")

	go service.SetupExchange()
	go service.ReceiveMessage("receiver_queue", redisClient)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/websocket/client/:client_id", controller.HandleClient(redisClient))
	r.GET("/websocket/service/:service_id", controller.HandleService(redisClient))
	r.GET("/waiting-clients", func(c *gin.Context) { controller.GetWaitingClients(c, redisClient) })
	r.GET("/waiting-clients-clear", func(c *gin.Context) { controller.RemoveWholeWaitingQueue(c, redisClient) })
	r.POST("/assign-client", func(c *gin.Context) { controller.AssignClient(c, redisClient) })

	r.POST("/elk/autocomplete", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		controller.AutoCompleteHandler(elasticsearchAddr)(c.Writer, c.Request)
	})

	r.GET("/elk/autocomplete", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		controller.AutoCompleteHandler(elasticsearchAddr)(c.Writer, c.Request)
	})

	log.Printf("Server running on port %s...", port)
	r.Run(":" + port)
}
