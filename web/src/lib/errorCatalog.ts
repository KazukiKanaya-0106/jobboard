import { ApiError } from "./apiCient";

const ERROR_MESSAGES: Record<string, string> = {
  INVALID_REQUEST: "入力内容を確認してください。",
  INVALID_CREDENTIALS: "クラスタIDまたはパスワードが正しくありません。",
  AUTH_MISSING_TOKEN: "認証情報が見つかりません。再度ログインしてください。",
  AUTH_INVALID_TOKEN: "セッションの有効期限が切れたか認証情報が不正です。",
  CLUSTER_ALREADY_EXISTS: "指定したクラスタIDは既に使用されています。",
  NODE_NOT_FOUND: "対象のノードが見つかりません。",
  JOB_NOT_FOUND: "対象のジョブが見つかりません。",
  JOB_ALREADY_RUNNING: "このノードではすでにジョブが実行中です。",
  JOB_NOT_RUNNING: "実行中のジョブはありません。",
  INTERNAL_ERROR: "サーバーで問題が発生しました。時間をおいて再度お試しください。",
};

export function resolveErrorMessage(error: unknown, fallback = "予期せぬエラーが発生しました。"): string {
  if (error instanceof ApiError) {
    if (error.code && ERROR_MESSAGES[error.code]) {
      return ERROR_MESSAGES[error.code];
    }
    if (error.message) {
      return error.message;
    }
  }

  if (error instanceof Error && error.message.trim() !== "") {
    return error.message;
  }

  return fallback;
}

export function resolveErrorCode(error: unknown): string | undefined {
  if (error instanceof ApiError) {
    return error.code;
  }
  return undefined;
}
