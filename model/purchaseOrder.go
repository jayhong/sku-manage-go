package model

import (
	"github.com/sirupsen/logrus"
	"sku-manage/mixin"
	"time"
)

//进货单
type Order struct {
	ID        uint32    `gorm:"primary_key" json:"role_id"`
	OrderName string    `json:"role_name"`
	Descript  string    `json:"descript"`
	SkuCount  int       `gorm:"-" json:"sku_count"`
	Total     int       `gorm:"-" json:"total"`
	UpdatedAt time.Time `json:"update_at"`
	CreatedAt time.Time `json:"create_at"`
}

func GetOrder(id uint32) (Order, mixin.ErrorCode) {
	var role Order

	result := db.First(&role, id)
	if err := result.Error; err != nil {
		if result.RecordNotFound() {
			return role, mixin.ErrorOrderNoExist
		}
		logrus.Errorf("[GetOrderName] error %s", err.Error)
		return role, mixin.ErrorServerDb
	}

	return role, mixin.StatusOK
}

func GetOrderByName(name string) (Order, mixin.ErrorCode) {
	var order Order

	result := db.Where("role_name = ?", name).First(&order)
	if err := result.Error; err != nil {
		if !result.RecordNotFound() {
			return order, mixin.ErrorOrderNoExist
		}
		logrus.Errorf("[GetOrderName] error %s", err.Error)
		return order, mixin.ErrorServerDb
	}

	return order, mixin.StatusOK
}

func CreateOrder(role Order) mixin.ErrorCode {
	if err := db.Create(&role).Error; err != nil {
		logrus.Errorf("[CreateOrder] error %s", err.Error)
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func UpdateOrder(role Order) mixin.ErrorCode {
	if err := db.Save(role).Error; err != nil {
		logrus.Errorf("[UpdateOrder] error %s", err.Error)
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func DeleteOrder(id uint32) mixin.ErrorCode {
	if err := db.Where("id = ?", id).Delete(Order{}).Error; err != nil {
		logrus.Errorf("[DeleteOrder] error %s", err.Error)
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func GetAllOrder() ([]Order, mixin.ErrorCode) {
	var roles []Order
	result := db.Find(&roles)
	if err := result.Error; err != nil {
		if result.RecordNotFound() {
			return nil, mixin.ErrorOrderNoExist
		}
		logrus.Errorf("[GetAllOrder] error %s", err.Error)
		return nil, mixin.ErrorServerDb
	}

	return roles, mixin.StatusOK
}

