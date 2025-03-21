package models

import (
	"stock-management/pkgs/response"

	"github.com/google/uuid"
)

type ProductCategory struct {
	ProductCategoryID   uuid.UUID             `gorm:"primaryKey;type:uuid;column:product_category_id" json:"product_category_id"`
	ProductCategoryName string                `gorm:"not null;column:product_category_name" json:"product_category_name"`
	Status              ProductCategoryStatus `gorm:"not null;column:status" json:"status"`
	CreatedAt           string                `gorm:"not null;column:created_at" json:"created_at"`
	UpdatedAt           string                `gorm:"not null;column:updated_at" json:"updated_at"`
}

func (p *ProductCategory) TableName() string {
	return "product_category"
}

type ProductCategoryStatus string

const (
	CategoryActive   ProductCategoryStatus = "active"
	CategoryInActive ProductCategoryStatus = "in active"
)

type ProductCategoryCreateReq struct {
	ProductCategoryName string                `json:"product_category_name"`
	Status              ProductCategoryStatus `json:"status"`
}

func (req *ProductCategoryCreateReq) Validate() response.RespCode {
	if req.ProductCategoryName == "" {
		return response.ErrInvalidName
	}
	if req.Status == "" || (ProductCategoryStatus(req.Status) != CategoryActive && ProductCategoryStatus(req.Status) != CategoryInActive) {
		return response.ErrInvalidStatus
	}
	return response.OkCode
}

type ProductCategorySearchReq struct {
	ProductCategoryName string         `json:"supplier_name,omitempty"`
	Status              SupplierStatus `json:"status,omitempty"`
	Pagination
}

func (req *ProductCategorySearchReq) Validate() response.RespCode {
	if req.Limit == 0 {
		req.Limit = 20
	}
	return response.OkCode
}
