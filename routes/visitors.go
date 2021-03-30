package routes

import (
	"github.com/basselalaraaj/unique-visitors/handlers"
	"github.com/gin-gonic/gin"
)

func AddVisitorsRoutes(rg *gin.Engine) {
	visitors := rg.Group("/visitors")
	visitors.GET("", handlers.VisitorsHandler)
	visitors.POST("", handlers.VisitorHandler)
}
