package model

import (
	"github.com/jinzhu/gorm"
	"sku-manage/mixin"
	"time"

	"github.com/sirupsen/logrus"
)

type Purchase struct {
	ID        uint32 `gorm:"primary_key" json:"id"`
	Name      string `json:"name"`
	Size      string `json:"size"`
	Sku       string `gorm:"type:varchar(64)" json:"sku"`
	Num       int    `json:"number"`
	RoleId    uint32 `json:"role_id"`
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

func GetPurchaseByRoleIdSku(sku string, roleId uint32) (Purchase, mixin.ErrorCode) {
	var p Purchase
	if err := db.Table("purchases").Where("sku = ? and role_id = ?", sku, roleId).Find(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return p, mixin.StatusOK
		}
		logrus.Errorf("[GetPurchaseByRoleIdSku] error %s", err.Error())
		return p, mixin.ErrorServerDb
	}
	return p, mixin.StatusOK
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

func GetRoleIdSkuNum() (map[uint32]SkuNum, mixin.ErrorCode) {
	resp := make(map[uint32]SkuNum)

	rows, err := db.Table("purchases").Select("role_id, count(sku), sum(num)").Group("role_id").Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, mixin.ErrorServerDb
	}

	for rows.Next() {
		var skuNum SkuNum
		var roleId uint32
		if err := rows.Scan(&roleId, &skuNum.SkuCount, &skuNum.Total); err != nil {
			logrus.Error(err.Error())
			return nil, mixin.ErrorServerDb
		}
		resp[roleId] = skuNum
	}

	return resp, mixin.StatusOK
}

// SELECT name, GROUP_CONCAT(sku), SUM(num) FROM purchases WHERE role_id = ? GROUP BY name;
func GetPurchaseByRoleId(roleId uint32) ([]Purchase, mixin.ErrorCode) {
	purchases := []Purchase{}

	rows, err := db.Table("purchases").
		Select("name, size, GROUP_CONCAT(sku), SUM(num)").
		Where("role_id = ?", roleId).
		Group("name, size").
		Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, mixin.ErrorServerDb
	}

	for rows.Next() {
		var purchase Purchase
		if err := rows.Scan(&purchase.Name, &purchase.Size, &purchase.Sku, &purchase.Num); err != nil {
			logrus.Error(err.Error())
			return nil, mixin.ErrorServerDb
		}
		purchases = append(purchases, purchase)
	}
	return purchases, mixin.StatusOK
}
