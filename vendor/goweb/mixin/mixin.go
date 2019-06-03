package mixin

import (
	"net/http"

	"github.com/unrolled/render"
)

type ResponseMixin struct {
	render *render.Render
}

func NewResponseMixin() *ResponseMixin {
	return &ResponseMixin{
		render: render.New(render.Options{IndentJSON: true}),
	}
}

func (this *ResponseMixin) ResponseOK(w http.ResponseWriter, value interface{}) {
	this.render.JSON(w, http.StatusOK, value)
}

func (this *ResponseMixin) ResponseError(w http.ResponseWriter, err error) {
	if errResponse, ok := err.(PearErrorResponse); ok {
		errResponse.ErrorResponse(w, this.renderError)
		return
	}
	NewHttpServerError(ErrorServerUnKnow).ErrorResponse(w, this.renderError)
}

func (this *ResponseMixin) renderError(w http.ResponseWriter, status int, errCode int, errMsg string) {
	this.render.JSON(w, status, NewPearError(errCode, errMsg))
}

func (this *ResponseMixin) ResponseErrCode(w http.ResponseWriter, errCode ErrorCode) {
	if errCode < ErrorServerUnKnow {
		this.ResponseError(w, NewHttpClientError(errCode))
		return
	}
	this.ResponseError(w, NewHttpServerError(errCode))
}

func (this *ResponseMixin) ResponseUnauthorized(w http.ResponseWriter) {
	this.renderError(w, http.StatusUnauthorized, int(ErrorClientUnauthorized), "unauthorized")
}
