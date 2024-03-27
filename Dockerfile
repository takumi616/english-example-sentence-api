#Creates binary which is included in deploy container
FROM golang:1.22.0-alpine3.18 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -trimpath -ldflags "-w -s" -o main

#-------------------------------------------------------------------------

#Deploy 
FROM alpine:latest as deploy
RUN apk update
COPY --from=builder /app/main .
CMD ["./main"]

#-------------------------------------------------------------------------

#Local development environment(hot reload)
FROM golang:1.22 as dev
WORKDIR /app
RUN go install github.com/cosmtrek/air@latest
CMD ["air"]
