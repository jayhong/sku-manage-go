package model

import (
	"sku-manage/mixin"

	"github.com/sirupsen/logrus"
)

type URL struct {
	ID     uint32 `gorm:"primary_key" json:"url_id"`
	Url    string `json:"url"`
	Status bool   `json:"status"`
	Type   string `gorm:"-" json:"type"` // 用于前端展示
}

func GetAllURL() ([]URL, mixin.ErrorCode) {
	var urls []URL

	if err := db.Find(&urls).Error; err != nil {
		logrus.Errorf("[GetAllURL] error %s", err.Error())
		return nil, mixin.ErrorServerDb
	}

	return urls, mixin.StatusOK
}

func CreateURL(url URL) (uint32, mixin.ErrorCode) {
	if err := db.Create(&url).Error; err != nil {
		logrus.Errorf("[CreateURL] create error %s", err.Error())
		return 0, mixin.ErrorServerDb
	}

	return url.ID, mixin.StatusOK
}

func UpdateURL(url URL) mixin.ErrorCode {
	if err := db.Save(url).Error; err != nil {
		logrus.Errorf("[UpdateRole] error %s", err.Error)
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func UpdateURLStatus(id uint32, status bool) mixin.ErrorCode {
	if err := db.Table("urls").Where("id = ?", id).Update("status", status).Error; err != nil {
		logrus.Errorf("[UpdateURLStatus] error %s", err.Error)
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func GetUrlByUrl(urlStr string) (URL, mixin.ErrorCode) {
	var url URL
	if err := db.Where("url = ?", urlStr).First(&url).Error; err != nil {
		logrus.Errorf("[GetUrlByUrl] error %s", err.Error)
		return url, mixin.ErrorServerDb
	}
	return url, mixin.StatusOK
}

func UpdateURLCollected(id uint32, collected bool) mixin.ErrorCode {
	if err := db.Table("urls").Where("id = ?", id).Update("collected", collected).Error; err != nil {
		logrus.Errorf("[UpdateURLCollected] error %s", err.Error)
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func DeleteURL(id uint32) mixin.ErrorCode {
	if err := db.Where("id = ?", id).Delete(URL{}).Error; err != nil {
		logrus.Errorf("[DeleteRole] error %s", err.Error)
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}
