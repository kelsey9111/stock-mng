package initialize

import (
	"fmt"
	"stock-management/global"
	"stock-management/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func checkErrorPannic(err error, errString string) {
	if err != nil {
		global.Logger.Fatal(fmt.Sprintf("%s: %v", errString, err))
	}
}

func InitPostgreSQL() {
	p := global.Config.PostgreSql

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
		p.Host, p.Username, p.Password, p.DBname, p.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	checkErrorPannic(err, "Open PostgreSQL error")

	global.Logger.Info("Init PostgreSQL success")
	global.Pdb = db

	err = installUUIDExtension()
	checkErrorPannic(err, "Install uuid-ossp extension error")

	setPostgrePool()
	migratePostgreTables()
	// SeedData()
}

func installUUIDExtension() error {
	sqlDb, err := global.Pdb.DB()
	if err != nil {
		return fmt.Errorf("error getting SQL DB: %v", err)
	}

	_, err = sqlDb.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err != nil {
		return fmt.Errorf("failed to create uuid-ossp extension: %v", err)
	}

	global.Logger.Info("uuid-ossp extension installed successfully")
	return nil
}

func setPostgrePool() {
	p := global.Config.PostgreSql

	sqlDb, err := global.Pdb.DB()
	if err != nil {
		fmt.Printf("PostgreSQL error: %s", err)
		return
	}

	sqlDb.SetConnMaxLifetime(time.Duration(p.MaxLifeTimeConns) * time.Second)
	sqlDb.SetMaxIdleConns(p.MaxIdleConns)
	sqlDb.SetMaxOpenConns(p.MaxOpenConns)
}

func migratePostgreTables() {
	err := global.Pdb.AutoMigrate(
		&models.Product{},
		&models.Supplier{},
		&models.ProductCategory{},
	)
	if err != nil {
		fmt.Printf("PostgreSQL migration error: %s", err)
		return
	}
	fmt.Print("PostgreSQL migration success")
}

func SeedData() {
	uid, _ := uuid.Parse("e2bb9114-6532-40c1-a9ec-ebc92274dc72")
	createdAt, _ := time.Parse("2006-01-02", "2025-03-07")
	updatedAt, _ := time.Parse("2006-01-02", "2025-03-07")

	categories := []models.ProductCategory{
		{ProductCategoryID: uid, ProductCategoryName: "Clothing", CreatedAt: createdAt, UpdatedAt: updatedAt, Status: models.CategoryActive},
	}

	if err := global.Pdb.Create(&categories).Error; err != nil {
		global.Logger.Fatal(fmt.Sprintf("Failed to seed categories: %v", err))
	} else {
		global.Logger.Info("Categories seeded successfully")
	}

	uid, _ = uuid.Parse("20dcd6f5-3ea5-41a6-a89d-59494b1e3170")
	suppliers := []models.Supplier{
		{SupplierID: uid, SupplierName: "Supplier A", Status: models.SupplierActive},
	}

	if err := global.Pdb.Create(&suppliers).Error; err != nil {
		global.Logger.Fatal(fmt.Sprintf("Failed to seed supplier: %v", err))
	} else {
		global.Logger.Info("Supplier seeded successfully")
	}

	global.Logger.Info("Seed data added successfully")
}
