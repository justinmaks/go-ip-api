package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

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

		c.JSON(http.StatusOK, gin.H{
			"ip":      ip,
			"country": location.CountryName,
			"region":  location.RegionName,
			"city":    location.City,
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
