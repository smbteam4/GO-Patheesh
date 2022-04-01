# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:latest

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
# RUN apk update && apk upgrade && \
#     apk add --no-cache bash git openssh

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .


# Build the app
RUN go build -o robot-apocalypse /app/main.go

# Document that the service listens on port 8080.
EXPOSE 8080

# Run the outyet command by default when the container starts.
ENTRYPOINT /app/robot-apocalypse
