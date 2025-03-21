package services

import (
	"context"
	"fmt"
	"stock-management/internal/models"
	"stock-management/internal/repo"
	"stock-management/pkgs/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProductService interface {
	GetProductList(ctx context.Context, req models.ProductSearchReq) (*models.ProductSearchRp, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error)
	CreateProduct(ctx context.Context, product models.ProductCreateReq) (models.Product, error)
	UpdateProduct(ctx context.Context, product models.ProductUpdateReq) (models.Product, error)
	GetProductPerCategory(ctx context.Context) (map[string]decimal.Decimal, error)
	GetProductPerSupplier(ctx context.Context) (map[string]decimal.Decimal, error)
	GetProductDistance(ctx context.Context, id uuid.UUID, ip string) (*models.ProductDistanceRp, error)
}

type productService struct {
	productRepo repo.ProductRepo
}

func newProductService(productRepo repo.ProductRepo) ProductService {
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

func (ps *productService) GetProductPerCategory(ctx context.Context) (map[string]decimal.Decimal, error) {
	return ps.productRepo.GetProductPercentagePerKey(ctx, models.CategoryProductsScanKey)
}

func (ps *productService) GetProductPerSupplier(ctx context.Context) (map[string]decimal.Decimal, error) {
	return ps.productRepo.GetProductPercentagePerKey(ctx, models.SupplierProductsScanKey)
}

func (ps *productService) GetProductDistance(ctx context.Context, id uuid.UUID, ip string) (*models.ProductDistanceRp, error) {
	p, err := ps.productRepo.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("product is invalid %s", id)
	}
	addressIp, err := utils.GetCoordinatesFromIP(ip)
	if err != nil {
		return nil, err
	}
	recheckGetAddressIp, err := utils.GetAddressFromLatLonOSM(addressIp)
	if err != nil {
		return nil, err
	}

	addressProduct, err := utils.GetCoordinatesFromCity(p.StockLocation)
	if err != nil {
		return nil, err
	}
	recheckGetAddressProduct, err := utils.GetAddressFromLatLonOSM(addressProduct)
	if err != nil {
		return nil, err
	}
	distance := utils.CalculateDistance(addressIp, addressProduct)

	rs := models.ProductDistanceRp{
		IpCurrentCity:     recheckGetAddressIp,
		StockLocationCity: recheckGetAddressProduct,
		Distance:          fmt.Sprintf("%.2f km", distance),
	}
	return &rs, nil
}
