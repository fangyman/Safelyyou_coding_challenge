package main

import (
	"log"
	"path/filepath"

	"safelyyou_coding_challenge_go/api"
	"safelyyou_coding_challenge_go/database"
	"safelyyou_coding_challenge_go/services"

	"github.com/gin-gonic/gin"
)

const (
	databaseURL = "./fleet_monitoring.db"
	csvPath     = "devices.csv"
	serverHost  = "127.0.0.1"
	serverPort  = "6733"
)

func main() {
	log.Println("Starting Fleet Monitoring API...")

	log.Println("Initializing database...")

	dbPath := filepath.FromSlash(databaseURL)

	dbOps, err := database.InitDatabase(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Println("Loading devices from CSV...")
	csvService := services.NewCSVService(dbOps)
	if err := csvService.LoadDevicesFromCSV(csvPath); err != nil {
		log.Fatalf("Failed to load devices from CSV: %v", err)
	}
	log.Println("Startup complete.")

	r := gin.Default()

	r.SetTrustedProxies([]string{"127.0.0.1"})

	metricsService := services.NewMetricsService(dbOps)
	h := api.NewAPIHandler(dbOps, metricsService)

	r.GET("/", h.Root)

	v1 := r.Group("/api/v1/devices")
	{
		v1.POST("/:device_id/heartbeat", h.ReceiveHeartbeat)
		v1.POST("/:device_id/stats", h.ReceiveStats)
		v1.GET("/:device_id/stats", h.GetDeviceStats)
	}

	// Run the server on port 6733
	serverAddr := serverHost + ":" + serverPort
	log.Printf("Server listening on http://%s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
