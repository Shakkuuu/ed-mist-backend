# ビルドステージ
FROM golang:1.25.1 AS builder

WORKDIR /app

# 依存関係のコピーとインストール
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# アプリケーションのビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

# 実行ステージ
FROM alpine:3.20
RUN apk update && apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/server /app/server
EXPOSE 8080
ENTRYPOINT ["/app/server"]
