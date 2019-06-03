package mixin

import (
	"net/http"
	"strconv"
)

type ErrorCode int32

var PearErrorMap map[ErrorCode]string

const (
	StatusOK                   ErrorCode = 0
	ErrorClientUnknow          ErrorCode = 4000
	ErrorClientUnauthorized    ErrorCode = 4001
	ErrorClientInvalidArgument ErrorCode = 4002
	ErrorClientSecretCheck     ErrorCode = 4003
	ErrorClientMacHasBind      ErrorCode = 4004
	ErrorClientUserOrPassword  ErrorCode = 4005
	ErrorUserOrEmailHasSignup  ErrorCode = 4006
	ErrorUserNotExist          ErrorCode = 4007
	ErrorNodeUnauthorized      ErrorCode = 4008
	ErrorSystemUnknow          ErrorCode = 4009
	ErrorProductUnknow         ErrorCode = 4010

	ErrorServerUnKnow       ErrorCode = 5000
	ErrorServerCreateSecret ErrorCode = 5001
	ErrorServerCreateToken  ErrorCode = 5002
	ErrorServerEncoding     ErrorCode = 5003
	ErrorServerDb           ErrorCode = 5004
	ErrorServerCache        ErrorCode = 5005
	ErrorServerRPC          ErrorCode = 5006
)

func init() {
	PearErrorMap = map[ErrorCode]string{
		ErrorClientUnknow:          "unknow error",
		ErrorClientUnauthorized:    "unauthorized",
		ErrorClientInvalidArgument: "invalid argument",
		ErrorClientSecretCheck:     "secret check error",
		ErrorClientMacHasBind:      "mac has bind",
		ErrorClientUserOrPassword:  "username or password error",
		ErrorUserOrEmailHasSignup:  "username or email has signup",
		ErrorUserNotExist:          "user not exist",
		ErrorNodeUnauthorized:      "node unauthorized",
		ErrorSystemUnknow:          "unknow system",
		ErrorProductUnknow:         "unknow product",

		ErrorServerUnKnow:       "unknow error",
		ErrorServerCreateSecret: "secret error",
		ErrorServerCreateToken:  "token error",
		ErrorServerEncoding:     "encoding error",
		ErrorServerDb:           "db error",
		ErrorServerCache:        "cache error",
		ErrorServerRPC:          "rpc error",
	}
}

type ErrorResponseFun func(http.ResponseWriter, int, int, string)

type PearErrorResponse interface {
	ErrorResponse(http.ResponseWriter, ErrorResponseFun)
}

type HttpServerError struct {
	code ErrorCode
}

func NewHttpServerError(code ErrorCode) *HttpServerError {
	return &HttpServerError{code: code}
}

func (this *HttpServerError) Error() string {
	errMsg := PearErrorMap[this.code]
	return errMsg
}

func (this *HttpServerError) ErrorResponse(w http.ResponseWriter, errFun ErrorResponseFun) {
	errCode, msg := pearGetCodeAndMsg(this.code, ErrorServerUnKnow)
	errFun(w, http.StatusInternalServerError, errCode, msg)
}

type HttpClientError struct {
	code ErrorCode
}

func NewHttpClientError(code ErrorCode) *HttpClientError {
	return &HttpClientError{code: code}
}

func (this *HttpClientError) Error() string {
	errMsg := PearErrorMap[this.code]
	return errMsg
}

func (this *HttpClientError) ErrorResponse(w http.ResponseWriter, errFun ErrorResponseFun) {
	errCode, msg := pearGetCodeAndMsg(this.code, ErrorClientUnknow)
	errFun(w, http.StatusBadRequest, errCode, msg)
}

func pearGetCodeAndMsg(errCode, defaultCode ErrorCode) (int, string) {
	if errMsg, ok := PearErrorMap[errCode]; ok {
		return int(errCode), errMsg
	}
	return int(defaultCode), PearErrorMap[defaultCode]
}

type PearError struct {
	Code string `json:"errorcode"`
	Msg  string `json:"msg"`
}

func NewPearError(errCode int, errMsg string) *PearError {
	return &PearError{
		Code: strconv.Itoa(errCode),
		Msg:  errMsg,
	}
}
