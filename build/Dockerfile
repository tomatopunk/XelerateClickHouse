FROM golang:1.19 as builder

WORKDIR /app

COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o cb cmd/main.go

# Run
CMD ["cb"]