package model

import (
	"github.com/sirupsen/logrus"
	"sku-manage/mixin"
)

type Department struct {
	ID           uint32   `gorm:"primary_key" json:"department_id"`
	Name         string   `gorm:"type:varchar(64);" json:"department"`
	Size         string   `gorm:"type:varchar(64);" json:"size"`
	PurchaseUrl  string   `gorm:"purchase_url" json:"purchase_url"`
	ImgUrl       string   `json:"image_url"`
	Skus         []string `gorm:"-" json:"skus"`
	Sizes        []string `gorm:"-" json:"sizes"` // 用于批量添加
	NameMerge    bool     `gorm:"-" json:"name_merge"`
	OriginSize   string   `gorm:"-" json:"original_size"`
	OriginalName string   `gorm:"-" json:"original_name"`
}

func GetDepartment(id uint32) (Department, mixin.ErrorCode) {
	var department Department
	if err := db.First(&department, id).Error; err != nil {
		logrus.Errorf("[GetDepartment] error %s", err.Error)
		return department, mixin.ErrorServerDb
	}
	return department, mixin.StatusOK
}

func CreateDepartment(department Department) (uint32, mixin.ErrorCode) {
	if err := db.Create(&department).Error; err != nil {
		logrus.Errorf("[CreateDepartment] error %s", err.Error)
		return department.ID, mixin.ErrorServerDb
	}
	return department.ID, mixin.StatusOK
}

func UpdateDepartment(department Department) mixin.ErrorCode {
	if err := db.Save(&department).Error; err != nil {
		logrus.Errorf("[UpdateDepartment] error %s", err.Error)
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func DeleteDepartment(id uint32) mixin.ErrorCode {
	if err := db.Where("id = ?", id).Delete(Department{}).Error; err != nil {
		logrus.Errorf("[DeleteDepartment] error %s", err.Error)
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func GetAllDepartment() ([]Department, mixin.ErrorCode) {
	var departments []Department

	if err := db.Find(&departments).Error; err != nil {
		logrus.Errorf("[GetAllDepartment] error %s", err.Error)
		return nil, mixin.ErrorServerDb
	}

	return departments, mixin.StatusOK
}

// SELECT name, GROUP_CONCAT(purchase_url), GROUP_CONCAT(img_url) FROM departments GROUP BY name, size
func GetDepMapGroupByName() (map[SkuMapKey]Department, mixin.ErrorCode) {
	resp := make(map[SkuMapKey]Department)

	rows, err := db.Table("departments").Select("name, size, GROUP_CONCAT(purchase_url), GROUP_CONCAT(img_url)").Group("name, size").Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, mixin.ErrorServerDb
	}

	for rows.Next() {
		var dep Department
		if err := rows.Scan(&dep.Name, &dep.Size, &dep.PurchaseUrl, &dep.ImgUrl); err != nil {
			return nil, mixin.ErrorServerDb
		}
		resp[SkuMapKey{dep.Name, dep.Size}] = dep
	}

	return resp, mixin.StatusOK
}

// TODO 以下方法可以去掉
func GetAllDepartmentMap() (map[uint32]string, mixin.ErrorCode) {
	departmentMap := make(map[uint32]string)

	var departments []Department
	if err := db.Find(&departments).Error; err != nil {
		logrus.Errorf("[GetAllDepartmentMap] error %s", err.Error)
		return nil, mixin.ErrorServerDb
	}

	for _, data := range departments {
		departmentMap[data.ID] = data.Name
	}
	return departmentMap, mixin.StatusOK
}

func GetDepartmentByName(name string) (Department, mixin.ErrorCode) {
	var department Department
	if err := db.Where("name = ?", name).Find(&department).Error; err != nil {
		logrus.Errorf("[GetDepartmentByName] error %s", err.Error)
		return department, mixin.ErrorServerDb
	}
	return department, mixin.StatusOK
}
