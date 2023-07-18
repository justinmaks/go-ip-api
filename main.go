package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// VisitLog represents a log entry for each visit to the "/" endpoint
type VisitLog struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
}

var visitLogs []VisitLog
var usRequestCounter int
var nonUSRequestCounter int

func main() {
	router := gin.Default()

	// Route handler for getting the user's IP address and location information
	router.GET("/", func(c *gin.Context) {
		ip := c.ClientIP()

		location, err := getLocationInfo(ip)
		if err != nil {
			fmt.Println("Failed to get location info:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		// Increment the appropriate counter based on the country
		if location.CountryName == "United States" {
			usRequestCounter++
		} else {
			nonUSRequestCounter++
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
	})

	// Endpoint to retrieve all visit logs
	router.GET("/visits", func(c *gin.Context) {
		c.JSON(http.StatusOK, visitLogs)
	})

	// New endpoint to get statistics for US and non-US requests
	router.GET("/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"us_requests":     usRequestCounter,
			"non_us_requests": nonUSRequestCounter,
		})
	})

	// Start the server
	router.Run(":8080")
}

func getLocationInfo(ip string) (*LocationInfo, error) {
	apiKey := os.Getenv("IPSTACK_API_KEY")
	url := fmt.Sprintf("http://api.ipstack.com/%s?access_key=%s", ip, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	var location LocationInfo
	err = json.NewDecoder(resp.Body).Decode(&location)
	if err != nil {
		return nil, err
	}

	return &location, nil
}

// LocationInfo represents the location information retrieved from the ipstack API
type LocationInfo struct {
	CountryName string `json:"country_name"`
	RegionName  string `json:"region_name"`
	City        string `json:"city"`
}
