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
	ErrorGroupHasUser          ErrorCode = 4011
	ErrorCompanyHasGroup       ErrorCode = 4012
	ErrorDisableUser           ErrorCode = 4013
	ErrorRoleHasUser           ErrorCode = 4014
	ErrorDepartmentHasUser     ErrorCode = 4015
	ErrorNameDump              ErrorCode = 4016
	ErrorCompanyName           ErrorCode = 4017
	ErrRoleNameExist           ErrorCode = 4018
	ErrorGroupNameExist        ErrorCode = 4019
	ErrorDelBill               ErrorCode = 4020
	ErrorRoleNoExist           ErrorCode = 4021
	ErrorNoPermission          ErrorCode = 4022
	ErrorSizeParamError        ErrorCode = 4023

	ErrorServerUnKnow       ErrorCode = 5000
	ErrorServerCreateSecret ErrorCode = 5001
	ErrorServerCreateToken  ErrorCode = 5002
	ErrorServerEncoding     ErrorCode = 5003
	ErrorServerDb           ErrorCode = 5004
	ErrorServerCache        ErrorCode = 5005
	ErrorServerRPC          ErrorCode = 5006
	ErrorNoSkus             ErrorCode = 5007
)

func init() {
	PearErrorMap = map[ErrorCode]string{
		StatusOK:                   "成功",
		ErrorClientUnknow:          "未知错误",
		ErrorClientUnauthorized:    "登陆失效, 请重新登陆",
		ErrorClientInvalidArgument: "参数错误",
		ErrorClientSecretCheck:     "secret check error",
		ErrorClientMacHasBind:      "mac has bind",
		ErrorClientUserOrPassword:  "username or password error",
		ErrorUserOrEmailHasSignup:  "username or email has signup",
		ErrorUserNotExist:          "user not exist",
		ErrorNodeUnauthorized:      "node unauthorized",
		ErrorSystemUnknow:          "unknow system",
		ErrorProductUnknow:         "unknow product",
		ErrorGroupHasUser:          "该分组下有用户，请先更新用户的分组，再执行删除操作",
		ErrorCompanyHasGroup:       "该公司下有用户，请先更新用户的公司，再执行删除操作",
		ErrorDisableUser:           "您的账号被禁止登陆，请联系管理员",
		ErrorRoleHasUser:           "有该角色的用户，请先更新用户的角色，再执行删除操作",
		ErrorDepartmentHasUser:     "有该部门的用户，请先更新用户的部门，再执行删除操作",
		ErrorNameDump:              "该名称已存在，请修改名称后重试",
		ErrorCompanyName:           "公司名已存在",
		ErrRoleNameExist:           "角色名已存在",
		ErrorGroupNameExist:        "分组名已存在",
		ErrorDelBill:               "只能修改处于草稿中的进件",
		ErrorRoleNoExist:           "用户角色不存在",
		ErrorNoPermission:          "无权限",
		ErrorNoSkus:                "sku不存在",
		ErrorSizeParamError:        "款式不能为空",

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
