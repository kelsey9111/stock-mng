package initialize

import (
	"fmt"
	"stock-management/global"
	"stock-management/internal/models"
	"time"

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
	// migratePostgreTables()
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
