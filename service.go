package main

import (
	"net/http"

	"sku-manage/mixin"
)

type AccountInterface interface {
	Signup(username, password, email string) (*LoginState, mixin.ErrorCode)
	Login(username, password string) (*LoginState, mixin.ErrorCode)
	UpdatePassword(userName, newPassword, oldPassword string) mixin.ErrorCode
	BindPhone(userName, phone string) mixin.ErrorCode
	BindEmail(userName, email string) mixin.ErrorCode
	CheckEmail(userName, email string) (bool, mixin.ErrorCode)
	CheckPassword(userName, password string) mixin.ErrorCode
	ForgetPassword(userName, email string) mixin.ErrorCode
	BindAccount(username, password, openId string) mixin.ErrorCode
	Login3party(state, code string) (openid string, err error)
	Login3partyCheck(openid, state string) (*LoginState, mixin.ErrorCode)
}

type RequestValidator interface {
	Validate(*http.Request, interface{}) error
}

type AccountService struct {
	*mixin.ResponseMixin
	validator RequestValidator
	_jwt      JWT
}

func NewAccountService(validator RequestValidator,
	jwt JWT) *AccountService {
	return &AccountService{
		ResponseMixin: mixin.NewResponseMixin(),
		validator:     validator,
		_jwt:          jwt,
	}
}
