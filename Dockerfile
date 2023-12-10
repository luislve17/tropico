# Set base image
FROM golang:1.20.6-alpine

# Build application
WORKDIR /app/
COPY ./source /app
RUN go mod tidy && go build -o /tropico

# Expose entrypoint
EXPOSE 8000

# Run application
CMD ["/tropico"]

