.PHONY: help dev build run clean test

# デフォルトターゲット
help:
	@echo "利用可能なコマンド:"
	@echo "  make dev    - airを使って開発モードで起動（ファイル変更を監視）"
	@echo "  make build  - アプリケーションをビルド"
	@echo "  make run    - アプリケーションを起動"
	@echo "  make clean  - ビルドファイルを削除"
	@echo "  make test   - テストを実行"
	@echo "  make swagger-  Swaggerドキュメントを生成"

# 開発モード（air使用）
dev:
	@echo "開発モードで起動中..."
	air

# ビルド
build:
	@echo "アプリケーションをビルド中..."
	go build -o ./tmp/main ./main/server.go

# 実行
run:
	@echo "アプリケーションを起動中..."
	go run main/server.go

# クリーンアップ
clean:
	@echo "ビルドファイルを削除中..."
	rm -rf ./tmp
	rm -f build-errors.log

# テスト
test:
	@echo "テストを実行中..."
	go test ./...

# Swaggerドキュメント生成
swagger:
	@echo "Swaggerドキュメントを生成中..."
	swag init -g main/server.go -o docs

# 依存関係の整理
deps:
	@echo "依存関係を整理中..."
	go mod tidy

# Docker起動
docker-up:
	@echo "Dockerでアプリケーションを起動中..."
	docker-compose up -d

# Docker停止
docker-down:
	@echo "Dockerでアプリケーションを停止中..."
	docker-compose down 