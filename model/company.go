package model

import (
	"sku-manage/mixin"

	"github.com/sirupsen/logrus"
)

type Company struct {
	ID          uint32 `gorm:"primary_key" json:"company_id"`
	CompanyName string `json:"company_name"`
}

func CreateCompany(company Company) mixin.ErrorCode {
	if db.NewRecord(company) {
		if err := db.Create(&company).Error; err != nil {
			logrus.Errorf("[CreateCompany] error %s", err.Error())
			return mixin.ErrorServerDb
		}
	}

	return mixin.StatusOK
}

func UpdateCompany(company Company) mixin.ErrorCode {
	if err := db.Save(&company).Error; err != nil {
		logrus.Errorf("[UpdateCompany] error %s", err.Error())
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func DeleteCompany(company Company) mixin.ErrorCode {
	if err := db.Delete(&company).Error; err != nil {
		logrus.Errorf("[DeleteCompany] error %s", err.Error())
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func GetAllCompany(param map[string]interface{}) ([]Company, mixin.ErrorCode) {
	var companys []Company
	if err := db.Where(param).Find(&companys).Error; err != nil {
		logrus.Errorf("[GetAllCompany] error %s", err.Error)
		return nil, mixin.ErrorServerDb
	}
	return companys, mixin.StatusOK
}

func GetAllCompanyMap() (map[uint32]string, mixin.ErrorCode) {
	companyMap := make(map[uint32]string)

	var companys []Company
	if err := db.Find(&companys).Error; err != nil {
		logrus.Errorf("[GetAllCompany] error %s", err.Error)
		return nil, mixin.ErrorServerDb
	}

	for _, data := range companys {
		companyMap[data.ID] = data.CompanyName
	}
	return companyMap, mixin.StatusOK
}

func GetCompany(id uint32) (Company, mixin.ErrorCode) {
	var company Company

	if err := db.First(&company, id).Error; err != nil {
		logrus.Errorf("[GetCompanyName] error %s", err.Error())
	}
	return company, mixin.StatusOK
}

func GetCompanyByName(username string) (Company, mixin.ErrorCode) {
	var company Company

	if err := db.Where("company_name = ?", username).Find(&company).Error; err != nil {
		logrus.Errorf("[GetCompanyName] error %s", err.Error())
	}
	return company, mixin.StatusOK
}
