package handlers

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type event struct {
	URL     string `json:"URL" binding:"required"`
	Visitor string `json:"Visitor" binding:"required"`
}

type visitors = map[string]int
type page = map[string]visitors

var hits = struct {
	sync.RWMutex
	page page
}{page: make(page)}

func addPageIfNotFound(url string) {
	hits.RLock()
	_, found := hits.page[url]
	hits.RUnlock()
	if found {
		return
	}

	hits.Lock()
	hits.page[url] = make(map[string]int)
	hits.Unlock()
}

func addVisitor(url string, visitor string) {
	addPageIfNotFound(url)

	hits.RLock()
	_, found := hits.page[url][visitor]
	hits.RUnlock()
	if found {
		return
	}

	hits.Lock()
	hits.page[url][visitor] = 1
	hits.Unlock()
}

func getVisitors(url string) int {
	hits.RLock()
	totalUniqueHits := len(hits.page[url])
	hits.RUnlock()
	return totalUniqueHits
}

func VisitorHandler(c *gin.Context) {
	var newEvent event
	if err := c.ShouldBindJSON(&newEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	addVisitor(newEvent.URL, newEvent.Visitor)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func VisitorsHandler(c *gin.Context) {
	url := c.DefaultQuery("url", "none")
	if url == "none" {
		c.String(http.StatusBadRequest, "missing required url query.")
		return
	}
	c.String(http.StatusOK, "%d", getVisitors(url))
}
