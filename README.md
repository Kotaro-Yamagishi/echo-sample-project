<div id="top"></div>

## 使用技術一覧

<!-- シールド一覧 -->
<!-- 該当するプロジェクトの中から任意のものを選ぶ-->
<p style="display: inline">
  <img src="https://img.shields.io/badge/-Go-76E1FE.svg?logo=go&style=for-the-badge">
  <img src="https://img.shields.io/badge/-Docker-1488C6.svg?logo=docker&style=for-the-badge">
</p></p>

## 目次

- [使用技術一覧](#使用技術一覧)
- [目次](#目次)
- [プロジェクト名](#プロジェクト名)
- [プロジェクトについて](#プロジェクトについて)
- [環境](#環境)
- [開発環境構築](#開発環境構築)
  - [コンテナの作成と起動](#コンテナの作成と起動)
  - [動作確認](#動作確認)
  - [環境変数の一覧](#環境変数の一覧)
  - [コマンド一覧](#コマンド一覧)
- [トラブルシューティング](#トラブルシューティング)

<!-- プロジェクト名を記載 -->

## プロジェクト名

Echo-Gorm Project

<!-- プロジェクトについて -->

## プロジェクトについて

Echoとgormを使用したデモアプリケーションです


<p align="right">(<a href="#top">トップへ</a>)</p>

## 環境

<!-- 言語、フレームワーク、ミドルウェア、インフラの一覧とバージョンを記載 -->

| 言語・フレームワーク  | バージョン |
| --------------------- | ---------- |
| Go                    | 1.18       |
| echo                  | 4.7.2      |
| gorm                  | 1.9.16     |
| MySQL                 | 8.0        |

その他のパッケージのバージョンは go.mod を参照してください

<p align="right">(<a href="#top">トップへ</a>)</p>

## 開発環境構築

<!-- コンテナの作成方法、パッケージのインストール方法など、開発環境構築に必要な情報を記載 -->

### コンテナの作成と起動

1. mysqlのdockerコンテナをビルド＆起動
```
cd docker/mysql
docker compose up
```
2. go アプリケーションの起動

### 動作確認

http://localhost:1323 にアクセスできるか確認
アクセスできたら成功

APIは3つ

get  http://localhost:1323/users

get  http://localhost:1323/users/:userId

post http://localhost:1323/users

post時のbody例
```
 {
   "id":"userId1",
   "name":"user1"
 }
 ```

### 環境変数の一覧
なし

### コマンド一覧
なし

## トラブルシューティング
なし

