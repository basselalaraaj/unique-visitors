package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAddVisitor(t *testing.T) {
	testCases := []struct {
		url     string
		visitor string
		want    int
	}{
		{"http://www.example.com", "123", 1},
		{"http://www.example.com", "123", 1},
		{"http://www.example.com", "123", 1},
		{"http://www.example.com", "456", 2},
		{"http://www.example.com", "456", 2},
		{"http://www.example.com", "789", 3},
		{"http://www.example.com/home", "456", 1},
		{"http://www.example.com/home", "456", 1},
		{"http://www.example.com/home", "789", 2},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("URL %s and visitor %s", tc.url, tc.visitor), func(t *testing.T) {
			addVisitor(tc.url, tc.visitor)
			if got := getVisitors((tc.url)); got != tc.want {
				t.Errorf("URL %s and visitor %s, got %d, want %d", tc.url, tc.visitor, got, tc.want)
			}
		})
	}
}

func TestVisitorsHandler(t *testing.T) {
	t.Run("Add visitor wrong body", func(t *testing.T) {
		router := gin.Default()
		gin.SetMode(gin.TestMode)
		router.POST("/visitors", VisitorHandler)

		body := strings.NewReader(`{
			"visitor": "123"
		}`)

		w := httptest.NewRecorder()
		ctx := context.Background()
		req, _ := http.NewRequestWithContext(ctx, "POST", "/visitors", body)
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, `{"error":"Key: 'event.URL' Error:Field validation for 'URL' failed on the 'required' tag"}`, w.Body.String())
	})

	t.Run("Add visitor", func(t *testing.T) {
		router := gin.Default()
		gin.SetMode(gin.TestMode)
		router.POST("/visitors", VisitorHandler)

		body := strings.NewReader(`{
			"url": "http://www.website.com",
			"visitor": "123"
		}`)

		w := httptest.NewRecorder()
		ctx := context.Background()
		req, _ := http.NewRequestWithContext(ctx, "POST", "/visitors", body)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, `{"status":"ok"}`, w.Body.String())
	})

	t.Run("Visitors is 1", func(t *testing.T) {
		router := gin.Default()
		gin.SetMode(gin.TestMode)
		router.GET("/visitors", VisitorsHandler)

		w := httptest.NewRecorder()
		ctx := context.Background()
		req, _ := http.NewRequestWithContext(ctx, "GET", "/visitors?url=http://www.website.com", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "1", w.Body.String())
	})

	t.Run("Visitors missing url", func(t *testing.T) {
		router := gin.Default()
		gin.SetMode(gin.TestMode)
		router.GET("/visitors", VisitorsHandler)

		w := httptest.NewRecorder()
		ctx := context.Background()
		req, _ := http.NewRequestWithContext(ctx, "GET", "/visitors", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, "missing required url query.", w.Body.String())
	})
}
