package model

import (
	"sku-manage/mixin"

	"github.com/sirupsen/logrus"
)

type Sku struct {
	UrlID     uint32 `json:"url_id"`
	SkuPropID uint32 `json:"sku_prop_id"`
	SizeID    uint32 `json:"size_id"`
	Sku       string `gorm:"type:varchar(64);primary_key" json:"sku"`
}

type SkuMapKey struct {
	SkuPropID uint32
	SizeID    uint32
}

func GetSku(skuStr string) (Sku, mixin.ErrorCode) {
	sku := Sku{}
	if err := db.Where("sku = ?", skuStr).First(&sku).Error; err != nil {
		return sku, mixin.ErrorServerDb
	}

	return sku, mixin.StatusOK
}

func CreateSku(sku Sku) mixin.ErrorCode {
	if err := db.Create(&sku).Error; err != nil {
		logrus.Errorf("[CreateSku] create error %s", err.Error())
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

func DeleteSku(sku string) mixin.ErrorCode {
	if err := db.Delete(Sku{}, "sku = ? ", sku).Error; err != nil {
		logrus.Errorf("[DeleteSku] Delete error %s", err.Error())
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

func DeleteSkuBySizeID(id uint32) mixin.ErrorCode {
	if err := db.Delete(Sku{}, "size_id = ? ", id).Error; err != nil {
		logrus.Errorf("[DeleteSkuBySizeID] Delete error %s", err.Error())
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

func DeleteSkuBySkuPropId(id uint32) mixin.ErrorCode {
	if err := db.Delete(Sku{}, "sku_prop_id = ? ", id).Error; err != nil {
		logrus.Errorf("[DeleteSkuBySkuPropId] Delete error %s", err.Error())
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

func GetUrlIdSkuMap(urlID uint32) (map[SkuMapKey]string, mixin.ErrorCode) {
	skuMap := make(map[SkuMapKey]string)
	rows, err := db.Table("skus").Select("sku_prop_id, size_id, GROUP_CONCAT(sku)").Where("url_id = ?", urlID).Group("sku_prop_id, size_id").Rows()
	if err != nil {
		return skuMap, mixin.ErrorServerDb
	}
	var (
		skuPropID uint32
		sizeID    uint32
		skuStr    string
	)
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&skuPropID, &sizeID, &skuStr); err != nil {
			return skuMap, mixin.ErrorServerDb
		}
		skuMap[SkuMapKey{skuPropID, sizeID}] = skuStr
	}

	return skuMap, mixin.StatusOK
}
