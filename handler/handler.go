package handler

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupHandler(app *gin.Engine) {
	app.GET("/metrics", gin.WrapH(promhttp.Handler()))
	
	// Test endpoints
	app.GET("/test/fast", FastResponse)
	app.GET("/test/slow", SlowResponse)
	app.GET("/test/random", RandomResponse)
	app.GET("/test/error", ErrorResponse)
	app.GET("/test/ok", OkResponse)
}

// FastResponse returns quickly for testing metrics
func FastResponse(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Fast response",
		"time":    time.Now().String(),
	})
}

// SlowResponse simulates a slow API for testing metrics
func SlowResponse(c *gin.Context) {
	// Sleep between 1-3 seconds
	sleepTime := time.Duration(1000+rand.Intn(4000)) * time.Millisecond
	time.Sleep(sleepTime)
	
	c.JSON(http.StatusOK, gin.H{
		"message":    "Slow response",
		"sleep_time": sleepTime.String(),
		"time":       time.Now().String(),
	})
}

// RandomResponse randomly succeeds or fails
func RandomResponse(c *gin.Context) {
	// Randomly generate status codes
	statuses := []int{
		http.StatusOK,
		http.StatusBadRequest,
		http.StatusInternalServerError,
		http.StatusOK,
		http.StatusOK,
	}
	
	status := statuses[rand.Intn(len(statuses))]
	
	c.JSON(status, gin.H{
		"message": "Random response",
		"status":  status,
		"time":    time.Now().String(),
	})
}

// ErrorResponse always returns an error
func ErrorResponse(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": "This endpoint always returns an error",
		"time":    time.Now().String(),
	})
}

// OkResponse always returns a successful response
func OkResponse(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"status":  "success",
		"time":    time.Now().String(),
	})
}
