package apierror

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorCode string

type Descriptor struct {
	Code    ErrorCode
	Status  int
	Message string
}

type ErrorBody struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message,omitempty"`
	Detail  string    `json:"detail,omitempty"`
}

type errorEnvelope struct {
	Error ErrorBody `json:"error"`
}

type options struct {
	status  int
	message string
	detail  string
}

type Option func(*options)

func WithStatus(status int) Option {
	return func(o *options) {
		o.status = status
	}
}

func WithMessage(message string) Option {
	return func(o *options) {
		o.message = message
	}
}

func WithDetail(detail string) Option {
	return func(o *options) {
		o.detail = detail
	}
}

func Write(c *gin.Context, desc Descriptor, opts ...Option) {
	if desc.Status == 0 {
		desc.Status = http.StatusInternalServerError
	}

	cfg := options{}
	for _, opt := range opts {
		opt(&cfg)
	}

	status := desc.Status
	if cfg.status != 0 {
		status = cfg.status
	}

	message := desc.Message
	if cfg.message != "" {
		message = cfg.message
	}

	body := ErrorBody{
		Code: desc.Code,
	}
	if message != "" {
		body.Message = message
	}
	if cfg.detail != "" {
		body.Detail = cfg.detail
	}

	c.AbortWithStatusJSON(status, errorEnvelope{Error: body})
}

const (
	CodeInvalidRequest       ErrorCode = "INVALID_REQUEST"
	CodeInvalidCredentials   ErrorCode = "INVALID_CREDENTIALS"
	CodeAuthMissingToken     ErrorCode = "AUTH_MISSING_TOKEN"
	CodeAuthInvalidToken     ErrorCode = "AUTH_INVALID_TOKEN"
	CodeClusterAlreadyExists ErrorCode = "CLUSTER_ALREADY_EXISTS"
	CodeNodeNotFound         ErrorCode = "NODE_NOT_FOUND"
	CodeJobNotFound          ErrorCode = "JOB_NOT_FOUND"
	CodeJobAlreadyRunning    ErrorCode = "JOB_ALREADY_RUNNING"
	CodeJobNotRunning        ErrorCode = "JOB_NOT_RUNNING"
	CodeInternalError        ErrorCode = "INTERNAL_ERROR"
)

var (
	InvalidRequest = Descriptor{
		Code:    CodeInvalidRequest,
		Status:  http.StatusBadRequest,
		Message: "リクエスト内容が正しくありません。",
	}
	InvalidCredentials = Descriptor{
		Code:    CodeInvalidCredentials,
		Status:  http.StatusUnauthorized,
		Message: "クラスタIDまたはパスワードが正しくありません。",
	}
	AuthMissingToken = Descriptor{
		Code:    CodeAuthMissingToken,
		Status:  http.StatusUnauthorized,
		Message: "認証情報が見つかりません。",
	}
	AuthInvalidToken = Descriptor{
		Code:    CodeAuthInvalidToken,
		Status:  http.StatusUnauthorized,
		Message: "認証情報が無効です。",
	}
	ClusterAlreadyExists = Descriptor{
		Code:    CodeClusterAlreadyExists,
		Status:  http.StatusConflict,
		Message: "指定されたクラスタIDは既に使用されています。",
	}
	NodeNotFound = Descriptor{
		Code:    CodeNodeNotFound,
		Status:  http.StatusNotFound,
		Message: "ノードが見つかりません。",
	}
	JobNotFound = Descriptor{
		Code:    CodeJobNotFound,
		Status:  http.StatusNotFound,
		Message: "ジョブが見つかりません。",
	}
	JobAlreadyRunning = Descriptor{
		Code:    CodeJobAlreadyRunning,
		Status:  http.StatusConflict,
		Message: "このノードではすでにジョブが実行中です。",
	}
	JobNotRunning = Descriptor{
		Code:    CodeJobNotRunning,
		Status:  http.StatusConflict,
		Message: "実行中のジョブがありません。",
	}
	Internal = Descriptor{
		Code:    CodeInternalError,
		Status:  http.StatusInternalServerError,
		Message: "サーバーで問題が発生しました。時間をおいて再度お試しください。",
	}
)
