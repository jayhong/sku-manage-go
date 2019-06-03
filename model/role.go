package model

import (
	"github.com/sirupsen/logrus"
	"sku-manage/mixin"
	"time"
)

//进货单
type Role struct {
	ID        uint32    `gorm:"primary_key" json:"role_id"`
	RoleName  string    `json:"role_name"`
	Descript  string    `json:"descript"`
	SkuCount  int       `gorm:"-" json:"sku_count"`
	Total     int       `gorm:"-" json:"total"`
	UpdatedAt time.Time `json:"update_at"`
	CreatedAt time.Time `json:"create_at"`
}

func GetRole(id uint32) (Role, mixin.ErrorCode) {
	var role Role

	result := db.First(&role, id)
	if err := result.Error; err != nil {
		if result.RecordNotFound() {
			return role, mixin.ErrorRoleNoExist
		}
		logrus.Errorf("[GetRoleName] error %s", err.Error)
		return role, mixin.ErrorServerDb
	}

	return role, mixin.StatusOK
}

func GetRoleByName(name string) (Role, mixin.ErrorCode) {
	var role Role

	result := db.Where("role_name = ?", name).First(&role)
	if err := result.Error; err != nil {
		if !result.RecordNotFound() {
			return role, mixin.ErrorRoleNoExist
		}
		logrus.Errorf("[GetRoleName] error %s", err.Error)
		return role, mixin.ErrorServerDb
	}

	return role, mixin.StatusOK
}

func CreateRole(role Role) mixin.ErrorCode {
	if err := db.Create(&role).Error; err != nil {
		logrus.Errorf("[CreateRole] error %s", err.Error)
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func UpdateRole(role Role) mixin.ErrorCode {
	if err := db.Save(role).Error; err != nil {
		logrus.Errorf("[UpdateRole] error %s", err.Error)
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func DeleteRole(id uint32) mixin.ErrorCode {
	if err := db.Where("id = ?", id).Delete(Role{}).Error; err != nil {
		logrus.Errorf("[DeleteRole] error %s", err.Error)
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func GetAllRole() ([]Role, mixin.ErrorCode) {
	var roles []Role
	result := db.Find(&roles)
	if err := result.Error; err != nil {
		if result.RecordNotFound() {
			return nil, mixin.ErrorRoleNoExist
		}
		logrus.Errorf("[GetAllRole] error %s", err.Error)
		return nil, mixin.ErrorServerDb
	}

	return roles, mixin.StatusOK
}

func GetAllRoleMap() (map[uint32]string, mixin.ErrorCode) {
	roleMap := make(map[uint32]string)

	var roles []Role
	if err := db.Find(&roles).Error; err != nil {
		logrus.Errorf("[GetAllRoleMap] error %s", err.Error)
		return nil, mixin.ErrorServerDb
	}

	for _, data := range roles {
		roleMap[data.ID] = data.RoleName
	}
	return roleMap, mixin.StatusOK
}
