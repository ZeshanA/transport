# Start with a lightweight Linux distro (alpine) and Golang
FROM golang:alpine

# SERVICE_NAME contains the name of the service currently being built
ARG SERVICE_NAME

# Set the default value of the DEBUG build arg to false
ARG DEBUG="false"

# Set the DEBUG environment value so that we can pass it to
# the entrypoint script
ENV DEBUG=${DEBUG}

RUN echo "DEBUG IS ${DEBUG}"

# Add compile flags for no optimisation if in Debug mode
RUN if [ "${DEBUG}"="true" ] ; then export COMPILE_FLAGS="all=-N -l" ; else export COMPILE_FLAGS="" ; fi

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

# If debugging, install Delve debugger
RUN if [ "${DEBUG}" = "true" ] ; then go get github.com/derekparker/delve/cmd/dlv && apk add --no-cache libc6-compat ; fi

# Visit each package and run `go test`
RUN for dir in ./*/; do cd ${dir} && go test && cd ../; done