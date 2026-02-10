.PHONY: build run test docker-up docker-down

# アプリケーションのビルド
build:
	go build -o api cmd/api/main.go

# アプリケーションの実行 (ソースから)
run:
	go run cmd/api/main.go

# 実行バイナリの起動
start: build
	./api

# テストの実行
test:
	go test ./tests/...

# 開発環境の起動 (Docker)
docker-up:
	docker compose up -d

# 開発環境の停止
docker-down:
	docker compose down
