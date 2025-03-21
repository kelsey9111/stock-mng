package controller

import (
	"stock-management/internal/services"
	rs "stock-management/pkgs/response"

	"github.com/gin-gonic/gin"
)

var Statistics = new(StatisticsController)

type StatisticsController struct{}

// @Summary Get percentage of products per category
// @Description Get percentage of products per category
// @Tags Statistics
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]decimal.Decimal
// @Router /api/statistics/products-per-category [get]
func (pc *StatisticsController) GetProductPerCategory(c *gin.Context) {
	results, err := services.Service.ProductService.GetProductPerCategory(c)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	rs.SuccessResponse(c, results)
}

// @Summary Get percentage of products per supplier
// @Description Get percentage of products per supplier
// @Tags Statistics
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]decimal.Decimal
// @Router /api/statistics/products-per-supplier [get]
func (pc *StatisticsController) GetProductPerSupplier(c *gin.Context) {
	results, err := services.Service.ProductService.GetProductPerSupplier(c)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	rs.SuccessResponse(c, results)
}
