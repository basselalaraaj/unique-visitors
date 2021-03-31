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

type visitors = map[string]struct{}
type page = map[string]visitors

var hits = struct {
	sync.RWMutex
	page page
}{page: make(page)}

func addVisitor(url string, visitor string) {
	hits.Lock()
	_, found := hits.page[url]
	if !found {
		hits.page[url] = make(map[string]struct{})
	}

	_, found = hits.page[url][visitor]
	if found {
		hits.Unlock()
		return
	}

	hits.page[url][visitor] = struct{}{}
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
	url := c.Query("url")
	if url == "" {
		c.String(http.StatusBadRequest, "missing required url query.")
		return
	}
	c.String(http.StatusOK, "%d", getVisitors(url))
}
