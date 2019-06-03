package model

import (
	"sku-manage/mixin"

	"github.com/sirupsen/logrus"
)

//alter  table  skus add unique index skus_index(sku,name);
type Sku struct {
	Sku  string `gorm:"type:varchar(64);primary_key"`
	Name string `gorm:"type:varchar(64);primary_key"`
	Size string `gorm:"type:varchar(64);primary_key"`
}

type SkuMapKey struct {
	Name string
	Size string
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

func DeleteSkuByNameSize(name, size string) mixin.ErrorCode {
	if err := db.Delete(Sku{}, "name = ? and size = ?", name, size).Error; err != nil {
		logrus.Errorf("[User.Delete] Delete error %s", err.Error())
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

func GetSkuMap() (map[SkuMapKey]string, mixin.ErrorCode) {
	skuMap := make(map[SkuMapKey]string)
	rows, err := db.Table("skus").Select("name, size, GROUP_CONCAT(sku)").Group("name, size").Rows()
	if err != nil {
		return skuMap, mixin.ErrorServerDb
	}
	var (
		name   string
		size   string
		skuStr string
	)
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&name, &size, &skuStr); err != nil {
			return skuMap, mixin.ErrorServerDb
		}
		skuMap[SkuMapKey{name, size}] = skuStr
	}

	return skuMap, mixin.StatusOK
}
