# Set base image
FROM golang:1.20.6-alpine

# Build application
WORKDIR /app/
COPY ./source /app
RUN go mod tidy

# Run tests
CMD ["go", "test",  "-v", "./..."]
