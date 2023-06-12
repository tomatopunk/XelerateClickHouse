FROM golang:1.19 as builder

WORKDIR /app

COPY . .

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o cb cmd/main.go

ENV PATH="/app:${PATH}"

# Run
CMD ["cb","--help"]