# 第一階段：build
FROM golang:1.24.1 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api_server .

# 第二階段：運行
FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/api_server .

# 默認 port 可改
EXPOSE 8080

# 執行你的 Go 程序
CMD ["./api_server"]