package services

import (
	"context"
	"stock-management/internal/models"
	"stock-management/internal/repo"

	"github.com/google/uuid"
)

type ProductService interface {
	GetProductList(ctx context.Context, req models.ProductSearchReq) (*models.ProductSearchRp, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error)
	CreateProduct(ctx context.Context, product models.ProductCreateReq) (models.Product, error)
	UpdateProduct(ctx context.Context, product models.ProductUpdateReq) (models.Product, error)
	GetProductPerCategory(ctx context.Context) (map[string]int64, error)
	GetProductPerSupplier(ctx context.Context) (map[string]int64, error)
	GetProductDistance(ctx context.Context, id uuid.UUID, ip string) (float64, error)
}

type productService struct {
	productRepo repo.ProductRepo
}

func NewProductService(productRepo repo.ProductRepo) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

func (ps *productService) GetProductList(ctx context.Context, req models.ProductSearchReq) (*models.ProductSearchRp, error) {
	rs, offset, err := ps.productRepo.GetProductList(ctx, req)
	if err != nil {
		return nil, err
	}
	result := &models.ProductSearchRp{
		Data: rs,
		Pagination: models.Pagination{
			Offset: offset,
			Limit:  req.Limit,
		},
	}
	return result, nil
}

func (ps *productService) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	return ps.productRepo.GetProduct(ctx, id)
}

func (ps *productService) CreateProduct(ctx context.Context, product models.ProductCreateReq) (models.Product, error) {
	return ps.productRepo.CreateProduct(ctx, product)
}

func (ps *productService) UpdateProduct(ctx context.Context, product models.ProductUpdateReq) (models.Product, error) {
	return ps.productRepo.UpdateProduct(ctx, product)
}

func (ps *productService) GetProductPerCategory(ctx context.Context) (map[string]int64, error) {
	return ps.productRepo.GetProductPercentagePerKey(ctx, models.CategoryProductsScanKey)
}

func (ps *productService) GetProductPerSupplier(ctx context.Context) (map[string]int64, error) {
	return ps.productRepo.GetProductPercentagePerKey(ctx, models.SupplierProductsScanKey)
}

func (ps *productService) GetProductDistance(ctx context.Context, id uuid.UUID, ip string) (float64, error) {
	// todo
	return 0, nil
}
