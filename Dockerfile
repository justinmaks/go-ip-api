FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go project files to the container's working directory
COPY . .

# Install necessary dependencies
RUN go get -u github.com/gin-gonic/gin
RUN go get -u github.com/oschwald/geoip2-golang
RUN go get -u github.com/oschwald/maxminddb-golang

# Build the Go application
RUN go build -o main .

# Expose the port the application listens on
EXPOSE 8080

# Command to run the application when the container starts
CMD ["./main"]
