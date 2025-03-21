package services

import (
	"context"
	"stock-management/internal/models"
	"stock-management/internal/repo"
)

type ProductCategoryService interface {
	GetProductCategoryList(ctx context.Context, req models.ProductCategorySearchReq) (*models.SearchRp, error)
	CreateProductCategory(ctx context.Context, productCategory models.ProductCategoryCreateReq) (*models.ProductCategory, error)
}

type productCategoryService struct {
	productCategoryRepo repo.ProductCategoryRepo
}

func newProductCategoryService(productCategoryRepo repo.ProductCategoryRepo) ProductCategoryService {
	return &productCategoryService{
		productCategoryRepo: productCategoryRepo,
	}
}

func (ps *productCategoryService) GetProductCategoryList(ctx context.Context, req models.ProductCategorySearchReq) (*models.SearchRp, error) {
	rs, offset, err := ps.productCategoryRepo.GetProductCategoryList(ctx, req)
	if err != nil {
		return nil, err
	}
	result := &models.SearchRp{
		Data: rs,
		Pagination: models.Pagination{
			Offset: offset,
			Limit:  req.Limit,
		},
	}
	return result, nil
}

func (ps *productCategoryService) CreateProductCategory(ctx context.Context, productCategory models.ProductCategoryCreateReq) (*models.ProductCategory, error) {
	return ps.productCategoryRepo.CreateProductCategory(ctx, productCategory)
}
