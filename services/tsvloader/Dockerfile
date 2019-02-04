# Start with a lightweight Linux distro (alpine) and Golang
FROM golang:alpine

# SERVICE_NAME contains the name of the service currently being built
ARG SERVICE_NAME

# Enable Go modules
ENV GO111MODULE=on

# Disable CGO to prevent `go build` asking for gcc
ENV CGO_ENABLED=0

# Install git to allow `go get`
RUN apk add --update --no-cache git

# Create an "app" directory in the container to isolate our files
RUN mkdir /app
WORKDIR /app

# Create a "lib" directory to hold the libraries module
RUN mkdir lib

# Copy our library module over
COPY lib lib

# Create a directory for our service, under the "services" directory in our container
RUN mkdir -p services/${SERVICE_NAME}

# Copy the current service folder into the "services" directory of our container
COPY services/${SERVICE_NAME} services/${SERVICE_NAME}

# Set the service's directory inside the container as our working directory
WORKDIR services/${SERVICE_NAME}

# Install required Go packages
RUN go get

# Build the Go executable inside our container
RUN go build -o executable ./...

# Open port 80 for requests
EXPOSE 80

# Run the compiled executable
CMD ["./executable"]