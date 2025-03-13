package main

import (
	"abuabdillatief/sample-golang-prometheus-grafana/handler"
	"abuabdillatief/sample-golang-prometheus-grafana/middleware"

	"github.com/gin-gonic/gin"
)

var (
	internalServerError = "Internal Server Error"
	notFound            = "Not Found"
)

func main() {
	r := gin.Default()
    
	middleware.PrometheusInit()        
	r.Use(middleware.TrackMetrics)
    handler.SetupHandler(r)

	r.Run("127.0.0.1:8080")
}
