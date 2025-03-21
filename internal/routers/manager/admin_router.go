package manager

import (
	"stock-management/global"
	controller "stock-management/internal/controllers"
	rs "stock-management/pkgs/response"

	"github.com/gin-gonic/gin"
)

type AdminRouter struct{}

func (r *AdminRouter) InitAdminRouter(router *gin.RouterGroup) {

	router.GET("ping", Pong)
	productRouter := router.Group("product")
	{
		productRouter.POST("list", controller.Product.GetProductList)
		productRouter.POST("detail", controller.Product.GetProduct)
		productRouter.POST("create", controller.Product.CreateProduct)
		productRouter.PUT("update", controller.Product.UpdateProduct)
		productRouter.POST("export", controller.Product.ExportProductsToPDF)
		productRouter.POST("distance", controller.Product.GetProductDistance)
	}

	productCategoryRouter := router.Group("product-category")
	{
		productCategoryRouter.POST("list", controller.ProductCategory.GetProductCategoryList)
		productCategoryRouter.POST("create", controller.ProductCategory.CreateProductCategory)
	}

	supplierRouter := router.Group("supplier")
	{
		supplierRouter.POST("list", controller.Supplier.SearchSupplierList)
		supplierRouter.POST("create", controller.Supplier.CreateSupplier)
	}

	statisticsRouter := router.Group("statistics")
	{
		statisticsRouter.GET("products-per-category", controller.Statistics.GetProductPerCategory)
		statisticsRouter.GET("products-per-supplier", controller.Statistics.GetProductPerSupplier)
	}
}

func Pong(c *gin.Context) {
	global.Logger.Info("access pong function")
	uid := c.Query("uid")
	rs.SuccessResponse(c, uid)
}
