# Use an official Go runtime as a parent image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy the local code into the container at the working directory
COPY . .

# Build the Go application
RUN go build -o app ./app

COPY configs /configs

EXPOSE 8080

# Specify the command to run on container start
CMD ["./app"]
