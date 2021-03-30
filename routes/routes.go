package routes

import (
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func GetRoutes() *gin.Engine {
	AddVisitorsRoutes(router)
	return router
}
