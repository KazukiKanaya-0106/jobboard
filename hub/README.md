# Hub - Job Board API Server

Gin + pgxpool + golang-migrate を使用したAPIサーバー

## 🚀 開発環境

```bash
# コンテナ起動
docker-compose up -d

# ログ確認
docker-compose logs -f hub

# コンテナに入る
docker-compose exec hub bash
```

## 📊 マイグレーション

```bash
# マイグレーション実行（最新まで）
docker-compose exec hub go run cmd/migrate/main.go -cmd=up

# マイグレーション実行（1つ戻す）
docker-compose exec hub go run cmd/migrate/main.go -cmd=down

# 現在のバージョン確認
docker-compose exec hub go run cmd/migrate/main.go -cmd=version

# 強制的にバージョン設定（エラー時）
docker-compose exec hub go run cmd/migrate/main.go -cmd=force -version=1
```

## 🔧 新しいマイグレーションの作成

```bash
# 例：jobs テーブルを作成
# migrations/000002_create_jobs_table.up.sql
# migrations/000002_create_jobs_table.down.sql
```

命名規則：`{version}_{description}.{up|down}.sql`

## 📡 API エンドポイント

- `GET /health` - ヘルスチェック（DB接続確認含む）
- `GET /` - API情報

## 🗄️ データベース接続

**pgxpool**を使用：
- 高パフォーマンス接続プール
- PostgreSQL特化の機能
- コンテキスト対応

環境変数で設定（`.env`参照）：
- `DB_HOST` - データベースホスト
- `DB_PORT` - データベースポート
- `DB_USER` - データベースユーザー
- `DB_PASSWORD` - データベースパスワード
- `DB_NAME` - データベース名

## 📦 使用パッケージ

- `github.com/gin-gonic/gin` - Webフレームワーク
- `github.com/jackc/pgx/v5/pgxpool` - PostgreSQL接続プール
- `github.com/golang-migrate/migrate/v4` - マイグレーションツール
