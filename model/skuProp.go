package model

import (
	"github.com/sirupsen/logrus"
	"sku-manage/mixin"
)

type SkuProp struct {
	ID     uint32 `gorm:"primary_key" json:"sku_prop_id"`
	UrlID  uint32 `json:"url_id"`
	Name   string `gorm:"type:varchar(64);" json:"name"`
	ImgUrl string `json:"image_url"`
}

func GetSkuProp(id uint32) (SkuProp, mixin.ErrorCode) {
	var skuProp SkuProp
	if err := db.First(&skuProp, id).Error; err != nil {
		logrus.Errorf("[GetSkuProp] error %s", err.Error)
		return skuProp, mixin.ErrorServerDb
	}
	return skuProp, mixin.StatusOK
}

func GetSkuPropByUrlID(id uint32) ([]SkuProp, mixin.ErrorCode) {
	var skuProp []SkuProp

	if err := db.Where("url_id = ?", id).Find(&skuProp).Error; err != nil {
		logrus.Errorf("[GetSkuPropByUrlID] error %s", err.Error)
		return nil, mixin.ErrorServerDb
	}

	return skuProp, mixin.StatusOK
}

func CreateSkuProp(skuProp SkuProp) (uint32, mixin.ErrorCode) {
	if err := db.Create(&skuProp).Error; err != nil {
		logrus.Errorf("[CreateSkuProp] error %s", err.Error)
		return skuProp.ID, mixin.ErrorServerDb
	}
	return skuProp.ID, mixin.StatusOK
}

func BatchCreateSkuProp(skuProps []SkuProp) (mixin.ErrorCode) {
	for _, skuProp := range skuProps {
		_, errCode := CreateSkuProp(skuProp)
		if errCode != mixin.StatusOK {
			return errCode
		}
	}

	return mixin.StatusOK
}

func UpdateSkuProp(skuProps SkuProp) mixin.ErrorCode {
	if err := db.Save(&skuProps).Error; err != nil {
		logrus.Errorf("[UpdateSkuProp] error %s", err.Error)
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func DeleteSkuProp(id uint32) mixin.ErrorCode {
	if err := db.Where("id = ?", id).Delete(SkuProp{}).Error; err != nil {
		logrus.Errorf("[DeleteSkuProp] error %s", err.Error)
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func DeleteSkuPropByUrlID(id uint32) mixin.ErrorCode {
	if err := db.Where("url_id = ?", id).Delete(SkuProp{}).Error; err != nil {
		logrus.Errorf("[DeleteSkuPropByUrlID] error %s", err.Error)
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}
