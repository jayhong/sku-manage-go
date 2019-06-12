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

	perm := make([]string, 0)

	response := &LoginResponse{
		UserId:     user.ID,
		UserName:   inParam.UserName,
		Token:      token,
		Company:    company,
		Group:      group,
		IP:         user.Ip,
		Time:       user.LastTime,
		Permission: perm,
	}
	model.UpdateUser(model.User{ID: user.ID, Ip: r.RemoteAddr, LastTime: time.Now().Unix()})

	this.ResponseOK(w, response)
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
	inParam := map[string]interface{}{
		"user_id":    r.Form.Get("user_id"),
		"username":   r.Form.Get("username"),
		"company_id": r.Form.Get("company_id"),
		"group_id":   r.Form.Get("group_id"),
		"enable":     r.Form.Get("enable"),
	}

	for key, value := range inParam {
		if value == "" {
			delete(inParam, key)
		}
	}

	users, errCode := model.UserList(nil)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	var list []UserListResponse

	companyMap, _ := model.GetAllCompanyMap()
	groupMap, _ := model.GetAllGroupMap()

	for _, data := range users {
		list = append(list, UserListResponse{
			UserId:      data.ID,
			UserName:    data.UserName,
			Ip:          data.Ip,
			LastTime:    time.Unix(data.LastTime, 0).Format("2006-01-02"),
			Enable:      data.Enable,
			Descript:    data.Descript,
			CreatedAt:   data.CreatedAt.Format("2006-01-02"),
			CompanyID:   data.CompanyID,
			CompanyName: companyMap[data.CompanyID],
			GroupID:     data.GroupID,
			GroupName:   groupMap[data.GroupID],
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

	this.ResponseOK(w, UserInfoResponse{
		UserId:    user.ID,
		UserName:  user.UserName,
		Ip:        user.Ip,
		LastTime:  user.LastTime,
		Enable:    user.Enable,
		Descript:  user.Descript,
		CreatedAt: user.CreatedAt.Unix(),
		Company:   company,
		Group:     group,
	})
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

func (this *AccountService) list_company_handle(w http.ResponseWriter, r *http.Request) {
	inParam := map[string]interface{}{
		"id": r.Form.Get("company_id"),
	}

	for key, value := range inParam {
		if value == "" {
			delete(inParam, key)
		}
	}

	var companies []model.Company
	companies, errCode := model.GetAllCompany(inParam)
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

	groups, errCode := model.GetGroupList(inParam)
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

func (this *AccountService) dict_handle(w http.ResponseWriter, r *http.Request) {
	companies, errCode := model.GetAllCompany(nil)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	groups, errCode := model.GetGroupList(nil)
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

	this.ResponseOK(w, map[string]interface{}{
		"company": companies,
		"group":   groups,
	})
}

func (this *AccountService) tree_handle(w http.ResponseWriter, r *http.Request) {
	var tree []UserTree

	companies, errCode := model.GetAllCompany(nil)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	groups, errCode := model.GetGroupList(nil)
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

func (this *AccountService) add_order_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &model.Order{}

	if err := this.validator.Validate(r, inParam); err != nil {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}
	role, _ := model.GetOrderByName(inParam.OrderName)
	if role.ID != 0 {
		this.ResponseErrCode(w, mixin.ErrRoleNameExist)
		return
	}

	if errCode := model.CreateOrder(*inParam); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}
func (this *AccountService) update_order_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &model.Order{}

	if err := this.validator.Validate(r, inParam); err != nil {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	order, _ := model.GetOrderByName(inParam.OrderName)
	if order.ID != 0 && order.ID != inParam.ID {
		this.ResponseErrCode(w, mixin.ErrRoleNameExist)
		return
	}

	if errCode := model.UpdateOrder(*inParam); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}
func (this *AccountService) delete_order_handle(w http.ResponseWriter, r *http.Request) {
	inParam := make(map[string]uint32)
	if err := util.JsonDecode(r, &inParam); err != nil {
		logrus.Errorf("[AccountService.delete_order_handle] %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	if errCode := model.DeleteOrder(inParam["order_id"]); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	if errCode := model.DeletePurchaseByOrderId(inParam["order_id"]); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}
func (this *AccountService) list_order_handle(w http.ResponseWriter, r *http.Request) {
	orders, errCode := model.GetAllOrder()
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	for i, order := range orders{
		orders[i].CreateAtStr =  orders[i].CreatedAt.Format("2006-01-02 15:04:05")
		orders[i].SkuCount, orders[i].Total, errCode = model.GetPurchaseCountByOrderId(order.ID)
		if errCode != mixin.StatusOK {
			this.ResponseErrCode(w, errCode)
			return
		}
	}

	this.ResponseOK(w, orders)
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
	}

	this.ResponseOK(w, urls)
}
func (this *AccountService) add_url_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &model.URL{}

	if err := this.validator.Validate(r, inParam); err != nil {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	urlID, errCode := model.CreateURL(*inParam)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	sps, err := getSkuProps(inParam.Url)
	if err != nil {
		this.ResponseErrCode(w, mixin.ErrorServerDb)
		return
	}

	var skuProps []model.SkuProp
	var sizes []model.Size
	if len(sps) >= 1 {
		for j, _ := range sps[0].Value {
			skuProps = append(skuProps, model.SkuProp{
				UrlID:  urlID,
				Name:   sps[0].Value[j].Name,
				ImgUrl: sps[0].Value[j].Url,
			})
		}
		if len(skuProps) > 0 {
			errCode := model.BatchCreateSkuProp(skuProps)
			if errCode != mixin.StatusOK {
				this.ResponseOK(w, errCode)
				return
			}
		}
	}

	if len(sps) >= 2 {
		for i, _ := range sps[1].Value {
			sizes = append(sizes, model.Size{
				UrlID: urlID,
				Name:  sps[1].Value[i].Name,
			})
		}
		if len(sizes) > 0 {
			errCode := model.BatchCreateSize(sizes)
			if errCode != mixin.StatusOK {
				this.ResponseOK(w, errCode)
				return
			}
		}
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

	errCode = model.DeleteSizeByUrlID(inParam["url_id"])
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	errCode = model.DeleteSkuPropByUrlID(inParam["url_id"])
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	this.ResponseOK(w, nil)
}

// TODO 为了不被封ip需要模拟浏览器
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

func (this *AccountService) list_sku_props(w http.ResponseWriter, r *http.Request) {
	urlID := cast.ToUint32(r.Form.Get("url_id"))
	skuProps, errCode := model.GetSkuPropByUrlID(urlID)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, skuProps)
}
func (this *AccountService) add_sku_props(w http.ResponseWriter, r *http.Request) {
	inParam := &model.SkuProp{}
	if err := this.validator.Validate(r, inParam); err != nil {
		logrus.Errorf("[AccountService] add_sku_props %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	if _, errCode := model.CreateSkuProp(*inParam); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}
func (this *AccountService) update_sku_props(w http.ResponseWriter, r *http.Request) {
	inParam := &model.SkuProp{}
	if err := this.validator.Validate(r, inParam); err != nil || inParam.ID == 0 {
		logrus.Errorf("[AccountService] add_company_handle %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	if errCode := model.UpdateSkuProp(*inParam); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}
func (this *AccountService) delete_sku_props(w http.ResponseWriter, r *http.Request) {
	propID := cast.ToUint32(r.Form.Get("prop_id"))

	if errCode := model.DeleteSkuProp(propID); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	if errCode := model.DeleteSkuBySkuPropId(propID); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	this.ResponseOK(w, nil)
}

func (this *AccountService) list_size(w http.ResponseWriter, r *http.Request) {
	urlID := cast.ToUint32(r.Form.Get("url_id"))
	sizes, errCode := model.GetSizeByUrlID(urlID)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, sizes)
}
func (this *AccountService) add_size(w http.ResponseWriter, r *http.Request) {
	inParam := &model.Size{}
	if err := this.validator.Validate(r, inParam); err != nil {
		logrus.Errorf("[AccountService] add_size_props %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	if _, errCode := model.CreateSize(*inParam); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}
func (this *AccountService) update_size(w http.ResponseWriter, r *http.Request) {
	inParam := &model.Size{}
	if err := this.validator.Validate(r, inParam); err != nil || inParam.ID == 0 {
		logrus.Errorf("[AccountService] add_company_handle %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	if errCode := model.UpdateSize(*inParam); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}
func (this *AccountService) delete_size(w http.ResponseWriter, r *http.Request) {
	sizeID := cast.ToUint32(r.Form.Get("size_id"))

	if errCode := model.DeleteSize(sizeID); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	if errCode := model.DeleteSkuBySizeID(sizeID); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	this.ResponseOK(w, nil)
}

type SkusInfo struct {
	urlID     uint32   `json:"url_id"`
	SkuPropID uint32   `json:"sku_prop_id"`
	Name      string   `json:"name"`
	ImgUrl    string   `json:"image_url"`
	SizeID    uint32   `json:"size_id"`
	Size      string   `json:"size"`
	Skus      []string `json:"skus"`
}

func (this *AccountService) list_skus(w http.ResponseWriter, r *http.Request) {
	var listResp []SkusInfo

	urlID := cast.ToUint32(r.Form.Get("url_id"))
	skuProps, errCode := model.GetSkuPropByUrlID(urlID)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	sizes, errCode := model.GetSizeByUrlID(urlID)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	skuMap, errCode := model.GetUrlIdSkuMap(urlID)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	if len(sizes) == 0 {
		for i, _ := range skuProps {
			skus := make([]string, 0)
			if skuStr, ok := skuMap[model.SkuMapKey{SkuPropID: skuProps[i].ID}]; ok {
				skus = strings.Split(skuStr, ",")
			}

			listResp = append(listResp, SkusInfo{
				SkuPropID: skuProps[i].ID,
				Name:      skuProps[i].Name,
				ImgUrl:    skuProps[i].ImgUrl,
				Skus:      skus,
			})
		}
		this.ResponseOK(w, listResp)
		return
	}

	for i, _ := range skuProps {
		for j, _ := range sizes {
			skus := make([]string, 0)
			if skuStr, ok := skuMap[model.SkuMapKey{SkuPropID: skuProps[i].ID, SizeID: sizes[j].ID}]; ok {
				skus = strings.Split(skuStr, ",")
			}

			listResp = append(listResp, SkusInfo{
				SkuPropID: skuProps[i].ID,
				Name:      skuProps[i].Name,
				ImgUrl:    skuProps[i].ImgUrl,
				SizeID:    sizes[j].ID,
				Size:      sizes[j].Name,
				Skus:      skus,
			})
		}
	}

	this.ResponseOK(w, listResp)
}
func (this *AccountService) add_sku(w http.ResponseWriter, r *http.Request) {
	inParam := &model.Sku{}

	if err := this.validator.Validate(r, inParam); err != nil {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}
	if errCode := model.CreateSku(*inParam); errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	this.ResponseOK(w, nil)
}
func (this *AccountService) del_sku(w http.ResponseWriter, r *http.Request) {
	sku := r.Form.Get("sku")

	errCode := model.DeleteSku(sku)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	this.ResponseOK(w, nil)
}

type PurchaseListResp struct {
	Url          []string              `json:"url"`
	PurchaseList []model.PurchasesItem `json:"purchase_list"`
}

func (this *AccountService) purchase_list(w http.ResponseWriter, r *http.Request) {
	orderId := cast.ToUint32(r.Form.Get("order_id"))
	if orderId == 0 {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	var resp []PurchaseListResp

	purchaseList, errCode := model.GetOrderIdPurchases(orderId)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	for url, value := range purchaseList {
		resp = append(resp, PurchaseListResp{
			Url:          []string{url},
			PurchaseList: value,
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

	p, errCode := model.GetPurchaseByUrlIdSku(inParam.Sku, inParam.OrderId)
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
