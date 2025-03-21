package services

import (
	"stock-management/global"
	"stock-management/internal/repo"
)

var Service *service

type service struct {
	CategoryService ProductCategoryService
	ProductService  ProductService
	SupplierService SupplierService
}

func InitService() {
	productCategoryRepo := repo.NewCategoryRepo(global.Pdb)
	supplierRepo := repo.NewSupplierRepo(global.Pdb)
	productRepo := repo.NewProductRepo(global.Pdb, global.Rdb, productCategoryRepo, supplierRepo)
	productServices := newProductService(productRepo)
	categoryServices := newProductCategoryService(productCategoryRepo)
	supplierService := newSupplierService(supplierRepo)

	Service = &service{
		CategoryService: categoryServices,
		ProductService:  productServices,
		SupplierService: supplierService,
	}
}
