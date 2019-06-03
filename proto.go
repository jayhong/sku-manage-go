package main

import (
	"sku-manage/model"
)

type LoginState struct {
	UserId   uint32 `json:"user_id"`
	UserName string `json:"username"`
	Token    string `json:"token"`
}

type UserState struct {
	UserId   uint32 `json:"user_id"`
	UserName string `json:"username"`
}

type LoginRequest struct {
	UserName string `json:"username" valid:"required"`
	Password string `json:"password" valid:"required,length(6|100)"`
}

type LoginResponse struct {
	UserId     uint32        `json:"user_id" valid:"required"`
	UserName   string        `json:"username"`
	Permission []string      `json:"permission"`
	Role       model.Role    `json:"role"`
	Company    model.Company `json:"company"`
	Group      model.Group   `json:"group"`
	Department string        `json:"department"`
	IP         string        `json:"last_ip"`
	Time       int64         `json:"last_time"`
	Token      string        `json:"token"`
}

type UserInfoResponse struct {
	UserId     uint32           `json:"id"`
	UserName   string           `json:"username"`
	Ip         string           `json:"last_ip"`
	LastTime   int64            `json:"last_time"`
	Enable     int              `json:"enable"`
	Descript   string           `json:"descript"`
	CreatedAt  int64            `json:"created_at"`
	Role       model.Role       `json:"role"`
	Company    model.Company    `json:"company"`
	Group      model.Group      `json:"group"`
	Department model.Department `json:"department"`
}

type UserListResponse struct {
	UserId         uint32 `json:"id"`
	UserName       string `json:"username"`
	Ip             string `json:"last_ip"`
	LastTime       string `json:"last_time"`
	Enable         int    `json:"enable"`
	Descript       string `json:"descript"`
	CreatedAt      string `json:"created_at"`
	RoleID         uint32 `json:"role_id"`
	RoleName       string `json:"role_name"`
	CompanyID      uint32 `json:"company_id"`
	CompanyName    string `json:"company_name"`
	GroupID        uint32 `json:"group_id"`
	GroupName      string `json:"group_name"`
	DepartmentID   uint32 `json:"department_id"`
	DepartmentName string `json:"department_name"`
}

type UpdatePasswordRequest struct {
	UserName    string `json:"username",valid:"required"`
	NewPassword string `json:"new_password" valid:"required,length(6|100)"`
	OldPassword string `json:"old_password" valid:"required,length(6|100)"`
}

type ListResp struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}

type ListGroupResp struct {
	ID          uint32 `json:"group_id"`
	GroupName   string `json:"group_name"`
	CompanyId   uint32 `json:"company_id"`
	CompanyName string `json:"company_name"`
}

type ListRoleResp struct {
	ID       uint32   `json:"role_id"`
	RoleName string   `json:"role_name"`
	Descript string   `json:"descript"`
	Perm     []string `json:"permission"`
}
