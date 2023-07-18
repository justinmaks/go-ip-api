# go-ip-api

## Getting Started

To run this API locally, you need to have Go installed on your system. If you don't have it installed, you can download it from the official website: https://golang.org/dl/

### Clone repo 

Clone this repository to your local machine using the following command:
```
git clone <repository-url>
cd <repository-name>
```

### Install deps

The API uses external packages for IP geolocation. Install these packages using go get:
```
go get -u github.com/gin-gonic/gin
go get -u github.com/oschwald/geoip2-golang
go get -u github.com/oschwald/maxminddb-golang
```

### IPStack API

To use IP geolocation, you need an API key from IPStack. If you don't have one, sign up at https://ipstack.com/ and obtain your API key.

Once you have the API key, set it as an environment variable in your terminal:
```
export IPSTACK_API_KEY=<your-api-key>
```

### Run the API

Run the API using the following command:
```
go run main.go
```

The API will start running on localhost:8080.

## Endpoints

### GET /

This endpoint returns the visitor's IP address, country, region, and city based on their IP geolocation.

Ex response:
```
{
  "ip": "203.0.113.1",
  "country": "United States",
  "region": "California",
  "city": "San Francisco"
}
```

### GET /visits

This endpoint returns a JSON array containing logs of all the visits to the "/" endpoint, including the IP addresses and corresponding countries.

Example response:
```
[
  {
    "ip": "203.0.113.1",
    "country": "United States"
  },
  {
    "ip": "198.51.100.5",
    "country": "Canada"
  }
]
```

## Data persistence 

As of now, the API uses an in-memory slice to store visit logs. Keep in mind that this data is not persistent, and it will be lost upon restarting the API. If you require data persistence, consider using a database solution like PostgreSQL, MySQL, or MongoDB.

