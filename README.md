# Echo Sample Project

<div id="top"></div>

## 使用技術一覧

<!-- シールド一覧 -->
<!-- 該当するプロジェクトの中から任意のものを選ぶ-->
<p style="display: inline">
  <img src="https://img.shields.io/badge/-Go-76E1FE.svg?logo=go&style=for-the-badge">
  <img src="https://img.shields.io/badge/-Docker-1488C6.svg?logo=docker&style=for-the-badge">
  <img src="https://img.shields.io/badge/-MySQL-4479A1.svg?logo=mysql&style=for-the-badge">
  <img src="https://img.shields.io/badge/-Echo-000000.svg?logo=echo&style=for-the-badge">
</p>

## 目次

- [使用技術一覧](#使用技術一覧)
- [目次](#目次)
- [プロジェクト名](#プロジェクト名)
- [プロジェクトについて](#プロジェクトについて)
- [環境](#環境)
- [開発環境構築](#開発環境構築)
  - [Dockerを使用した起動](#dockerを使用した起動)
  - [ローカル開発環境](#ローカル開発環境)
  - [動作確認](#動作確認)
  - [環境変数の一覧](#環境変数の一覧)
  - [コマンド一覧](#コマンド一覧)
- [トラブルシューティング](#トラブルシューティング)

<!-- プロジェクト名を記載 -->

## プロジェクト名

Echo Sample Project

<!-- プロジェクトについて -->

## プロジェクトについて

EchoとSQLBoilerを使用したサンプルアプリケーションです。Sakilaデータベースを使用して、都市と国の情報を管理するAPIを提供します。

<p align="right">(<a href="#top">トップへ</a>)</p>

## 環境

<!-- 言語、フレームワーク、ミドルウェア、インフラの一覧とバージョンを記載 -->

| 言語・フレームワーク  | バージョン |
| --------------------- | ---------- |
| Go                    | 1.24.2     |
| Echo                  | 4.9.0      |
| SQLBoiler             | 4.19.5     |
| MySQL                 | 8.0        |

その他のパッケージのバージョンは go.mod を参照してください

<p align="right">(<a href="#top">トップへ</a>)</p>

## 開発環境構築

<!-- コンテナの作成方法、パッケージのインストール方法など、開発環境構築に必要な情報を記載 -->

### Dockerを使用した起動

#### 本番環境用（推奨）
シンプルで軽量なDockerfileを使用します。

1. 環境変数ファイルの作成
```bash
# .envファイルを作成
cat > .env << EOF
DB_HOST=db
DB_PORT=3306
DB_NAME=sakila
DB_USER=user
DB_PASSWORD=passw0rd
TZ=Asia/Tokyo
DEBUG=true
LOG_LEVEL=info
EOF
```

2. アプリケーションとMySQLデータベースを同時に起動
```bash
docker-compose up -d
```

#### 開発環境用
開発ツール（linter、swagger、wire等）を含むDockerfileを使用します。

```bash
# 開発用Dockerfileでビルド
docker build -f Dockerfile.dev -t echo-sample-dev .

# 開発用コンテナを起動
docker run -d --name echo-sample-dev \
  --env-file .env \
  -p 1324:1324 \
  echo-sample-dev
```

2. ログの確認
```bash
docker-compose logs -f app
```

3. アプリケーションの停止
```bash
docker-compose down
```

### ローカル開発環境

1. 環境変数ファイルの作成（オプション）
```bash
# .envファイルを作成して環境変数を設定
cat > .env << EOF
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=sakila
DB_USER=user
DB_PASSWORD=passw0rd
TZ=Asia/Tokyo
DEBUG=true
LOG_LEVEL=info
EOF
```

2. MySQLのdockerコンテナをビルド＆起動
```bash
cd docker/mysql
docker compose up -d
```

3. Goアプリケーションの起動
```bash
cd main
go run server.go
```

### 動作確認

アプリケーションが起動したら、以下のエンドポイントで動作確認ができます：

- ヘルスチェック: http://localhost:1324/health
- 都市一覧: http://localhost:1324/cities
- 国一覧: http://localhost:1324/countries
- Swagger UI: http://localhost:1324/swagger/

### 環境変数の一覧

| 変数名 | デフォルト値 | 説明 |
|--------|-------------|------|
| DB_HOST | 127.0.0.1 | データベースホスト |
| DB_PORT | 3306 | データベースポート |
| DB_NAME | sakila | データベース名 |
| DB_USER | root | データベースユーザー |
| DB_PASSWORD | password | データベースパスワード |

### コマンド一覧

```bash
# アプリケーションのビルド
go build -o echo-sample main/server.go

# 依存関係の整理
go mod tidy

# SQLBoilerのコード生成
sqlboiler --config main/infra/things/sqlboiler/sqlboiler.toml mysql

# Wireのコード生成
wire main/app/DI/

# Swaggerドキュメントの生成
swag init -g main/server.go -o main/docs
```

<p align="right">(<a href="#top">トップへ</a>)</p>

## トラブルシューティング

### Goのバージョンエラーが発生する場合

```bash
# Goのキャッシュをクリア
go clean -cache -modcache

# 依存関係を再整理
go mod tidy
```

### データベース接続エラーが発生する場合

1. MySQLコンテナが起動しているか確認
```bash
docker-compose ps
```

2. データベースのログを確認
```bash
docker-compose logs db
```

3. データベースに直接接続して確認
```bash
docker exec -it sakila_mysql mysql -u root -ppassword sakila
```

<p align="right">(<a href="#top">トップへ</a>)</p>