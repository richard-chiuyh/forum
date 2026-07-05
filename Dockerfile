# FROM golang:1.23.0-alpine AS builder  #   使用更加新的版本
FROM golang:1.25.0-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct
ENV GOPRIVATE=git.168tzi.com

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go mod tidy
RUN go build -ldflags="-s -w" -o /app/forum forum.go


FROM alpine

# RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata    #   移到上面
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/forum /app/forum
COPY etc /app/etc

CMD ["./forum", "-f", "etc/forum.yaml"]
