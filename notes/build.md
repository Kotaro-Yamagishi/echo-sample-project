## Goのビルドプロセスの概要
Goソースコード
   ↓
go build
   ↓
① コンパイル（.go → 中間表現 SSA → 機械語）
② リンカでバイナリ生成（外部パッケージなどを含めて一体化）
③ 実行ファイル（静的リンクされたバイナリ）が生成される

## 学習テーマ
### Goビルドの仕組み
go build, go install, go mod, go run
実際のワークフロー

### コンパイラ理論
go tool comile
SSA（Static Single Assignment）や中間表現（IR）、型チェックの仕組み

### リンカの仕組み
go build -ldflags
Goの性的リンクと動的リンクの違い、リンカフラグ（-ldflangs）

### Goランタイム
GoDEBUG変数など
ガベージコレクション、スケジューラ（M:Nスケジューリング）、ヒープ管理

### アセンブリとの関係
Goコード→ASM
go tool compile -S, objdump を使って機械語を見る方法

### ツールチェインの理解
go tool compile, go tool link, go tool objdump の使い方


## コンパイラ理論
### Goコンパイラのステージ
Goコード
  ↓ (パーサ)
抽象構文木 (AST)
  ↓ (型チェック)
型付きAST
  ↓ (SSA生成)
中間表現 (IR)
  ↓ (最適化/コード生成)
マシンコード

### SSAとは
各変数が一度だけ代入される形式に変換した中間表現のこと