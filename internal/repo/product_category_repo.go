package repo

import (
	"context"
	"errors"
	"stock-management/internal/models"
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type ProductCategoryRepo interface {
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.ProductCategory, error)
	GetCategoriesByIds(ctx context.Context, ids []uuid.UUID) ([]models.ProductCategory, error)
	CreateProductCategory(ctx context.Context, req models.ProductCategoryCreateReq) (*models.ProductCategory, error)
	GetProductCategoryList(ctx context.Context, req models.ProductCategorySearchReq) ([]models.ProductCategory, int, error)
}
type productCategoryRepo struct {
	pdb *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) ProductCategoryRepo {
	return &productCategoryRepo{
		pdb: db,
	}
}

func (cr *productCategoryRepo) GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.ProductCategory, error) {
	var s models.ProductCategory
	if err := cr.pdb.WithContext(ctx).First(&s, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (cr *productCategoryRepo) GetCategoriesByIds(ctx context.Context, ids []uuid.UUID) ([]models.ProductCategory, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if len(ids) == 0 {
		return []models.ProductCategory{}, nil
	}
	var categories []models.ProductCategory
	if err := cr.pdb.WithContext(ctx).Where("product_category_id IN ?", ids).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (cr *productCategoryRepo) GetProductCategoryList(ctx context.Context, req models.ProductCategorySearchReq) ([]models.ProductCategory, int, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var productCategorys []models.ProductCategory
	q := cr.pdb.WithContext(ctx).Model(&models.ProductCategory{})
	if req.ProductCategoryName != "" {
		q = q.Where("product_category_name LIKE ?", "%"+req.ProductCategoryName+"%")
	}
	if req.Status != "" {
		q = q.Where("status = ?", req.Status)
	}
	if err := q.Offset(req.Offset).Limit(req.Limit).Find(&productCategorys).Error; err != nil {
		return nil, 0, err
	}
	var totalCount int64
	err := cr.pdb.WithContext(ctx).Model(&models.ProductCategory{}).Where(q.Statement.SQL.String(), q.Statement.Vars...).Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}
	return productCategorys, int(totalCount), nil
}

func (cr *productCategoryRepo) CreateProductCategory(ctx context.Context, req models.ProductCategoryCreateReq) (*models.ProductCategory, error) {
	productCategory := models.ProductCategory{
		ProductCategoryID:   uuid.New(),
		ProductCategoryName: req.ProductCategoryName,
		Status:              models.ProductCategoryStatus(req.Status),
		CreatedAt:           time.Now().Format("2006-01-02"),
		UpdatedAt:           time.Now().Format("2006-01-02"),
	}
	if err := cr.pdb.WithContext(ctx).Create(&productCategory).Error; err != nil {
		return nil, err
	}
	return &productCategory, nil
}
