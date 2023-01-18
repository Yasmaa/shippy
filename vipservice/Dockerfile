# Base Image
FROM golang:1.16-alpine

# Make app directory
RUN mkdir /app

# Copy all content to the app directory
ADD . /app

# Make app directory the working directory
WORKDIR /app

# Copy
COPY go.mod ./
COPY go.sum ./

# Download any required modules
RUN go mod download

# Build the program to create an executable binary
RUN go build -o webserver .

EXPOSE 9092
# Set the startup command
CMD ["/app/webserver"]