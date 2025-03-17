package models

import (
	"stock-management/pkgs/response"

	"github.com/google/uuid"
)

type Supplier struct {
	SupplierID   uuid.UUID      `gorm:"primaryKey;type:uuid;column:supplier_id" json:"supplier_id"`
	SupplierName string         `gorm:"not null;column:supplier_name" json:"supplier_name"`
	Status       SupplierStatus `gorm:"not null;column:status" json:"status"`
}

func (Supplier) TableName() string {
	return "supplier"
}

type SupplierStatus string

const (
	SupplierActive   SupplierStatus = "active"
	SupplierInActive SupplierStatus = "in active"
)

type SupplierCreateReq struct {
	SupplierName string `json:"supplier_name"`
	Status       string `json:"status"`
}

func (req *SupplierCreateReq) Validate() response.RespCode {
	if req.SupplierName == "" {
		return response.ErrInvalidName
	}
	if req.Status == "" || (SupplierStatus(req.Status) != SupplierActive && SupplierStatus(req.Status) != SupplierInActive) {
		return response.ErrInvalidStatus
	}
	return response.OkCode
}

type SupplierSearchReq struct {
	SupplierName string         `json:"supplier_name,omitempty"`
	Status       SupplierStatus `json:"status,omitempty"`
	Pagination
}

func (req *SupplierSearchReq) Validate() response.RespCode {
	if req.Limit == 0 {
		req.Limit = 20
	}
	return response.OkCode
}
