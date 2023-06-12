FROM golang:1.19 as builder

WORKDIR /app

COPY . .

RUN ls -all
# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/main.go

# Run
CMD ["cb"]