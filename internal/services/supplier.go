package services

import (
	"context"
	"stock-management/internal/models"
	"stock-management/internal/repo"
)

type SupplierService interface {
	SearchSupplierList(ctx context.Context, req models.SupplierSearchReq) (*models.SearchRp, error)
	CreateSupplier(ctx context.Context, supplier models.SupplierCreateReq) (*models.Supplier, error)
}

type supplierService struct {
	supplierRepo repo.SupplierRepo
}

func newSupplierService(supplierRepo repo.SupplierRepo) SupplierService {
	return &supplierService{
		supplierRepo: supplierRepo,
	}
}

func (ps *supplierService) SearchSupplierList(ctx context.Context, req models.SupplierSearchReq) (*models.SearchRp, error) {
	rs, offset, err := ps.supplierRepo.GetSupplierList(ctx, req)
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

func (ps *supplierService) CreateSupplier(ctx context.Context, supplier models.SupplierCreateReq) (*models.Supplier, error) {
	return ps.supplierRepo.CreateSupplier(ctx, supplier)
}
