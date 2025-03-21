package repo

import (
	"context"
	"errors"
	"stock-management/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SupplierRepo interface {
	GetSupplier(ctx context.Context, id uuid.UUID) (*models.Supplier, error)
	GetSuppliersByIds(ctx context.Context, ids []uuid.UUID) ([]models.Supplier, error)
	CreateSupplier(ctx context.Context, req models.SupplierCreateReq) (*models.Supplier, error)
	GetSupplierList(ctx context.Context, req models.SupplierSearchReq) ([]models.Supplier, int, error)
}

type supplierRepo struct {
	pdb *gorm.DB
}

func NewSupplierRepo(db *gorm.DB) SupplierRepo {
	return &supplierRepo{
		pdb: db,
	}
}

func (sr *supplierRepo) GetSupplierList(ctx context.Context, req models.SupplierSearchReq) ([]models.Supplier, int, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var suppliers []models.Supplier
	q := sr.pdb.WithContext(ctx).Model(&models.Supplier{})
	if req.SupplierName != "" {
		q = q.Where("supplier_name LIKE ?", "%"+req.SupplierName+"%")
	}
	if req.Status != "" {
		q = q.Where("status = ?", req.Status)
	}
	if err := q.Offset(req.Offset).Limit(req.Limit).Find(&suppliers).Error; err != nil {
		return nil, 0, err
	}
	var totalCount int64
	err := sr.pdb.WithContext(ctx).Model(&models.Supplier{}).Where(q.Statement.SQL.String(), q.Statement.Vars...).Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}
	return suppliers, int(totalCount), nil
}

func (sr *supplierRepo) GetSupplier(ctx context.Context, id uuid.UUID) (*models.Supplier, error) {
	var s models.Supplier
	if err := sr.pdb.WithContext(ctx).First(&s, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (sr *supplierRepo) GetSuppliersByIds(ctx context.Context, ids []uuid.UUID) ([]models.Supplier, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if len(ids) == 0 {
		return []models.Supplier{}, nil
	}
	var suppliers []models.Supplier
	if err := sr.pdb.WithContext(ctx).Where("supplier_id IN ?", ids).Find(&suppliers).Error; err != nil {
		return nil, err
	}
	return suppliers, nil
}

func (sr *supplierRepo) CreateSupplier(ctx context.Context, req models.SupplierCreateReq) (*models.Supplier, error) {
	supplier := models.Supplier{
		SupplierID:   uuid.New(),
		SupplierName: req.SupplierName,
		Status:       models.SupplierStatus(req.Status),
	}
	if err := sr.pdb.WithContext(ctx).Create(&supplier).Error; err != nil {
		return nil, err
	}
	return &supplier, nil
}
