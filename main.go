package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// VisitLog represents a log entry for each visit to the "/" endpoint
type VisitLog struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
}

type LocationInfo struct {
	CountryName string `json:"country_name"`
	RegionName  string `json:"region_name"`
	City        string `json:"city"`
}

var visitLogs []VisitLog
var usRequestCounter int
var nonUSRequestCounter int
var apiKey string
var counterMutex = &sync.Mutex{} // Mutex to avoid race condition when incrementing counters

func main() {
	// Load API Key at start
	apiKey = os.Getenv("IPSTACK_API_KEY")
	if apiKey == "" {
		log.Fatal("IPSTACK_API_KEY not set")
	}

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST"}
	config.AllowHeaders = []string{"Origin"}
	router.Use(cors.New(config))

	// Route handler for getting the user's IP address and location information
	router.GET("/", handleRoot)

	// Endpoint to retrieve all visit logs
	router.GET("/visits", handleVisits)

	// New endpoint to get statistics for US and non-US requests
	router.GET("/stats", handleStats)

	// Start the server
	router.Run(":8080")
}

func handleRoot(c *gin.Context) {
	ip := c.ClientIP()

	locationCh := make(chan *LocationInfo, 1)
	errCh := make(chan error, 1)
	go getLocationInfo(context.Background(), ip, locationCh, errCh)

	location := <-locationCh
	if err := <-errCh; err != nil {
		log.Println("Failed to get location info:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	// Increment the appropriate counter based on the country
	if location.CountryName == "United States" {
		counterMutex.Lock() // Avoid race condition when incrementing
		usRequestCounter++
		counterMutex.Unlock()
	} else {
		counterMutex.Lock() // Avoid race condition when incrementing
		nonUSRequestCounter++
		counterMutex.Unlock()
	}

	// Create a new log entry and add it to the visitLogs slice
	visitLogs = append(visitLogs, VisitLog{
		IP:      ip,
		Country: location.CountryName,
	})

	c.JSON(http.StatusOK, gin.H{
		"ip":      ip,
		"country": location.CountryName,
		"region":  location.RegionName,
		"city":    location.City,
	})
}

func handleVisits(c *gin.Context) {
	c.JSON(http.StatusOK, visitLogs)
}

func handleStats(c *gin.Context) {
	counterMutex.Lock() // Avoid race condition when reading counters
	stats := gin.H{
		"us_requests":     usRequestCounter,
		"non_us_requests": nonUSRequestCounter,
	}
	counterMutex.Unlock()
	c.JSON(http.StatusOK, stats)
}

func getLocationInfo(ctx context.Context, ip string, locationCh chan<- *LocationInfo, errCh chan<- error) {
	url := fmt.Sprintf("http://api.ipstack.com/%s?access_key=%s", ip, apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		errCh <- err
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errCh <- err
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errCh <- fmt.Errorf("unexpected response status: %d", resp.StatusCode)
		return
	}

	var location LocationInfo
	err = json.NewDecoder(resp.Body).Decode(&location)
	if err != nil {
		errCh <- err
		return
	}

	locationCh <- &location
	errCh <- nil
}
