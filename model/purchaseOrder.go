package model

import (
	"github.com/sirupsen/logrus"
	"sku-manage/mixin"
	"time"
)

//进货单
type Order struct {
	ID        uint32    `gorm:"primary_key" json:"order_id"`
	OrderName string    `json:"order_name"`
	Descript  string    `json:"descript"`
	SkuCount  int       `gorm:"-" json:"sku_count"`
	Total     int       `gorm:"-" json:"total"`
	UpdatedAt time.Time `json:"update_at"`
	CreatedAt time.Time `json:"create_at"`
}

func GetOrder(id uint32) (Order, mixin.ErrorCode) {
	var order Order

	result := db.First(&order, id)
	if err := result.Error; err != nil {
		if result.RecordNotFound() {
			return order, mixin.ErrorOrderNoExist
		}
		logrus.Errorf("[GetOrderName] error %s", err.Error)
		return order, mixin.ErrorServerDb
	}

	return order, mixin.StatusOK
}

func GetOrderByName(name string) (Order, mixin.ErrorCode) {
	var order Order

	result := db.Where("order_name = ?", name).First(&order)
	if err := result.Error; err != nil {
		if !result.RecordNotFound() {
			return order, mixin.ErrorOrderNoExist
		}
		logrus.Errorf("[GetOrderName] error %s", err.Error)
		return order, mixin.ErrorServerDb
	}

	return order, mixin.StatusOK
}

func CreateOrder(order Order) mixin.ErrorCode {
	if err := db.Create(&order).Error; err != nil {
		logrus.Errorf("[CreateOrder] error %s", err.Error)
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func UpdateOrder(order Order) mixin.ErrorCode {
	if err := db.Save(order).Error; err != nil {
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
	var orders []Order
	result := db.Find(&orders)
	if err := result.Error; err != nil {
		if result.RecordNotFound() {
			return nil, mixin.ErrorOrderNoExist
		}
		logrus.Errorf("[GetAllOrder] error %s", err.Error)
		return nil, mixin.ErrorServerDb
	}

	return orders, mixin.StatusOK
}

