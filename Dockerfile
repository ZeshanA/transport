# Start with a lightweight Linux distro (alpine) and Golang
FROM golang:alpine

# Create a directory called "app" in our container
RUN mkdir /app

# Copy the backend module into the "app" directory of our container
COPY /backend /app/

# Set the "app" directory of the container as our working directory
WORKDIR /app

# Build the Go backend inside our container
RUN go build -o main .

# Run the compiled executable
CMD ["/app/main"]