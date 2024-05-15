FROM golang:1.22
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code and configuration.
COPY *.go ./
COPY *.ini ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /oauth2-proxy-service

EXPOSE 8080

# Run
CMD ["/oauth2-proxy-service"]