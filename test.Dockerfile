# Set base image
FROM golang:1.20.6-alpine

# Build application
WORKDIR /app/
COPY ./source /app
RUN go mod tidy
RUN go install gotest.tools/gotestsum@latest

# Run tests
ENTRYPOINT gotestsum --format standard-verbose

