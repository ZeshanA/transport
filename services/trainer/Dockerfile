# Start with a lightweight Linux distro (alpine) and Golang
FROM python:3-alpine

# SERVICE_NAME contains the name of the service currently being built
ARG SERVICE_NAME

# Create an "app" directory in the container to isolate our files
RUN mkdir /app
WORKDIR /app

# Create a directory for our service, under the "services" directory in our container
RUN mkdir -p services/${SERVICE_NAME}

# Copy the current service folder into the "services" directory of our container
COPY services/${SERVICE_NAME} services/${SERVICE_NAME}

# Set the service's directory inside the container as our working directory
WORKDIR services/${SERVICE_NAME}

# Install required Go packages
RUN pip install -r requirements.txt

# Visit each package and run any tests
RUN python -m unittest discover -s .

# Run the main file
CMD ["python3", "main/main.py"]