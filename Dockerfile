# Start from golang base image
FROM golang:1.18-alpine as builder

# Enable go modules
ENV GO111MODULE=on

# Set current working directory
WORKDIR /app

# Note here: To avoid downloading dependencies every time we
# build image. Here, we are caching all the dependencies by
# first copying go.mod and go.sum files and downloading them,
# to be used every time we build the image if the dependencies
# are not changed.
COPY go.mod .
# Download all dependencies.
RUN go mod download

COPY . .

# Note here: CGO_ENABLED is disabled for cross system compilation
# It is also a common best practise.

# Build the application.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/main.go


# Finally our multi-stage to build a small image
# Start a new stage from scratch
FROM scratch

#EXPOSE 9090

# Copy the Pre-built binary file
COPY --from=builder /app/main .

# Run executable
CMD ["/main"]