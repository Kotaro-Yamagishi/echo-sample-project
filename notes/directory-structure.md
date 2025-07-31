# ディレクトリ構成
## app
frameworkに依存したコード
echoやwire等

### app/controller
HTTPリクエストを受け取り、ビジネスロジック（Service層）を呼び出し、レスポンスを返す役割

### app/DI
DI関連。
今回はwireを使った依存関係依存関係解決コードの解決コードの自動生成ツール配置

### app/router
APIのエンドポイント作成

## domain
### controller
app/controllerのinterface

### datasource
infra/datasourceのinterface

### entity
app,usecaseで利用するstructの管理

### model
sqlboilerが自動生成したmodelを管理

### repository
infra/repositoryのinterfaceを管理

### service
usecaseのinterface

## infra
### datasource
複雑なクエリを生成し、呼び出し
それ以外の処理は行わない

### repository
呼び出したクエリに対して整形やキャッシュ管理等

### things
datasourceやrepositoryで実現できないことを管理
ex）リトライ処理等

## usecase
ビジネスロジック処理


# 設計思想的アドバイス
- 責務の分離
  - パッケージごとに「1つの役割」を持たせ、責務を混在させない。
- 依存の方向
  - controller → service → repository の一方向。逆方向依存はNG（interface + DIで実現）。
- interface重視
  - 特にservice/repositoryはinterfaceを用意し、実装と分離することでテストしやすくなる
- Echo, GORMなどの外部依存を隔離
  - 外部ライブラリに依存するコードは、必ずインフラ層（repositoryやhandler）に留める。
- 単体テストしやすい構成に	
  - モック差し替え可能に、interfaceと初期化コードを工夫する。

# Interface
## go言語におけるInterfaceの役割
- 「このパッケージ（層）がこれだけの機能を必要とする」という要求の表明.「この型がこの機能を提供する」設計ではない
- 利用者が意図を理解するため、interfaceは利用者側パッケージに置くのが原則

## 構造体ファイルとinterfaceファイルは分けるべき？
密結合を避けるため、「interface」と「struct」は分けて管理


## レイヤードアーキテクチャ

## クリーンアーキテクチャ
