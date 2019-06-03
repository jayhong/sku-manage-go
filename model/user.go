package model

import (
	"fmt"
	"sku-manage/mixin"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint32    `gorm:"primary_key" json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	UserName     string    `gorm:"type:varchar(64);index;unique" json:"username" `
	Password     string    `json:"password"`
	RawPass      string    `json:"raw_password"`
	CompanyID    uint32    `json:"company_id"`
	GroupID      uint32    `json:"group_id"`
	RoleID       uint32    `json:"role_id"`
	DepartmentID uint32    `json:"department_id"`
	Enable       int       `json:"enable"`
	Descript     string    `json:"descript"`
	Ip           string    `json:"last_ip"`
	LastTime     int64     `json:"last_time"`
}

func CreateUser(user User) mixin.ErrorCode {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("[User.Create] bcrypt.GenerateFromPassword err: %s", err.Error())
		return mixin.ErrorServerCreateSecret
	}

	user.Password = string(hashedPassword)

	if err := db.Create(&user).Error; err != nil {
		logrus.Errorf("[User.Create] create record error %s", err.Error())
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

func DeleteUser(id uint32) mixin.ErrorCode {

	if err := db.Delete(&User{ID: id}).Error; err != nil {
		logrus.Errorf("[User.Delete] error %s", err.Error())
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

//修改一下
func UpdateUser(user User) mixin.ErrorCode {

	if err := db.Model(&user).Updates(user).Error; err != nil {
		logrus.Errorf("[User.Update] Updates err: %s", err.Error())
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

func UpdateUserRole(userId uint32) mixin.ErrorCode {

	if err := db.Model(&User{}).Where("id = ?", userId).Updates(map[string]interface{}{"role_id": 7}).Error; err != nil {
		logrus.Errorf("[User.UpdateUserRole] Updates err: %s", err.Error())
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

func ResetPassword(userID uint32) mixin.ErrorCode {
	user := User{
		ID: userID,
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("[User.Update] bcrypt.GenerateFromPassword err: %s", err.Error())
		return mixin.ErrorServerCreateSecret
	}

	if err := db.Model(&user).Update("password", hashedPassword).Error; err != nil {
		logrus.Errorf("[User.Update] update err: %s", err.Error())
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

func UpdatePassword(userName, oldPwd, newPwd string) mixin.ErrorCode {

	user, errCode := CheckPassword(userName, oldPwd)
	if errCode != mixin.StatusOK {
		return errCode
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("[User.Update] bcrypt.GenerateFromPassword err: %s", err.Error())
		return mixin.ErrorServerCreateSecret
	}

	if err := db.Model(&user).Update("password", hashedPassword).Error; err != nil {
		logrus.Errorf("[User.Update] update err: %s", err.Error())
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

func CheckPassword(userName, password string) (User, mixin.ErrorCode) {
	user, errCode := UserInfo(userName)
	if errCode != mixin.StatusOK {
		return user, errCode
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logrus.Errorf("[CheckPassword] bcrypt.CompareHashAndPassword %s, %s", user.Password, password)
		return user, mixin.ErrorClientUserOrPassword
	}
	return user, mixin.StatusOK
}

func ListUser(inParam map[string]interface{}) ([]User, mixin.ErrorCode) {
	var users []User
	var result *gorm.DB
	if username, ok := inParam["username"]; ok {
		delete(inParam, "username")
		result = db.Where("user_name like ?", "%"+username.(string)+"%").Where(inParam).Find(&users)
	} else {
		result = db.Where(inParam).Find(&users)
	}
	if result.RecordNotFound() {
		logrus.Debugf("[User.ListAll] error %s", result.Error.Error())
		return nil, mixin.ErrorServerDb
	}

	return users, mixin.StatusOK
}

func UserCount(countType string, id uint32) (uint, mixin.ErrorCode) {
	var count uint

	sqlStr := fmt.Sprintf("%s = ?", countType)
	if err := db.Where(sqlStr, id).Model(&User{}).Count(&count).Error; err != nil {
		logrus.Debugf("[User.UserNumByCompanyId] error %s", err.Error())
		return count, mixin.ErrorServerDb
	}

	return count, mixin.StatusOK

}

func UserInfoById(id uint32) (User, mixin.ErrorCode) {
	var user User
	result := db.Where("id = ?", id).First(&user)
	if err := result.Error; err != nil {
		if result.RecordNotFound() {
			return user, mixin.ErrorUserNotExist
		}
		logrus.Errorf("[User.Create] get user indo error %s", err.Error())
		return user, mixin.ErrorServerDb
	}
	return user, mixin.StatusOK
}

func UserInfo(userName string) (User, mixin.ErrorCode) {
	var user User
	result := db.Where("user_name = ?", userName).First(&user)
	if err := result.Error; err != nil {
		if result.RecordNotFound() {
			return user, mixin.ErrorUserNotExist
		}
		logrus.Errorf("[User.Create] get user indo error %s", err.Error())
		return user, mixin.ErrorServerDb
	}
	return user, mixin.StatusOK
}
