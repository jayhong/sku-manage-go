package model

import (
	"sku-manage/mixin"

	"github.com/sirupsen/logrus"
)

type Group struct {
	ID        uint32 `gorm:"primary_key" json:"group_id"`
	GroupName string `json:"group_name"`
	LeaderID  uint32 `json:"leader_id"`
	CompanyId uint32 `json:"company_id"`
}

func CreateGroup(group Group) mixin.ErrorCode {
	if err := db.Create(&group).Error; err != nil {
		logrus.Errorf("[CreateGroup] error %s", err.Error())
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

func UpdateGroup(group Group) mixin.ErrorCode {
	if err := db.Save(&group).Error; err != nil {
		logrus.Errorf("[UpdateGroup] error %s", err.Error())
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func UpdateGroupLeader(groupId, leaderId uint32) mixin.ErrorCode {
	if err := db.Model(&Group{ID: groupId}).Update(map[string]interface{}{
		"leader_id": leaderId,
	}).Error; err != nil {
		logrus.Errorf("[UpdateGroupLeader] error %s", err.Error())
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func DeleteGroup(id uint32) mixin.ErrorCode {

	if err := db.Delete(&Group{ID: id}).Error; err != nil {
		logrus.Errorf("[DeleteGroup] error %s", err.Error())
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

func GetGroupList(param map[string]interface{}) ([]Group, mixin.ErrorCode) {
	var groups []Group

	if err := db.Where(param).Find(&groups).Error; err != nil {
		logrus.Errorf("[GetAllGroup] error %s", err.Error())
		return nil, mixin.ErrorServerDb
	}

	return groups, mixin.StatusOK
}

func GetAllGroupMap() (map[uint32]string, mixin.ErrorCode) {
	groupMap := make(map[uint32]string)

	var groups []Group
	if err := db.Find(&groups).Error; err != nil {
		logrus.Errorf("[GetAllDepartmentMap] error %s", err.Error)
		return nil, mixin.ErrorServerDb
	}

	for _, data := range groups {
		groupMap[data.ID] = data.GroupName
	}
	return groupMap, mixin.StatusOK
}

func GetCompanyGroup(companyId string) ([]Group, mixin.ErrorCode) {
	var groups []Group

	if err := db.Where("company_id = ?", companyId).Find(&groups).Error; err != nil {
		logrus.Errorf("[GetCompanyGroup] error %s", err.Error())
		return nil, mixin.ErrorServerDb
	}
	return groups, mixin.StatusOK
}

func GetCompanyGroupNum(companyId uint32) (uint, mixin.ErrorCode) {
	var count uint

	if err := db.Where("company_id = ?", companyId).Model(&Group{}).Count(&count).Error; err != nil {
		logrus.Errorf("[GetCompanyGroupNum] error %s", err.Error())
		return count, mixin.ErrorServerDb
	}
	return count, mixin.StatusOK
}

func GetGroup(id uint32) (Group, mixin.ErrorCode) {
	var group Group
	if err := db.First(&group, id).Error; err != nil {
		logrus.Errorf("[GetGroupName] error %s", err.Error())
	}

	return group, mixin.StatusOK
}

func GetGroupByName(name string) (Group, mixin.ErrorCode) {
	var group Group
	if err := db.Where("group_name = ?", name).First(&group).Error; err != nil {
		logrus.Errorf("[GetGroupName] error %s", err.Error())
	}

	return group, mixin.StatusOK
}
