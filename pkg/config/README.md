# Config Package

`github.com/kelseyhightower/envconfig`を使用して環境変数からアプリケーション設定をロードするパッケージです。

## 機能

- 環境変数からの設定ロード
- ファイルからの設定ロード（.env 形式）
- 設定値のバリデーション
- デフォルト値のサポート
- 構造体ベースの設定管理
- 環境判定ヘルパー関数

## 使用方法

### 基本的な使用方法

```go
package main

import (
    "log"
    "fmt"
    "github.com/pannpers/go-backend-scaffold/pkg/config"
)

func main() {
    // 環境変数から設定をロード
    cfg, err := config.Load("APP")
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // 設定をバリデーション
    if err := cfg.Validate(); err != nil {
        log.Fatalf("Invalid configuration: %v", err)
    }

    // 設定値を使用
    fmt.Printf("Server Port: %d\n", cfg.Server.Port)
    fmt.Printf("Database Name: %s\n", cfg.Database.Name)
}
```

### 環境変数の設定

以下の環境変数を設定できます：

#### 基本設定

- `APP_ENVIRONMENT`: 環境（development, staging, production）
- `APP_DEBUG`: デバッグモード（true/false）

#### サーバー設定

- `APP_SERVER_PORT`: サーバーポート（デフォルト: 8080）
- `APP_SERVER_HOST`: サーバーホスト（デフォルト: localhost）
- `APP_SERVER_READ_TIMEOUT`: 読み取りタイムアウト（秒）
- `APP_SERVER_WRITE_TIMEOUT`: 書き込みタイムアウト（秒）
- `APP_SERVER_IDLE_TIMEOUT`: アイドルタイムアウト（秒）

#### データベース設定

- `APP_DATABASE_HOST`: データベースホスト（デフォルト: localhost）
- `APP_DATABASE_PORT`: データベースポート（デフォルト: 5432）
- `APP_DATABASE_NAME`: データベース名（必須）
- `APP_DATABASE_USER`: データベースユーザー（必須）
- `APP_DATABASE_PASSWORD`: データベースパスワード（必須）
- `APP_DATABASE_SSL_MODE`: SSL モード（デフォルト: disable）
- `APP_DATABASE_MAX_OPEN_CONNS`: 最大接続数（デフォルト: 25）
- `APP_DATABASE_MAX_IDLE_CONNS`: 最大アイドル接続数（デフォルト: 5）
- `APP_DATABASE_CONN_MAX_LIFETIME`: 接続最大生存時間（秒、デフォルト: 300）

#### ログ設定

- `APP_LOGGING_LEVEL`: ログレベル（debug, info, warn, error）
- `APP_LOGGING_FORMAT`: ログ形式（json, text）
- `APP_LOGGING_STRUCTURED`: 構造化ログ（true/false）
- `APP_LOGGING_INCLUDE_CALLER`: 呼び出し元情報を含む（true/false）

### ファイルからの設定ロード

```go
// .envファイルから設定をロード
cfg, err := config.LoadFromFile("APP", ".env")
if err != nil {
    log.Fatalf("Failed to load configuration: %v", err)
}
```

.env ファイルの例：

```
APP_ENVIRONMENT=production
APP_DEBUG=true
APP_SERVER_PORT=9090
APP_DATABASE_NAME=myapp
APP_DATABASE_USER=dbuser
APP_DATABASE_PASSWORD=dbpass
APP_LOGGING_LEVEL=debug
```

### 環境判定ヘルパー

```go
if cfg.IsDevelopment() {
    // 開発環境用の処理
}

if cfg.IsProduction() {
    // 本番環境用の処理
}

if cfg.IsStaging() {
    // ステージング環境用の処理
}
```

### データベース接続文字列の取得

```go
dsn := cfg.Database.GetDSN()
// 例: "host=localhost port=5432 user=dbuser password=dbpass dbname=myapp sslmode=disable"
```

## 設定構造体

### Config

メインの設定構造体です。

```go
type Config struct {
    Server     ServerConfig
    Database   DatabaseConfig
    Logging    LoggingConfig
    Environment string
    Debug      bool
}
```

### ServerConfig

サーバー関連の設定です。

```go
type ServerConfig struct {
    Port         int
    Host         string
    ReadTimeout  int
    WriteTimeout int
    IdleTimeout  int
}
```

### DatabaseConfig

データベース関連の設定です。

```go
type DatabaseConfig struct {
    Host            string
    Port            int
    Name            string
    User            string
    Password        string
    SSLMode         string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime int
}
```

### LoggingConfig

ログ関連の設定です。

```go
type LoggingConfig struct {
    Level         string
    Format        string
    Structured    bool
    IncludeCaller bool
}
```

## バリデーション

設定値は以下のルールでバリデーションされます：

- サーバーポート: 1-65535 の範囲
- データベースポート: 1-65535 の範囲
- 環境: development, staging, production のいずれか
- ログレベル: debug, info, warn, error のいずれか
- ログ形式: json, text のいずれか
- 必須フィールド: データベース名、ユーザー、パスワード

## テスト

```bash
go test ./pkg/config
```

## 例

詳細な使用例は `example/main.go` を参照してください。

```bash
go run pkg/config/example/main.go
```
