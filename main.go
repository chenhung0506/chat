package main

import (
	"chat/internal/controller"
	"chat/internal/redis"
	"log"
	"os"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	redisAddr := getEnvWithDefault("REDIS_ADDR", "localhost:6379")
	elasticsearchAddr := getEnvWithDefault("ELASTICSEARCH_ADDR", "http://139.162.2.175:9200")
	port := getEnvWithDefault("SERVER_PORT", "3002")

	redisClient := redis.NewClient(redisAddr, "")

	r := gin.Default()

	// 啟用 CORS
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
