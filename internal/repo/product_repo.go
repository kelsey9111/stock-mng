package repo

import (
	"context"
	"errors"
	"fmt"
	"stock-management/internal/models"
	"stock-management/pkgs/utils"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductRepo interface {
	GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error)
	GetProductList(ctx context.Context, req models.ProductSearchReq) ([]models.Product, int, error)
	CreateProduct(ctx context.Context, product models.ProductCreateReq) (models.Product, error)
	UpdateProduct(ctx context.Context, product models.ProductUpdateReq) (models.Product, error)
	GetProductPercentagePerKey(ctx context.Context, key string) (map[string]decimal.Decimal, error)
}

type productRepo struct {
	pdb                 *gorm.DB
	cache               *redis.Client
	productCategoryRepo ProductCategoryRepo
	supplierRepo        SupplierRepo
}

func NewProductRepo(db *gorm.DB, redis *redis.Client, cr ProductCategoryRepo, sp SupplierRepo) ProductRepo {
	return &productRepo{
		pdb:                 db,
		cache:               redis,
		productCategoryRepo: cr,
		supplierRepo:        sp,
	}
}

func (pr *productRepo) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var product models.Product
	err := pr.pdb.WithContext(ctx).Preload("Supplier").Preload("ProductCategory").First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (pr *productRepo) GetProductList(ctx context.Context, req models.ProductSearchReq) ([]models.Product, int, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var categoryUUIDs []uuid.UUID
	var supplierUUIDs []uuid.UUID

	wg := utils.NewWgGroup()
	wg.Go(func() error {
		if len(req.ProductCategoryUUIDs) > 0 {
			exists, err := pr.productCategoryRepo.GetCategoriesByIds(ctx, req.ProductCategoryUUIDs)
			if err != nil {
				return err
			}
			for _, item := range exists {
				categoryUUIDs = append(categoryUUIDs, item.ProductCategoryID)
			}
		}
		return nil
	})

	wg.Go(func() error {
		if len(req.SupplierUUIDs) > 0 {
			exists, err := pr.supplierRepo.GetSuppliersByIds(ctx, req.SupplierUUIDs)
			if err != nil {
				return err
			}
			for _, item := range exists {
				supplierUUIDs = append(supplierUUIDs, item.SupplierID)
			}
		}
		return nil
	})

	err := wg.Wait()
	if err != nil {
		return nil, 0, err
	}

	var products []models.Product

	if len(req.ProductCategoryUUIDs) > 0 && len(categoryUUIDs) == 0 {
		return products, 0, nil
	}
	if len(req.SupplierUUIDs) > 0 && len(supplierUUIDs) == 0 {
		return products, 0, nil
	}

	q := pr.pdb.WithContext(ctx).Model(&models.Product{})

	q = pr.applyFilters(q, req, categoryUUIDs, supplierUUIDs)

	// If Offset and Limit are both 0, fetch all products without pagination
	if req.Offset == 0 && req.Limit == 0 {
		err = q.Preload("Supplier").Preload("ProductCategory").Find(&products).Error
	} else {
		// Apply pagination if Offset and Limit are specified
		err = q.Order("date_created desc").
			Limit(req.Limit).
			Offset(req.Offset).
			Preload("Supplier").
			Preload("ProductCategory").
			Find(&products).Error
	}

	if err != nil {
		return nil, 0, err
	}

	totalCount, err := pr.getTotalCount(ctx, req, categoryUUIDs, supplierUUIDs)
	if err != nil {
		return nil, 0, err
	}

	nextOffset := req.Offset + req.Limit
	if nextOffset >= int(totalCount) {
		nextOffset = 0
	}

	return products, nextOffset, nil
}

func (pr *productRepo) applyFilters(q *gorm.DB, req models.ProductSearchReq, categoryUUIDs, supplierUUIDs []uuid.UUID) *gorm.DB {
	if len(req.ProductNames) > 0 {
		q = q.Where("product_name IN (?)", req.ProductNames)
	}
	if len(req.ProductReferences) > 0 {
		q = q.Where("product_reference IN (?)", req.ProductReferences)
	}
	if len(req.Status) > 0 {
		q = q.Where("status IN (?)", req.Status)
	}
	if req.PriceFrom > 0 {
		q = q.Where("price >= ?", req.PriceFrom)
	}
	if req.PriceTo > 0 {
		q = q.Where("price <= ?", req.PriceTo)
	}
	if len(req.StockLocations) > 0 {
		q = q.Where("stock_location IN (?)", req.StockLocations)
	}
	if req.DateCreatedFrom != "" {
		q = q.Where("date_created >= ?", req.DateCreatedFrom)
	}
	if req.DateCreatedTo != "" {
		q = q.Where("date_created <= ?", req.DateCreatedTo)
	}

	if len(categoryUUIDs) > 0 {
		q = q.Where("product_category_id IN (?)", categoryUUIDs)
	}
	if len(supplierUUIDs) > 0 {
		q = q.Where("supplier_id IN (?)", supplierUUIDs)
	}

	return q
}

func (pr *productRepo) getTotalCount(ctx context.Context, req models.ProductSearchReq, categoryUUIDs, supplierUUIDs []uuid.UUID) (int64, error) {
	countQuery := pr.pdb.WithContext(ctx).Model(&models.Product{})

	// Apply filters to count query
	countQuery = pr.applyFilters(countQuery, req, categoryUUIDs, supplierUUIDs)

	var totalCount int64
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return 0, err
	}

	return totalCount, nil
}

func (pr *productRepo) CreateProduct(ctx context.Context, req models.ProductCreateReq) (models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx := pr.pdb.Begin()
	if tx.Error != nil {
		return models.Product{}, tx.Error
	}

	product := models.Product{
		ProductID:         uuid.New(),
		ProductName:       req.ProductName,
		ProductReference:  req.ProductReference,
		Status:            req.Status,
		ProductCategoryID: uuid.MustParse(req.ProductCategoryID),
		Price:             req.Price,
		StockLocation:     req.StockLocation,
		Quantity:          req.Quantity,
		SupplierID:        uuid.MustParse(req.SupplierID),
	}

	var err error

	wg := utils.NewWgGroup()
	wg.Go(func() error {
		cacheKey := fmt.Sprintf(models.SupplierProductsKey, product.SupplierID)
		_, err := pr.cache.Get(ctx, cacheKey).Result()
		if err == nil {
			return nil
		} else if err == redis.Nil {
			var supplier *models.Supplier
			supplier, err = pr.supplierRepo.GetSupplier(ctx, product.SupplierID)
			if err != nil {
				return err
			}
			if supplier == nil {
				return errors.New("invalid supplier")
			}
		}
		return nil
	})

	wg.Go(func() error {
		cacheKey := fmt.Sprintf(models.CategoryProductsKey, product.ProductCategoryID)
		_, err := pr.cache.Get(ctx, cacheKey).Result()
		if err == nil {
			return nil
		} else if err == redis.Nil {
			var productCategory *models.ProductCategory
			productCategory, err = pr.productCategoryRepo.GetCategoryByID(ctx, product.ProductCategoryID)
			if err != nil {
				return err
			}
			if productCategory == nil {
				return errors.New("invalid product category")
			}
		}
		return nil
	})

	err = wg.Wait()
	if err != nil {
		tx.Rollback()
		return models.Product{}, err
	}

	if err := tx.WithContext(ctx).Create(&product).Error; err != nil {
		tx.Rollback()
		return models.Product{}, err
	}

	//
	pipe := pr.cache.Pipeline()
	pipe.Incr(ctx, fmt.Sprintf(models.SupplierProductsKey, product.SupplierID))
	pipe.Incr(ctx, models.TotalProductsKey)
	pipe.Incr(ctx, fmt.Sprintf(models.CategoryProductsKey, product.ProductCategoryID))
	_, err = pipe.Exec(ctx)
	if err != nil {
		tx.Rollback()
		return models.Product{}, err
	}
	//

	if err := tx.Commit().Error; err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func (pr *productRepo) UpdateProduct(ctx context.Context, req models.ProductUpdateReq) (models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx := pr.pdb.Begin()
	if tx.Error != nil {
		return models.Product{}, tx.Error
	}

	var product models.Product
	if err := tx.WithContext(ctx).First(&product, req.ProductID).Error; err != nil {
		tx.Rollback()
		return models.Product{}, err
	}

	preProductCategoryID := product.ProductCategoryID
	newProductCategoryID := uuid.MustParse(req.ProductCategoryID)
	preSupplierID := product.SupplierID
	newSupplierID := uuid.MustParse(req.SupplierID)

	product.ProductName = req.ProductName
	product.ProductReference = req.ProductReference
	product.Status = req.Status
	product.Price = req.Price
	product.Quantity = req.Quantity
	product.SupplierID = uuid.MustParse(req.SupplierID)
	product.StockLocation = req.StockLocation
	product.ProductCategoryID = uuid.MustParse(req.ProductCategoryID)

	wg := utils.NewWgGroup()
	var err error

	if preSupplierID != newSupplierID {
		wg.Go(func() error {
			cacheKey := fmt.Sprintf(models.SupplierProductsKey, product.SupplierID)
			_, err := pr.cache.Get(ctx, cacheKey).Result()
			if err == nil {
				return nil
			} else if err == redis.Nil {
				var supplier *models.Supplier
				supplier, err = pr.supplierRepo.GetSupplier(ctx, product.SupplierID)
				if err != nil {
					return err
				}
				if supplier == nil {
					return errors.New("invalid supplier")
				}
			}
			return nil
		})
	}

	if preProductCategoryID != newProductCategoryID {
		wg.Go(func() error {
			cacheKey := fmt.Sprintf(models.CategoryProductsKey, product.ProductCategoryID)
			_, err := pr.cache.Get(ctx, cacheKey).Result()
			if err == nil {
				return nil
			} else if err == redis.Nil {
				var productCategory *models.ProductCategory
				productCategory, err = pr.productCategoryRepo.GetCategoryByID(ctx, product.ProductCategoryID)
				if err != nil {
					return err
				}
				if productCategory == nil {
					return errors.New("invalid product category")
				}
			}
			return nil
		})
	}

	err = wg.Wait()
	if err != nil {
		tx.Rollback()
		return models.Product{}, err
	}

	if err := tx.WithContext(ctx).Save(&product).Error; err != nil {
		tx.Rollback()
		return models.Product{}, err
	}

	// Update cache
	pipe := pr.cache.Pipeline()
	if preSupplierID != newSupplierID {
		pipe.Incr(ctx, fmt.Sprintf(models.SupplierProductsKey, newSupplierID))
		pipe.Decr(ctx, fmt.Sprintf(models.SupplierProductsKey, preSupplierID))
	}
	if preProductCategoryID != newProductCategoryID {
		pipe.Incr(ctx, fmt.Sprintf(models.CategoryProductsKey, newProductCategoryID))
		pipe.Decr(ctx, fmt.Sprintf(models.CategoryProductsKey, preProductCategoryID))
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		tx.Rollback()
		return models.Product{}, err
	}
	//

	if err := tx.Commit().Error; err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func (pr *productRepo) GetProductPercentagePerKey(ctx context.Context, key string) (map[string]decimal.Decimal, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	countList := make(map[string]int64)
	percentages := make(map[string]decimal.Decimal)

	totalProducts, err := pr.getTotalProducts(ctx)
	if err != nil {
		return nil, err
	}

	if totalProducts == 0 {
		return percentages, nil
	}

	var cursor uint64
	for {
		keys, newCursor, err := pr.cache.Scan(ctx, cursor, key, 100).Result()
		if err != nil {
			continue
		}

		for _, key := range keys {
			parts := strings.Split(key, ":")
			if len(parts) < 2 {
				continue
			}
			id := parts[1]

			countp, err := pr.cache.Get(ctx, key).Result()
			if err != nil {
				continue
			}

			count, parseErr := strconv.ParseInt(countp, 10, 64)
			if parseErr != nil {
				continue
			}

			countList[id] = count
		}

		cursor = newCursor
		if cursor == 0 {
			break
		}
	}

	for id, count := range countList {
		percentages[id] = decimal.NewFromInt(count).Mul(decimal.NewFromInt(100)).Div(decimal.NewFromInt(totalProducts)).Round(2)
	}

	return percentages, nil
}

func (pr *productRepo) getTotalProducts(ctx context.Context) (int64, error) {
	val, err := pr.cache.Get(ctx, models.TotalProductsKey).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	total, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting value to int64: %v", err)
	}
	return total, nil
}
