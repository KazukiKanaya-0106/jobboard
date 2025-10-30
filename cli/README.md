# jobboard-cli

Hub の node trigger API を呼び出しつつ、任意コマンドの実行結果を Slack に通知する CLI です。

## 使い方

```bash
jobboard [フラグ] -- <実行コマンド> [引数...]
```

例:

```bash
jobboard \
  --hub-url http://localhost:8080 \
  --node-token abc123 \
  --tag nightly \
  --timeout 90s \
  -- python scripts/import.py --force
```

## フラグと環境変数

| フラグ             | 環境変数                 | デフォルト              | 説明 |
|--------------------|--------------------------|-------------------------|------|
| `--hub-url`        | `JOBBOARD_HUB_URL`       | `http://localhost:8080` | Hub のベース URL |
| `--node-token`     | `JOBBOARD_NODE_TOKEN`    | なし                    | node trigger API 呼び出しに使用するトークン |
| `--tag`            | なし                     | なし                    | Hub に渡す任意タグ |
| `--slack-webhook`  | `JOBBOARD_SLACK_WEBHOOK` | なし (必須)             | Slack Incoming Webhook URL |
| `--timeout`        | なし                     | `60s`                   | Hub API 呼び出しのタイムアウト |

優先順位は **フラグ → 環境変数 → デフォルト** です。

## 動作フロー

1. フラグと環境変数を解決。Slack Webhook が未設定なら警告を出したうえでエラー終了します。
2. Slack Webhook が設定され、Node token が無い場合は `warning` を出し、Hub 連携をスキップして Slack 通知のみ行います。逆に Node token が設定され、Slack Webhook がない場合は `warning` を出し、Slack 通知をスキップして Hub 連携 のみ行い、どちらも設定されていない場合はエラーを出します。
3. Node token がある場合は Hub `/api/job-trigger/start` を呼び出し (`node_token`, `tag`, `started_at`)。
4. `--` 以降で指定したコマンドを実行し、終了コード・標準出力/標準エラーをそのまま返します。
5. Hub で開始済みの場合は `/api/job-trigger/finish` を呼び出し (`node_token`, `status`, `finished_at`, `duration_hours`)。
6. Slack Webhook へ結果を通知。本文にはコマンド、タグ、開始/終了時刻、所要時間、成功/失敗、エラー内容を含めます。

## エラーハンドリング

- Hub 呼び出しエラーや Slack 送信エラーは標準エラーに詳細を出し、CLI は実行コマンドの終了コードで終了します。
- Slack Webhook 未設定時は `warning` を表示後に即終了します。
- Node token 未設定時は `warning` を表示し、Hub 連携をスキップします。
- どちらも未設定の場合は `error` を表示し、終了する

## 留意点

- Slack 通知は Hub の成否にかかわらず必ず実行されます。
- CLI の終了コードは実行コマンドに追従するため、ワークフローから呼び出す際はコマンドの終了コードを確認してください。

