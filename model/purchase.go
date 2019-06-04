package model

import (
	"github.com/jinzhu/gorm"
	"sku-manage/mixin"
	"time"

	"github.com/sirupsen/logrus"
)

type Purchase struct {
	ID        uint32 `gorm:"primary_key" json:"id"`
	Sku       string `gorm:"type:varchar(64)" json:"sku"`
	Num       int    `json:"number"`
	OrderId   uint32 `json:"order_id"`
	CreatedAt time.Time
	UpdateAt  time.Time
}

func CreatePurchase(purchase Purchase) mixin.ErrorCode {
	if err := db.Create(&purchase).Error; err != nil {
		logrus.Errorf("[CreatePurchase] create error %s", err.Error())
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

func UpdatePurchaseNum(id uint32, num int) mixin.ErrorCode {
	if err := db.Model(Purchase{}).Where("id = ?", id).Update("num", num).Error; err != nil {
		logrus.Errorf("[UpdatePurchase] error %s", err.Error())
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func DeletePurchase(id uint32) mixin.ErrorCode {
	if err := db.Where("id = ?", id).Delete(Purchase{}).Error; err != nil {
		logrus.Errorf("[DeletePurchase] error %s", err.Error())
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func GetPurchaseByUrlIdSku(sku string, orderID uint32) (Purchase, mixin.ErrorCode) {
	var p Purchase
	if err := db.Table("purchases").Where("sku = ? and order_id = ?", sku, orderID).Find(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return p, mixin.StatusOK
		}
		logrus.Errorf("[GetPurchaseByUrlIdSku] error %s", err.Error())
		return p, mixin.ErrorServerDb
	}
	return p, mixin.StatusOK
}

func GetPurchaseByOrderId(orderID uint32) ([]Purchase, mixin.ErrorCode) {
	purchases := []Purchase{}

	if err := db.Table("purchasees").Where("order_id = ?", orderID).Find(&purchases).Error; err != nil {
		logrus.Errorf(err.Error())
		return nil, mixin.ErrorServerDb
	}
	return purchases, mixin.StatusOK
}

func DeletePurchaseByRoleId(roleId uint32) mixin.ErrorCode {
	if err := db.Where("role_id = ?", roleId).Delete(Purchase{}); err != nil {
		logrus.Errorf("[DeletePurchaseRoleId] Delete error %s", err.Error)
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

func GetPurchaseSkusByRoleId(roleId uint32) ([]Purchase, mixin.ErrorCode) {
	purchases := []Purchase{}

	if err := db.Where("role_id = ?", roleId).Find(&purchases).Error; err != nil {
		logrus.Error(err.Error())
		return nil, mixin.ErrorServerDb
	}

	return purchases, mixin.StatusOK
}

type SkuNum struct {
	SkuCount int
	Total    int
}

func GetOrderIdSkuNum() (map[uint32]SkuNum, mixin.ErrorCode) {
	resp := make(map[uint32]SkuNum)

	rows, err := db.Table("purchases").Select("order_id, count(sku), sum(num)").Group("order_id").Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, mixin.ErrorServerDb
	}

	for rows.Next() {
		var skuNum SkuNum
		var orderID uint32
		if err := rows.Scan(&orderID, &skuNum.SkuCount, &skuNum.Total); err != nil {
			logrus.Error(err.Error())
			return nil, mixin.ErrorServerDb
		}
		resp[orderID] = skuNum
	}

	return resp, mixin.StatusOK
}


