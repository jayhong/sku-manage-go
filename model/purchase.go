package model

import (
	"fmt"
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

	if err := db.Table("purchases").Where("order_id = ?", orderID).Find(&purchases).Error; err != nil {
		logrus.Errorf(err.Error())
		return nil, mixin.ErrorServerDb
	}
	return purchases, mixin.StatusOK
}

func DeletePurchaseByOrderId(roleId uint32) mixin.ErrorCode {
	if err := db.Where("order_id = ?", roleId).Delete(Purchase{}).Error; err != nil {
		logrus.Errorf("[DeletePurchaseRoleId] Delete error %s", err.Error())
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

func GetPurchaseSkusByOrderId(roleId uint32) ([]Purchase, mixin.ErrorCode) {
	purchases := []Purchase{}

	if err := db.Where("order_id = ?", roleId).Find(&purchases).Error; err != nil {
		logrus.Error(err.Error())
		return nil, mixin.ErrorServerDb
	}

	return purchases, mixin.StatusOK
}

type PurchasesItem struct {
	Sku         string `json:"sku"`
	Num         int    `json:"number"`
	SkuPropName string `json:"sku_prop_name"`
	ImageUrl    string `json:"image_url"`
	SizeName    string `json:"size_name"`
	Url         string `json:"url"`
	SkuNum      string `json:"sku_num"`
}

func GetOrderIdPurchases(orderID uint32) (map[string][]PurchasesItem, mixin.ErrorCode) {
	resp := make(map[string][]PurchasesItem)

	sqlStr := `SELECT p.sku, p.num, sku_props.img_url, sku_props.name, sizes.name, urls.url FROM purchases AS p INNER JOIN skus ON p.sku = skus.sku
INNER JOIN sku_props ON skus.sku_prop_id = sku_props.id
INNER JOIN sizes ON skus.size_id = sizes.id
INNER JOIN urls ON skus.url_id = urls.id where p.order_id = ?;`

	rows, err := db.Raw(sqlStr, orderID).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, mixin.ErrorServerDb
	}

	for rows.Next() {
		var item PurchasesItem
		if err := rows.Scan(&item.Sku, &item.Num,  &item.ImageUrl,&item.SkuPropName, &item.SizeName, &item.Url); err != nil {
			logrus.Error(err.Error())
			return nil, mixin.ErrorServerDb
		}

		item.SkuNum = fmt.Sprintf("%s 数量: %d", item.Sku, item.Num)
		resp[item.Url] = append(resp[item.Url], item)
	}

	return resp, mixin.StatusOK
}
