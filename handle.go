package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"sku-manage/mixin"
	"sku-manage/model"
	"sku-manage/util"
	"strconv"
	"strings"
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

func (this *AccountService) user(r *http.Request) (int32, string, bool) {
	token := r.Header.Get("X-Inspect-Token")
	if token == "" {
		return 0, "", false
	}
	userId, userName, _, ok := this._jwt.PublicJWT().Decode(token)
	if !ok {
		return 0, "", false
	}
	return cast.ToInt32(userId), userName, true
}

func (this *AccountService) login_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &LoginRequest{}
	if err := this.validator.Validate(r, inParam); err != nil {
		logrus.Errorf("[AccountService.login_handle] validate err: %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	user, errCode := model.CheckPassword(inParam.UserName, inParam.Password)
	if errCode != mixin.StatusOK {
		logrus.Errorf("[AccountService.login_handle] owner.Login name:%s, password:%s, err_code: %d", inParam.UserName, inParam.Password, errCode)
		this.ResponseErrCode(w, mixin.ErrorClientUserOrPassword)
		return
	}

	if user.Enable == 0 {
		this.ResponseErrCode(w, mixin.ErrorDisableUser)
		return
	}

	token, err := this._jwt.PublicJWT().Encode(fmt.Sprint(user.ID), user.UserName)
	if err != nil {
		logrus.Errorf("[AccountService.login_handle] gen token error %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorServerCreateToken)
		return
	}

	company, _ := model.GetCompany(user.CompanyID)
	group, _ := model.GetGroup(user.GroupID)
	role, _ := model.GetRole(user.RoleID)
	department, _ := model.GetDepartment(user.DepartmentID)

	perm := make([]string, 0)

	response := &LoginResponse{
		UserId:     user.ID,
		UserName:   inParam.UserName,
		Token:      token,
		Company:    company,
		Group:      group,
		Department: department.Name,
		IP:         user.Ip,
		Time:       user.LastTime,
		Role:       role,
		Permission: perm,
	}
	model.UpdateUser(model.User{ID: user.ID, Ip: r.RemoteAddr, LastTime: time.Now().Unix()})

	this.ResponseOK(w, response)
}

func (this *AccountService) create_user_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &model.User{}
	if err := this.validator.Validate(r, inParam); err != nil {
		logrus.Errorf(err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	tmp, _ := model.UserInfo(inParam.UserName)
	if tmp.ID != 0 {
		this.ResponseErrCode(w, mixin.ErrorUserOrEmailHasSignup)
		return
	}

	errCode := model.CreateUser(*inParam)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	if inParam.RoleID == 6 {
		tmp, _ := model.UserInfo(inParam.UserName)
		group, _ := model.GetGroup(inParam.GroupID)
		model.UpdateGroupLeader(group.ID, tmp.ID)
	}

	this.ResponseOK(w, nil)
}

func (this *AccountService) update_user_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &model.User{}
	if err := this.validator.Validate(r, inParam); err != nil || inParam.ID == 0 {
		//logrus.Errorf(err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	inParam.Password, inParam.RawPass = "", ""
	if inParam.UserName != "" {
		tmp, _ := model.UserInfo(inParam.UserName)
		if tmp.ID != 0 && tmp.UserName != inParam.UserName {
			this.ResponseErrCode(w, mixin.ErrorUserOrEmailHasSignup)
			return
		}
	}

	if inParam.RoleID == 6 {
		group, _ := model.GetGroup(inParam.GroupID)
		model.UpdateUserRole(group.LeaderID)
		model.UpdateGroupLeader(group.ID, inParam.ID)
	}

	errCode := model.UpdateUser(*inParam)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

func (this *AccountService) update_password_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &UpdatePasswordRequest{}
	if err := this.validator.Validate(r, inParam); err != nil {
		logrus.Errorf(err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	errCode := model.UpdatePassword(inParam.UserName, inParam.OldPassword, inParam.NewPassword)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	this.ResponseOK(w, map[string]interface{}{
		"errorcode": 0,
		"msg":       "success",
	})
}

func (this *AccountService) reset_password_handle(w http.ResponseWriter, r *http.Request) {
	idStr := r.Form.Get("user_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Errorf("[AccountService] del_company_handle %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}
	errCode := model.ResetPassword(uint32(id))
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	this.ResponseOK(w, nil)
}

func (this *AccountService) get_captcha_handle(w http.ResponseWriter, r *http.Request) {

	ConfigDigit := base64Captcha.ConfigDigit{
		Height:     80,
		Width:      240,
		CaptchaLen: 4,
		MaxSkew:    0.8,
		DotCount:   80,
	}

	//GenerateCaptcha 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
	captchaId, digitCap := base64Captcha.GenerateCaptcha("", ConfigDigit)
	base64Png := base64Captcha.CaptchaWriteToBase64Encoding(digitCap)

	body := map[string]interface{}{"errorcode": 0, "data": base64Png, "captcha_id": captchaId, "msg": "success"}
	this.ResponseOK(w, body)
}

func (this *AccountService) add_role_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &ListRoleResp{}

	if err := this.validator.Validate(r, inParam); err != nil {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}
	role, _ := model.GetRoleByName(inParam.RoleName)
	if role.ID != 0 {
		this.ResponseErrCode(w, mixin.ErrRoleNameExist)
		return
	}

	roleInfo := model.Role{
		RoleName: inParam.RoleName,
		Descript: inParam.Descript,
	}

	if errCode := model.CreateRole(roleInfo); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

func (this *AccountService) update_role_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &ListRoleResp{}

	if err := this.validator.Validate(r, inParam); err != nil {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	role, _ := model.GetRoleByName(inParam.RoleName)
	if role.ID != 0 && role.ID != inParam.ID {
		this.ResponseErrCode(w, mixin.ErrRoleNameExist)
		return
	}

	roleInfo := model.Role{
		ID:       inParam.ID,
		RoleName: inParam.RoleName,
		Descript: inParam.Descript,
	}

	if errCode := model.UpdateRole(roleInfo); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

//删除角色的时候需要检查一下是否有用户
func (this *AccountService) delete_role_handle(w http.ResponseWriter, r *http.Request) {

	inParam := make(map[string]uint32)
	if err := util.JsonDecode(r, &inParam); err != nil {
		logrus.Errorf("[AccountService.del_company_handle] %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	if errCode := model.DeleteRole(inParam["role_id"]); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	if errCode := model.DeletePurchaseByRoleId(inParam["role_id"]); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

func (this *AccountService) list_role_handle(w http.ResponseWriter, r *http.Request) {
	roles, errCode := model.GetAllRole()
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	skuCountMap, errCode := model.GetRoleIdSkuNum()
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	for i, _ := range roles {
		roles[i].SkuCount = skuCountMap[roles[i].ID].SkuCount
		roles[i].Total = skuCountMap[roles[i].ID].Total
	}

	this.ResponseOK(w, roles)
}

type PurchaseListResp struct {
	Name         string           `json:"name"`
	Size         string           `json:"size"`
	Sku          []string         `json:"sku"`
	Num          int              `json:"number"`
	ImageUrl     string           `json:"image_url"`
	PurchaseUrl  []string         `json:"purchase_url"`
	PurchaseList []model.Purchase `json:"purchase_list"`
}

func (this *AccountService) purchase_list(w http.ResponseWriter, r *http.Request) {
	roleIdStr := r.Form.Get("role_id")
	if roleIdStr == "" {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}
	roleId := cast.ToUint32(roleIdStr)

	ps, errCode := model.GetPurchaseSkusByRoleId(roleId)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	purMap := make(map[string]model.Purchase)
	for i, p := range ps {
		purMap[p.Sku] = ps[i]
	}

	purchases, errCode := model.GetPurchaseByRoleId(roleId)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	depMap, errCode := model.GetDepMapGroupByName()
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	var resp []PurchaseListResp
	for _, p := range purchases {
		key := model.SkuMapKey{p.Name, p.Size}
		var imgUrl string
		imgUrls := strings.Split(depMap[key].ImgUrl, ",")
		if len(imgUrls) > 0 {
			imgUrl = imgUrls[0]
		}

		skus := strings.Split(p.Sku, ",")
		var purchases []model.Purchase

		for i, _ := range skus {
			purchases = append(purchases, purMap[skus[i]])
		}

		resp = append(resp, PurchaseListResp{
			Name:         p.Name,
			Sku:          skus,
			Num:          p.Num,
			Size:         p.Size,
			ImageUrl:     imgUrl,
			PurchaseUrl:  strings.Split(depMap[key].PurchaseUrl, ","),
			PurchaseList: purchases,
		})
	}

	this.ResponseOK(w, resp)
}

func (this *AccountService) add_purchase(w http.ResponseWriter, r *http.Request) {
	inParam := &model.Purchase{}

	if err := this.validator.Validate(r, inParam); err != nil {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}
	var errCode mixin.ErrorCode
	sku, errCode := getSkuNameSize(inParam.Sku)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, mixin.ErrorNoSkus)
		return
	}

	inParam.Name = sku.Name
	inParam.Size = sku.Size

	// 判断sku是否存在如果存在就更新num
	p, errCode := model.GetPurchaseByRoleIdSku(inParam.Sku, inParam.RoleId)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	if p.ID != 0 {
		if errCode := model.UpdatePurchaseNum(p.ID, p.Num+inParam.Num); errCode != mixin.StatusOK {
			this.ResponseErrCode(w, errCode)
			return
		}
		this.ResponseOK(w, nil)
		return
	}

	if errCode := model.CreatePurchase(*inParam); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

func getSkuNameSize(skuStr string) (model.Sku, mixin.ErrorCode) {
	sku, errCode := model.GetSku(skuStr)
	if errCode != mixin.StatusOK {
		return model.Sku{}, errCode
	}
	return sku, mixin.StatusOK
}

func (this *AccountService) purchase_delete(w http.ResponseWriter, r *http.Request) {
	inParam := make(map[string]uint32)
	if err := util.JsonDecode(r, &inParam); err != nil {
		logrus.Errorf("[AccountService.del_company_handle] %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	if errCode := model.DeletePurchase(inParam["purchase_id"]); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

func (this *AccountService) purchase_update(w http.ResponseWriter, r *http.Request) {
	inParam := &model.Purchase{}

	if err := this.validator.Validate(r, inParam); err != nil {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	if errCode := model.UpdatePurchaseNum(inParam.ID, inParam.Num); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

type ConfigJsonBody struct {
	Id          string `json:"captcha_id"`
	VerifyValue string `json:"verify_value"`
}

func (this *AccountService) verify_captcha_handle(w http.ResponseWriter, r *http.Request) {

	param := &ConfigJsonBody{}
	if err := this.validator.Validate(r, param); err != nil {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	verifyResult := base64Captcha.VerifyCaptcha(param.Id, param.VerifyValue)

	body := map[string]interface{}{"errorcode": -1, "data": "验证失败", "msg": "captcha failed"}
	if verifyResult {
		body = map[string]interface{}{"errorcode": 0, "data": "验证通过", "msg": "captcha verified"}
	}

	this.ResponseOK(w, body)
}

func (this *AccountService) delete_user_handle(w http.ResponseWriter, r *http.Request) {

	inParam := make(map[string]uint32)
	if err := util.JsonDecode(r, &inParam); err != nil {
		logrus.Errorf("[AccountService.delete_user_handle] %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	if inParam["user_id"] == 0 {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	errCode := model.DeleteUser(inParam["user_id"])
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

func (this *AccountService) user_list_handle(w http.ResponseWriter, r *http.Request) {
	_, userName, _ := this.user(r)

	inParam := map[string]interface{}{
		"user_id":       r.Form.Get("user_id"),
		"username":      r.Form.Get("username"),
		"company_id":    r.Form.Get("company_id"),
		"group_id":      r.Form.Get("group_id"),
		"role_id":       r.Form.Get("role_id"),
		"department_id": r.Form.Get("department_id"),
		"enable":        r.Form.Get("enable"),
	}

	for key, value := range inParam {
		if value == "" {
			delete(inParam, key)
		}
	}

	userInfo, errCode := model.UserInfo(userName)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	role, errCode := model.GetRole(userInfo.RoleID)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	var users []model.User
	switch role.ID {
	case 1:
		users, errCode = model.ListUser(inParam)
	case 2:
		inParam["company_id"] = userInfo.CompanyID
		users, errCode = model.ListUser(inParam)
	default:
		errCode = mixin.ErrorClientUnauthorized
	}
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	var list []UserListResponse

	companyMap, _ := model.GetAllCompanyMap()
	groupMap, _ := model.GetAllGroupMap()
	roleMap, _ := model.GetAllRoleMap()
	departmentMap, _ := model.GetAllDepartmentMap()

	for _, data := range users {
		list = append(list, UserListResponse{
			UserId:         data.ID,
			UserName:       data.UserName,
			Ip:             data.Ip,
			LastTime:       time.Unix(data.LastTime, 0).Format("2006-01-02"),
			Enable:         data.Enable,
			Descript:       data.Descript,
			CreatedAt:      data.CreatedAt.Format("2006-01-02"),
			RoleID:         data.RoleID,
			RoleName:       roleMap[data.RoleID],
			CompanyID:      data.CompanyID,
			CompanyName:    companyMap[data.CompanyID],
			GroupID:        data.GroupID,
			GroupName:      groupMap[data.GroupID],
			DepartmentID:   data.DepartmentID,
			DepartmentName: departmentMap[data.DepartmentID],
		})
	}

	this.ResponseOK(w, list)
}

func (this *AccountService) user_info_handle(w http.ResponseWriter, r *http.Request) {
	userdID := r.Form.Get("user_id")
	if userdID == "" {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}
	user, errCode := model.UserInfoById(cast.ToUint32(userdID))
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	company, _ := model.GetCompany(user.CompanyID)
	group, _ := model.GetGroup(user.GroupID)
	role, _ := model.GetRole(user.RoleID)
	department, _ := model.GetDepartment(user.DepartmentID)

	this.ResponseOK(w, UserInfoResponse{
		UserId:     user.ID,
		UserName:   user.UserName,
		Ip:         user.Ip,
		LastTime:   user.LastTime,
		Enable:     user.Enable,
		Descript:   user.Descript,
		CreatedAt:  user.CreatedAt.Unix(),
		Role:       role,
		Company:    company,
		Group:      group,
		Department: department,
	})
}

func (this *AccountService) list_company_handle(w http.ResponseWriter, r *http.Request) {
	_, userName, _ := this.user(r)

	inParam := map[string]interface{}{
		"id": r.Form.Get("company_id"),
	}

	for key, value := range inParam {
		if value == "" {
			delete(inParam, key)
		}
	}

	userInfo, errCode := model.UserInfo(userName)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	role, errCode := model.GetRole(userInfo.RoleID)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	var companies []model.Company

	switch role.ID {
	case 1:
		companies, errCode = model.GetAllCompany(inParam)
	case 2:
		inParam := map[string]interface{}{
			"id": userInfo.CompanyID,
		}
		companies, errCode = model.GetAllCompany(inParam)
	default:
		this.ResponseErrCode(w, mixin.ErrorClientUnauthorized)
		return
	}
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	this.ResponseOK(w, companies)
}

func (this *AccountService) add_company_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &model.Company{}
	if err := this.validator.Validate(r, inParam); err != nil {
		logrus.Errorf("[AccountService] add_company_handle %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	tmp, _ := model.GetCompanyByName(inParam.CompanyName)
	if tmp.ID != 0 {
		this.ResponseErrCode(w, mixin.ErrorCompanyName)
		return
	}

	if errCode := model.CreateCompany(*inParam); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

func (this *AccountService) update_company_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &model.Company{}
	if err := this.validator.Validate(r, inParam); err != nil {
		logrus.Errorf("[AccountService] add_company_handle %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}
	tmp, _ := model.GetCompanyByName(inParam.CompanyName)
	if tmp.ID != 0 && tmp.ID != inParam.ID {
		this.ResponseErrCode(w, mixin.ErrorCompanyName)
		return
	}

	if errCode := model.UpdateCompany(*inParam); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

func (this *AccountService) del_company_handle(w http.ResponseWriter, r *http.Request) {

	inParam := make(map[string]uint32)
	if err := util.JsonDecode(r, &inParam); err != nil {
		logrus.Errorf("[AccountService.del_company_handle] %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	num, errCode := model.GetCompanyGroupNum(inParam["company_id"])
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, mixin.ErrorServerDb)
		return
	}
	if num != 0 {
		this.ResponseErrCode(w, mixin.ErrorCompanyHasGroup)
		return
	}

	if errCode := model.DeleteCompany(model.Company{ID: inParam["company_id"]}); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

func (this *AccountService) list_group_handle(w http.ResponseWriter, r *http.Request) {
	_, userName, _ := this.user(r)

	inParam := map[string]interface{}{
		"group_name": r.Form.Get("company_name"),
		"company_id": r.Form.Get("company_id"),
		"id":         r.Form.Get("group_id"),
	}

	for key, value := range inParam {
		if value == "" {
			delete(inParam, key)
		}
	}

	userInfo, errCode := model.UserInfo(userName)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	role, errCode := model.GetRole(userInfo.RoleID)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	var groups []model.Group

	switch role.ID {
	case 1:
		groups, errCode = model.GetGroupList(inParam)
	case 2:
		inParam["company_id"] = userInfo.CompanyID
		groups, errCode = model.GetGroupList(inParam)
	default:
		this.ResponseErrCode(w, mixin.ErrorClientUnauthorized)
		return
	}
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	var list []ListGroupResp
	companyMap, _ := model.GetAllCompanyMap()
	for _, data := range groups {
		list = append(list, ListGroupResp{
			ID:          data.ID,
			GroupName:   data.GroupName,
			CompanyId:   data.CompanyId,
			CompanyName: companyMap[data.CompanyId],
		})
	}

	this.ResponseOK(w, list)
}

func (this *AccountService) group_list(w http.ResponseWriter, r *http.Request) {

	inParam := map[string]interface{}{
		"group_name": r.Form.Get("company_name"),
		"company_id": r.Form.Get("company_id"),
	}

	for key, value := range inParam {
		if value == "" {
			delete(inParam, key)
		}
	}

	var errCode mixin.ErrorCode

	company_id := r.Form.Get("company_id")

	var groups []model.Group
	if company_id != "" {
		groups, errCode = model.GetGroupList(inParam)
	} else {
		inParam["company_id"] = company_id
		groups, errCode = model.GetGroupList(inParam)
	}
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	var lists []ListResp

	for _, data := range groups {
		lists = append(lists, ListResp{
			Name: data.GroupName,
			ID:   data.ID,
		})
	}

	this.ResponseOK(w, lists)
}

func (this *AccountService) update_group_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &model.Group{}
	if err := this.validator.Validate(r, inParam); err != nil {
		logrus.Errorf("[AccountService] update_group_handle %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}
	group, _ := model.GetGroupByName(inParam.GroupName)
	if group.ID != 0 && group.ID != inParam.ID {
		this.ResponseErrCode(w, mixin.ErrorGroupNameExist)
		return
	}
	if errCode := model.UpdateGroup(*inParam); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

func (this *AccountService) add_group_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &model.Group{}
	if err := this.validator.Validate(r, inParam); err != nil {
		logrus.Errorf("[AccountService] update_group_handle %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}
	group, _ := model.GetGroupByName(inParam.GroupName)
	if group.ID != 0 {
		this.ResponseErrCode(w, mixin.ErrorGroupNameExist)
		return
	}
	if errCode := model.UpdateGroup(*inParam); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

func (this *AccountService) del_group_handle(w http.ResponseWriter, r *http.Request) {
	inParam := make(map[string]uint32)
	if err := util.JsonDecode(r, &inParam); err != nil {
		logrus.Errorf("[AccountService.del_company_handle] %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	count, errCode := model.UserCount("group_id", inParam["group_id"])
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, mixin.ErrorServerDb)
		return
	}

	if count != 0 {
		this.ResponseErrCode(w, mixin.ErrorGroupHasUser)
		return
	}

	if inParam["group_id"] == 0 {
		this.ResponseOK(w, nil)
		return
	}

	if errCode := model.DeleteGroup(inParam["group_id"]); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

func (this *AccountService) list_department_handle(w http.ResponseWriter, r *http.Request) {
	deps, errCode := model.GetAllDepartment()
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	skuMap, _ := model.GetSkuMap()
	for i, _ := range deps {
		skus := strings.Split(skuMap[model.SkuMapKey{deps[i].Name, deps[i].Size}], ",")
		if len(skus) > 0 && skus[0] != "" {
			deps[i].Skus = skus
		}
	}

	this.ResponseOK(w, deps)
}

func (this *AccountService) add_department_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &model.Department{}

	if err := this.validator.Validate(r, inParam); err != nil {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	if inParam.Size == "" {
		this.ResponseErrCode(w, mixin.ErrorSizeParamError)
		return
	}

	// 判断链接是否存在
	url, errCode := model.GetUrlByUrl(inParam.PurchaseUrl)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	if url.ID == 0 {
		errCode = model.CreateURL(model.URL{
			Url:       inParam.PurchaseUrl,
			Status:    true,
			Collected: true,
		})
		if errCode != mixin.StatusOK {
			this.ResponseErrCode(w, errCode)
			return
		}
	}

	_, errCode = model.CreateDepartment(*inParam)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	for i, _ := range inParam.Skus {
		if errCode := model.CreateSku(model.Sku{
			Sku:  inParam.Skus[i],
			Name: inParam.Name,
			Size: inParam.Size,
		}); errCode != mixin.StatusOK {
			this.ResponseErrCode(w, errCode)
			return
		}
	}

	this.ResponseOK(w, nil)
}

type BatchAddDepartmentReq struct {
	Url   string              `json:"url"`
	Skus  []DownloadImgStruct `json:"skus"`
	Sizes []string            `json:"sizes"`
}

func (this *AccountService) batch_add_department_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &BatchAddDepartmentReq{}

	if err := this.validator.Validate(r, inParam); err != nil {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	// 判断链接是否存在
	url, errCode := model.GetUrlByUrl(inParam.Url)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	errCode = model.UpdateURLCollected(url.ID, true)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	for i, _ := range inParam.Skus {
		for j, _ := range inParam.Sizes {
			_, errCode = model.CreateDepartment(model.Department{
				Name:        inParam.Skus[i].Name,
				PurchaseUrl: inParam.Url,
				ImgUrl:      inParam.Skus[i].Url,
				Size:        inParam.Sizes[j],
			})
			if errCode != mixin.StatusOK {
				this.ResponseErrCode(w, errCode)
				return
			}
		}
	}

	this.ResponseOK(w, nil)
}

func (this *AccountService) update_department_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &model.Department{}

	if err := this.validator.Validate(r, inParam); err != nil {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	if !inParam.NameMerge {
		errCode := model.DeleteSkuByNameSize(inParam.OriginalName, inParam.OriginSize)
		if errCode != mixin.StatusOK {
			this.ResponseErrCode(w, errCode)
			return
		}
	}

	errCode := model.UpdateDepartment(*inParam)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	for i, _ := range inParam.Skus {
		if errCode := model.CreateSku(model.Sku{
			Sku:  inParam.Skus[i],
			Name: inParam.Name,
			Size: inParam.Size,
		}); errCode != mixin.StatusOK {
			continue
		}
	}

	this.ResponseOK(w, nil)
}

func (this *AccountService) del_department_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &model.Department{}

	if err := this.validator.Validate(r, inParam); err != nil {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}
	errCode := model.DeleteDepartment(inParam.ID)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	errCode = model.DeleteSkuByNameSize(inParam.Name, inParam.Size)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

func (this *AccountService) dict_handle(w http.ResponseWriter, r *http.Request) {

	_, userName, _ := this.user(r)

	userInfo, errCode := model.UserInfo(userName)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	role, errCode := model.GetRole(userInfo.RoleID)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	var companies []model.Company
	var groups []model.Group

	switch role.ID {
	case 1:
		companies, errCode = model.GetAllCompany(nil)
		groups, errCode = model.GetGroupList(nil)

	case 2:
		var company model.Company
		company, errCode = model.GetCompany(userInfo.CompanyID)
		companies = append(companies, company)
		groups, errCode = model.GetGroupList(map[string]interface{}{"company_id": userInfo.CompanyID})

	default:
		this.ResponseErrCode(w, mixin.ErrorClientUnauthorized)
		return
	}
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	var groupList []ListResp
	for _, data := range groups {
		groupList = append(groupList, ListResp{
			Name: data.GroupName,
			ID:   data.ID,
		})
	}

	deps, errCode := model.GetAllDepartment()
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	roles, errCode := model.GetAllRole()
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	var roleList []ListResp
	for _, data := range roles {
		roleList = append(roleList, ListResp{
			Name: data.RoleName,
			ID:   data.ID,
		})
	}
	this.ResponseOK(w, map[string]interface{}{
		"company":    companies,
		"group":      groups,
		"department": deps,
		"role":       roles,
	})
}

type UserTree struct {
	ID       uint32     `json:"id"`
	Lable    string     `json:"label"`
	Children []UserTree `json:"children"`
}

func (this *AccountService) tree_handle(w http.ResponseWriter, r *http.Request) {

	_, userName, _ := this.user(r)

	userInfo, errCode := model.UserInfo(userName)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	role, errCode := model.GetRole(userInfo.RoleID)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	var tree []UserTree
	var companies []model.Company
	var groups []model.Group
	switch role.ID {
	case 1:
		companies, errCode = model.GetAllCompany(nil)
		groups, errCode = model.GetGroupList(nil)
	case 2:
		var company model.Company
		company, errCode = model.GetCompany(userInfo.CompanyID)
		companies = append(companies, company)
		groups, errCode = model.GetGroupList(map[string]interface{}{"company_id": userInfo.CompanyID})
	default:
		this.ResponseErrCode(w, mixin.ErrorClientUnauthorized)
		return
	}
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	var groupList []ListResp
	for _, data := range groups {
		groupList = append(groupList, ListResp{
			Name: data.GroupName,
			ID:   data.ID,
		})
	}

	for _, c := range companies {
		tmp := []UserTree{}
		for _, g := range groups {
			if g.CompanyId == c.ID {
				tmp = append(tmp, UserTree{
					ID:       g.ID,
					Lable:    g.GroupName,
					Children: []UserTree{},
				})
			}
		}
		tree = append(tree, UserTree{
			ID:       c.ID,
			Lable:    c.CompanyName,
			Children: tmp,
		})
	}

	this.ResponseOK(w, map[string]interface{}{"userTree": tree})
}

func (this *AccountService) list_url_handle(w http.ResponseWriter, r *http.Request) {
	urls, errCode := model.GetAllURL()
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	for i, _ := range urls {
		if !urls[i].Status {
			urls[i].Type = "danger"
		}
		if !urls[i].Collected {
			urls[i].CollectedType = "danger"
		}
	}

	this.ResponseOK(w, urls)
}

func (this *AccountService) add_url_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &model.URL{}

	if err := this.validator.Validate(r, inParam); err != nil {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	errCode := model.CreateURL(*inParam)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	this.ResponseOK(w, nil)
}

func (this *AccountService) update_url_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &model.URL{}

	if err := this.validator.Validate(r, inParam); err != nil {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}
	if inParam.ID == 0 {
		this.ResponseOK(w, mixin.ErrorClientInvalidArgument)
		return
	}

	errCode := model.UpdateURL(*inParam)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	this.ResponseOK(w, nil)
}

func (this *AccountService) del_url_handle(w http.ResponseWriter, r *http.Request) {
	inParam := make(map[string]uint32)
	if err := util.JsonDecode(r, &inParam); err != nil {
		logrus.Errorf("[AccountService.del_url_handle] %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	if inParam["url_id"] == 0 {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	errCode := model.DeleteURL(inParam["url_id"])
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

func (this *AccountService) check_url_status(w http.ResponseWriter, r *http.Request) {
	urls, errCode := model.GetAllURL()
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	for i, url := range urls {
		t := time.Duration(rand.Intn(10) * 300)
		time.Sleep(t * time.Millisecond)
		resp, err := http.DefaultClient.Get(url.Url)
		if err != nil {
			logrus.Errorf(err.Error())
		}
		if err != nil || resp.StatusCode != 200 {
			model.UpdateURLStatus(urls[i].ID, false)
		}
	}

	this.ResponseOK(w, nil)
}

func (this *AccountService) getUrlSkus(w http.ResponseWriter, r *http.Request) {
	url := r.Form.Get("url")

	skuProps, err := getSkuProps(url)
	if err != nil {
		this.ResponseErrCode(w, mixin.ErrorServerDb)
		return
	}
	resp := make(map[string]interface{})
	switch len(skuProps) {
	case 1:
		resp["skus"] = skuProps[0]
	case 2:
		resp["skus"] = skuProps[0]

		var sizes []string
		for i, _ := range skuProps[1].Value {
			sizes = append(sizes, skuProps[1].Value[i].Name)
		}
		resp["sizes"] = sizes
	}

	this.ResponseOK(w, resp)
}

func (this *AccountService) download_url_pic(w http.ResponseWriter, r *http.Request) {
	url := r.Form.Get("url")
	fileName, err := DownloadImgs(url)
	if err != nil {
		this.ResponseErrCode(w, mixin.ErrorServerUnKnow)
		return
	}
	//需要把文件删除
	defer os.Remove(fileName)

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		this.ResponseErrCode(w, mixin.ErrorServerDb)
		return
	}

	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("content-disposition", "attachment;filename="+fileName)

	fmt.Fprint(w, string(content))
}
