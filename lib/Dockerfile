# Start with a lightweight Linux distro (alpine) and Golang
FROM golang:alpine

# Enable Go modules
ENV GO111MODULE=on

# Disable CGO to prevent `go build` asking for gcc
ENV CGO_ENABLED=0

# Install git to allow `go get`
RUN apk add --update --no-cache git

# Create a directory called "app" in our container
RUN mkdir /app

# Copy the current module into the "app" directory of our container
COPY . /app/

# Set the "app" directory of the container as our working directory
WORKDIR /app

# Install required Go packages
RUN go get

# Visit each package and run `go test`
RUN for dir in ./*/; do cd ${dir} && go test && cd ../; done