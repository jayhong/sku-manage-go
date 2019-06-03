package main

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"sku-manage/server"
)

func (s *AccountService) RegisterRoutes(router *mux.Router, prefix string) {
	subRouter := router.PathPrefix(prefix).Subrouter()
	server.AddRoutes(s.GetRoutes(), subRouter)
}

func (s *AccountService) GetRoutes() []server.Route {
	return []server.Route{
		server.Route{
			Name:        "login",
			Method:      "POST",
			Pattern:     "/login",
			HandlerFunc: s.login_handle,
		},

		server.Route{
			Name:        "reset password",
			Method:      "GET",
			Pattern:     "/{user_id:[0-9]+}/password/reset",
			HandlerFunc: s.reset_password_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "Captcha",
			Method:      "GET",
			Pattern:     "/get_captcha",
			HandlerFunc: s.get_captcha_handle,
		},

		server.Route{
			Name:        "Captcha",
			Method:      "POST",
			Pattern:     "/verify_captcha",
			HandlerFunc: s.verify_captcha_handle,
		},

		server.Route{
			Name:        "update password",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/password/update",
			HandlerFunc: s.update_password_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "create user",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/user/create",
			HandlerFunc: s.create_user_handle,
			//			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "update user",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/user/update",
			HandlerFunc: s.update_user_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "delete user",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/user/delete",
			HandlerFunc: s.delete_user_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "user list",
			Method:      "GET",
			Pattern:     "/{user_id:[0-9]+}/user/list",
			HandlerFunc: s.user_list_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "user info",
			Method:      "GET",
			Pattern:     "/{user_id:[0-9]+}/user/info",
			HandlerFunc: s.user_info_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "user info",
			Method:      "GET",
			Pattern:     "/user/info",
			HandlerFunc: s.user_info_handle,
		},

		server.Route{
			Name:        "list company",
			Method:      "GET",
			Pattern:     "/{user_id:[0-9]+}/company",
			HandlerFunc: s.list_company_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "add company",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/company/add",
			HandlerFunc: s.add_company_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "update company",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/company/update",
			HandlerFunc: s.update_company_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "delete company",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/company/del",
			HandlerFunc: s.del_company_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "add group",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/group/add",
			HandlerFunc: s.add_group_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},
		server.Route{
			Name:        "update group",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/group/update",
			HandlerFunc: s.update_group_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},
		server.Route{
			Name:        "delete group",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/group/del",
			HandlerFunc: s.del_group_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},
		server.Route{
			Name:        "list group",
			Method:      "GET",
			Pattern:     "/{user_id:[0-9]+}/group",
			HandlerFunc: s.list_group_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "group list",
			Method:      "GET",
			Pattern:     "/{user_id:[0-9]+}/group/list",
			HandlerFunc: s.group_list,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "dict",
			Method:      "GET",
			Pattern:     "/{user_id:[0-9]+}/dict",
			HandlerFunc: s.dict_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "tree",
			Method:      "GET",
			Pattern:     "/{user_id:[0-9]+}/tree",
			HandlerFunc: s.tree_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},


		//PC上传文件接口
		server.Route{
			Name:        "upload file",
			Method:      "POST",
			Pattern:     "/upload/file",
			HandlerFunc: s.upload_file_handle,
		},

		// 款式管理
		server.Route{
			Name:        "list department",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/department",
			HandlerFunc: s.list_department_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},
		server.Route{
			Name:        "add department",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/department/add",
			HandlerFunc: s.add_department_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},
		server.Route{
			Name:        "add departments",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/department/mutiadd",
			HandlerFunc: s.batch_add_department_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},
		server.Route{
			Name:        "department",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/department/update",
			HandlerFunc: s.update_department_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},
		server.Route{
			Name:        "del department",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/department/del",
			HandlerFunc: s.del_department_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "add role",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/role/add",
			HandlerFunc: s.add_role_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},
		server.Route{
			Name:        "update role",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/role/update",
			HandlerFunc: s.update_role_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},
		server.Route{
			Name:        "delete role",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/role/del",
			HandlerFunc: s.delete_role_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},
		server.Route{
			Name:        "list role",
			Method:      "GET",
			Pattern:     "/{user_id:[0-9]+}/role",
			HandlerFunc: s.list_role_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "add purchase",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/purchase/add",
			HandlerFunc: s.add_purchase,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "list purchase",
			Method:      "GET",
			Pattern:     "/{user_id:[0-9]+}/purchase/list",
			HandlerFunc: s.purchase_list,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},


		server.Route{
			Name:        "delete purchase",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/purchase/delete",
			HandlerFunc: s.purchase_delete,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "update purchase",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/purchase/update",
			HandlerFunc: s.purchase_update,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},
		server.Route{
			Name:        "add url",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/url/add",
			HandlerFunc: s.add_url_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "list url",
			Method:      "GET",
			Pattern:     "/{user_id:[0-9]+}/url/list",
			HandlerFunc: s.list_url_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "delete url",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/url/delete",
			HandlerFunc: s.del_url_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "update url ",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/url/update",
			HandlerFunc: s.update_url_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},


		server.Route{
			Name:        "Check url status",
			Method:      "GET",
			Pattern:     "/{user_id:[0-9]+}/url/status",
			HandlerFunc: s.check_url_status,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},
		server.Route{
			Name:        "get url sku",
			Method:      "GET",
			Pattern:     "/{user_id:[0-9]+}/url/sku",
			HandlerFunc: s.getUrlSkus,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},
		server.Route{
			Name:        "downloadurl",
			Method:      "GET",
			Pattern:     "/{user_id:[0-9]+}/url/download",
			HandlerFunc: s.download_url_pic,
			//Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},
	}
}
