package models

import (
	"errors"
	"stock-management/pkgs/response"
	"stock-management/pkgs/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ProductID         uuid.UUID `gorm:"primaryKey;type:uuid;column:product_id" json:"product_id"`
	ProductName       string    `gorm:"not null;column:product_name" json:"product_name"`
	ProductReference  string    `gorm:"not null;column:product_reference" json:"product_reference"`
	Status            string    `gorm:"not null;column:status" json:"status"`
	ProductCategoryID uuid.UUID `gorm:"not null;column:product_category_id" json:"product_category_id"`
	Price             int       `gorm:"not null;column:price" json:"price"`
	StockLocation     string    `gorm:"not null;column:stock_location" json:"stock_location"`
	SupplierID        uuid.UUID `gorm:"not null;column:supplier_id" json:"supplier_id"`
	Quantity          int       `gorm:"not null;column:quantity" json:"quantity"`
	DateCreated       time.Time `gorm:"not null;column:date_created" json:"date_created"`

	//
	Supplier        Supplier        `json:"supplier"`
	ProductCategory ProductCategory `json:"product_category"`
}

func (p *Product) TableName() string {
	return "product"
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	p.DateCreated = time.Now().UTC()
	return nil
}

type ProductStatus string

const (
	ProductStatusAvailable  ProductStatus = "Available"
	ProductStatusOnOrder    ProductStatus = "On Order"
	ProductStatusOutOfStock ProductStatus = "Out of Stock"
)

var validProductStatus = []ProductStatus{ProductStatusAvailable, ProductStatusOutOfStock, ProductStatusOnOrder}

type ProductCreateReq struct {
	ProductName       string `json:"product_name"`
	ProductReference  string `json:"product_reference"`
	Status            string `json:"status"`
	ProductCategoryID string `json:"product_category_id"`
	Price             int    `json:"price"`
	StockLocation     string `json:"stock_location"`
	Quantity          int    `json:"quantity"`
	SupplierID        string `json:"supplier_id"`
}

func (req *ProductCreateReq) Validate() response.RespCode {
	if req.ProductName == "" {
		return response.ErrInvalidName
	}
	if req.ProductReference == "" {
		req.ProductReference = uuid.New().String()
	}
	req.ProductReference = "PROD-" + time.Now().Format("200601") + "-" + req.ProductReference

	if req.Status == "" || (ProductStatus(req.Status) != ProductStatusAvailable && ProductStatus(req.Status) != ProductStatusOnOrder && ProductStatus(req.Status) != ProductStatusOutOfStock) {
		return response.ErrInvalidStatus
	}
	if req.ProductCategoryID == "" || !utils.IsValidUUID(req.ProductCategoryID) {
		return response.ErrInvalidCategory
	}
	if req.StockLocation == "" {
		return response.ErrInvalidStockLocation
	}
	if req.SupplierID == "" || !utils.IsValidUUID(req.SupplierID) {
		return response.ErrInvalidSupplier
	}
	return response.OkCode
}

type ProductUpdateReq struct {
	ProductID         string `json:"product_id"`
	ProductName       string `json:"product_name"`
	ProductReference  string `json:"product_reference"`
	Status            string `json:"status"`
	ProductCategoryID string `json:"product_category_id"`
	Price             int    `json:"price"`
	StockLocation     string `json:"stock_location"`
	Quantity          int    `json:"quantity"`
	SupplierID        string `json:"supplier_id"`
}

func (req *ProductUpdateReq) Validate() response.RespCode {
	if req.ProductID == "" || !utils.IsValidUUID(req.ProductID) {
		return response.ErrInvalidProduct
	}
	if req.ProductName == "" {
		return response.ErrInvalidName
	}
	if req.ProductReference == "" {
		return response.ErrInvalidReference
	}
	if req.Status == "" || (ProductStatus(req.Status) != ProductStatusAvailable && ProductStatus(req.Status) != ProductStatusOnOrder && ProductStatus(req.Status) != ProductStatusOutOfStock) {
		return response.ErrInvalidStatus
	}
	if req.ProductCategoryID == "" || !utils.IsValidUUID(req.ProductCategoryID) {
		return response.ErrInvalidCategory
	}
	if req.StockLocation == "" {
		return response.ErrInvalidStockLocation
	}
	if req.SupplierID == "" || !utils.IsValidUUID(req.SupplierID) {
		return response.ErrInvalidSupplier
	}
	return response.OkCode
}

type ProductSearchReq struct {
	ProductReferences  []string `json:"product_references"`
	ProductNames       []string `json:"product_names,omitempty"`
	Status             []string `json:"status,omitempty"`
	ProductCategoryIDs []string `json:"product_category_ids,omitempty"`
	SupplierIDs        []string `json:"supplier_ids,omitempty"`
	PriceFrom          int      `json:"price_from"`
	PriceTo            int      `json:"price_to"`
	StockLocations     []string `json:"stock_locations,omitempty"`

	DateCreatedFrom string `json:"date_created_from,omitempty"`
	DateCreatedTo   string `json:"date_created_to,omitempty"`
	Pagination

	//convert
	ProductCategoryUUIDs []uuid.UUID
	SupplierUUIDs        []uuid.UUID
}

type ProductSearchRp struct {
	Data []Product `json:"data"`
	Pagination
}

type ProductByIdReq struct {
	ProductID string `json:"product_id"`
}

type ProductDistanceByIdReq struct {
	ProductID string `json:"product_id"`
}

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type SearchRp struct {
	Data interface{} `json:"data"`
	Pagination
}

func (req *ProductSearchReq) Validate() response.RespCode {
	if req.DateCreatedFrom != "" {
		if err := validateDateFormat(req.DateCreatedFrom); err != nil {
			return response.ErrInvalidDate
		}
	}
	if req.DateCreatedTo != "" {
		if err := validateDateFormat(req.DateCreatedTo); err != nil {
			return response.ErrInvalidDate
		}
	}
	if req.DateCreatedFrom != "" && req.DateCreatedTo != "" {
		dateFrom, err := time.Parse(dateFormat, req.DateCreatedFrom)
		if err != nil {
			return response.ErrInvalidDate
		}

		dateTo, err := time.Parse(dateFormat, req.DateCreatedTo)
		if err != nil {
			return response.ErrInvalidDate
		}
		if dateFrom.After(dateTo) {
			return response.ErrInvalidDate
		}
	}

	if len(req.ProductCategoryIDs) > 0 {
		req.ProductCategoryUUIDs = getUUIDs(req.ProductCategoryIDs)
	}
	if len(req.SupplierIDs) > 0 {
		req.SupplierUUIDs = getUUIDs(req.SupplierIDs)
	}
	return response.OkCode
}

const (
	dateFormat string = "2006-01-02"
)

func validateDateFormat(date string) error {
	_, err := time.Parse(dateFormat, date)
	if err != nil {
		return errors.New("invalid date format")
	}
	return nil
}
func getUUIDs(l []string) []uuid.UUID {
	res := []uuid.UUID{}
	for _, s := range l {
		uid, err := uuid.Parse(s)
		if err == nil {
			res = append(res, uid)
		}
	}
	return res
}

const (
	CategoryProductsKey     string = "a_category_products:%v"
	SupplierProductsKey     string = "a_supplier_products:%v"
	CategoryProductsScanKey string = "a_category_products:*"
	SupplierProductsScanKey string = "a_supplier_products:*"
	TotalProductsKey        string = "a_product_total"
)
