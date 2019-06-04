package model

import (
	"github.com/sirupsen/logrus"
	"sku-manage/mixin"
)

type Size struct {
	ID    uint32 `gorm:"primary_key" json:"size_id"`
	UrlID uint32 `json:"url_id"`
	Name  string `gorm:"type:varchar(64);" json:"name"`
}

func GetSize(id uint32) (Size, mixin.ErrorCode) {
	var size Size
	if err := db.First(&size, id).Error; err != nil {
		logrus.Errorf("[GetSize] error %s", err.Error)
		return size, mixin.ErrorServerDb
	}
	return size, mixin.StatusOK
}

func CreateSize(size Size) (uint32, mixin.ErrorCode) {
	if err := db.Create(&size).Error; err != nil {
		logrus.Errorf("[CreateSize] error %s", err.Error)
		return size.ID, mixin.ErrorServerDb
	}
	return size.ID, mixin.StatusOK
}

func BatchCreateSize(sizes []Size) (mixin.ErrorCode) {
	for _, size := range sizes {
		_, errCode := CreateSize(size)
		if errCode != mixin.StatusOK {
			return errCode
		}
	}

	return mixin.StatusOK
}

func UpdateSize(department Size) mixin.ErrorCode {
	if err := db.Save(&department).Error; err != nil {
		logrus.Errorf("[UpdateSize] error %s", err.Error)
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func DeleteSize(id uint32) mixin.ErrorCode {
	if err := db.Where("id = ?", id).Delete(Size{}).Error; err != nil {
		logrus.Errorf("[DeleteSize] error %s", err.Error)
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func DeleteSizeByUrlID(id uint32) mixin.ErrorCode {
	if err := db.Where("url_id = ?", id).Delete(Size{}).Error; err != nil {
		logrus.Errorf("[DeleteSize] error %s", err.Error)
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func GetSizeByUrlID(id uint32) ([]Size, mixin.ErrorCode) {
	var sizes []Size

	if err := db.Where("url_id = ?", id).Find(&sizes).Error; err != nil {
		logrus.Errorf("[GetSizeByUrlID] error %s", err.Error)
		return nil, mixin.ErrorServerDb
	}

	return sizes, mixin.StatusOK
}