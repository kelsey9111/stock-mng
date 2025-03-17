package controller

import (
	"stock-management/internal/models"
	"stock-management/internal/services"
	rs "stock-management/pkgs/response"

	"github.com/gin-gonic/gin"
)

type ProductCategoryController struct {
	productCategoryService services.ProductCategoryService
}

func NewProductCategoryController(sv services.ProductCategoryService) *ProductCategoryController {
	return &ProductCategoryController{
		productCategoryService: sv,
	}
}

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
	productCategories, err := pc.productCategoryService.GetProductCategoryList(c, req)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	rs.SuccessResponse(c, productCategories)
}

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

	productCategory, err := pc.productCategoryService.CreateProductCategory(c, req)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	rs.SuccessResponse(c, productCategory)
}
