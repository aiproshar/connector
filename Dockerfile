FROM golang:1.17.6 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go env -w GOPROXY="https://goproxy.io,direct"
RUN go mod download
COPY . .
RUN go build -o /app/bin/app .


FROM ubuntu:latest
RUN apt-get update && apt-get install -y nocache git vim curl ca-certificates && update-ca-certificates
WORKDIR /app
RUN mkdir  -p "/klovercloud/connector"
COPY --from=builder /app/bin /app
#EXPOSE 8080
CMD ["./app"]
