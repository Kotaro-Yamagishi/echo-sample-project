## SQLBoilerとは
既存のデータベーススキーマをもとに、Goコードを自動生成してくれるORMライブラリ
自動生成されるコードは型安全、IDE補完も効きやすい、実行前のコンパイルエラーで問題の検出が可能

cd /echo-sample-project && sqlboiler --config main/infra/things/sqlboiler/sqlboiler.toml mysql