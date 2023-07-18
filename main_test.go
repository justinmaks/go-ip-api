package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetLocationInfo(t *testing.T) {
	// Provide a mock IP address for testing
	ip := "8.8.8.8"

	// Provide a mock LocationInfo response
	mockLocation := LocationInfo{
		CountryName: "United States",
		RegionName:  "California",
		City:        "Mountain View",
	}

	// Mock the HTTP response from the IPStack API
	mockResponse := struct {
		LocationInfo
	}{
		LocationInfo: mockLocation,
	}
	responseJSON, _ := json.Marshal(mockResponse)
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	}))
	defer mockServer.Close()

	// Set the mock IPStack API URL to the test server URL
	os.Setenv("IPSTACK_API_KEY", "test")
	os.Setenv("IPSTACK_API_URL", mockServer.URL)

	// Test the getLocationInfo function
	location, err := getLocationInfo(ip)
	if err != nil {
		t.Errorf("Error fetching location info: %v", err)
	}

	// Compare the location obtained with the expected mock location
	if location.CountryName != mockLocation.CountryName ||
		location.RegionName != mockLocation.RegionName ||
		location.City != mockLocation.City {
		t.Errorf("getLocationInfo returned incorrect location info: got %+v, want %+v", location, mockLocation)
	}
}

func TestAPIEndpoints(t *testing.T) {
	// Create a new Gin router with the same routes as the main API
	router := gin.New()
	router.GET("/", handleGetLocation)
	router.GET("/visits", handleGetVisits)

	// Test the "/" endpoint
	w := performRequest(router, "GET", "/")
	assertResponseCode(t, w.Code, http.StatusOK)

	// Test the "/visits" endpoint
	w = performRequest(router, "GET", "/visits")
	assertResponseCode(t, w.Code, http.StatusOK)
}

// Helper functions

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func assertResponseCode(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Expected response code %d, got %d", want, got)
	}
}
