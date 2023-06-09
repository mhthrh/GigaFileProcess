package View

import (
	"github.com/gin-gonic/gin"
	"github.com/mhthrh/GigaFileProcess/api/Server"
	"net/http"
)

func RunSync() http.Handler {
	router := gin.New()
	router.Use(gin.Recovery())

	router.POST("/gateway/file", Server.Run)
	router.POST("/gateway/investigate", Server.Run)
	router.GET("/version", Server.Version)
	router.NoRoute(Server.NotFound)

	return router
}
