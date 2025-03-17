package manager

import (
	"stock-management/global"
	controller "stock-management/internal/controllers"
	"stock-management/internal/repo"
	"stock-management/internal/services"
	rs "stock-management/pkgs/response"

	"github.com/gin-gonic/gin"
)

type AdminRouter struct{}

func (r *AdminRouter) InitAdminRouter(router *gin.RouterGroup) {
	productCategoryRp := repo.NewCategoryRepo()
	supplierRp := repo.NewSupplierRepo()
	pr := repo.NewProductRepo(productCategoryRp, supplierRp)

	ps := services.NewProductService(pr)
	pc := controller.NewProductController(ps)

	ss := services.NewSupplierService(supplierRp)
	sc := controller.NewSupplierController(ss)

	cs := services.NewProductCategoryService(productCategoryRp)
	cc := controller.NewProductCategoryController(cs)

	router.GET("ping", Pong)

	productRouter := router.Group("product")
	{
		productRouter.POST("search", pc.GetProductList)
		productRouter.POST("detail", pc.GetProduct)
		productRouter.POST("create", pc.CreateProduct)
		productRouter.PUT("update", pc.UpdateProduct)
		productRouter.POST("export", pc.ExportProductsToPDF)
		productRouter.GET("distance", pc.GetProductDistance)
	}

	productCategoryRouter := router.Group("product-category")
	{
		productCategoryRouter.POST("search", cc.GetProductCategoryList)
		productCategoryRouter.POST("create", cc.CreateProductCategory)
	}

	supplierRouter := router.Group("supplier")
	{
		supplierRouter.POST("search", sc.SearchSupplierList)
		supplierRouter.POST("create", sc.CreateSupplier)
	}

	statisticsRouter := router.Group("statistics")
	{
		statisticsRouter.GET("products-per-category", pc.GetProductPerCategory)
		statisticsRouter.GET("products-per-supplier", pc.GetProductPerSupplier)
	}
}

func Pong(c *gin.Context) {
	global.Logger.Info("access pong function")
	uid := c.Query("uid")
	rs.SuccessResponse(c, uid)
}
