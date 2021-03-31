package routes

import (
	"github.com/gin-gonic/gin"
)

func GetRoutes() *gin.Engine {
	router := gin.Default()
	addVisitorsRoutes(router)
	return router
}
