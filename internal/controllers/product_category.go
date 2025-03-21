package controller

import (
	"stock-management/internal/models"
	"stock-management/internal/services"
	rs "stock-management/pkgs/response"

	"github.com/gin-gonic/gin"
)

var ProductCategory = new(ProductCategoryController)

type ProductCategoryController struct{}

// SearchProductCategoryList retrieves a list of product categorys based on search criteria
// @Summary Retrieve product category list
// @Description Returns a list of product categorys matching the search request
// @Tags ProductCategory
// @Accept  json
// @Produce  json
// @Param request body models.ProductCategorySearchReq true "ProductCategory search details"
// @Success 200 {object} models.SearchRp
// @Router /api/product-category/list [post]
func (pc *ProductCategoryController) GetProductCategoryList(c *gin.Context) {
	var req models.ProductCategorySearchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	code := req.Validate()
	if code != rs.OkCode {
		rs.FailResponseWithCode(c, code)
		return
	}
	productCategories, err := services.Service.CategoryService.GetProductCategoryList(c, req)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	rs.SuccessResponse(c, productCategories)
}

// CreateProductCategory creates a new product category
// @Summary Create product category
// @Description Creates a new product category and returns the created product category details
// @Tags ProductCategory
// @Accept  json
// @Produce  json
// @Param request body models.ProductCategoryCreateReq true "ProductCategory creation details"
// @Success 201 {object} models.ProductCategory
// @Router /api/product-category/create [post]
func (pc *ProductCategoryController) CreateProductCategory(c *gin.Context) {
	var req models.ProductCategoryCreateReq

	if err := c.ShouldBindJSON(&req); err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	code := req.Validate()
	if code != rs.OkCode {
		rs.FailResponseWithCode(c, code)
		return
	}

	productCategory, err := services.Service.CategoryService.CreateProductCategory(c, req)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	rs.SuccessResponse(c, productCategory)
}
