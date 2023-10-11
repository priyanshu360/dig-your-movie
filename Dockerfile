# Use the official Go image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Install dnsutils package to get 'dig'
RUN apt-get update && apt-get install -y dnsutils

# Copy the Go application source code to the container
COPY . .

# Build the Go application
RUN go build -o my-go-app


# Define the command to run your Go application
CMD ["./my-go-app"]
