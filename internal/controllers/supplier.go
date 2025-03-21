package controller

import (
	"stock-management/internal/models"
	"stock-management/internal/services"
	rs "stock-management/pkgs/response"

	"github.com/gin-gonic/gin"
)

var Supplier = new(SupplierController)

type SupplierController struct{}

// SearchSupplierList retrieves a list of suppliers based on search criteria
// @Summary Retrieve supplier list
// @Description Returns a list of suppliers matching the search request
// @Tags Supplier
// @Accept  json
// @Produce  json
// @Param request body models.SupplierSearchReq true "Supplier search details"
// @Success 200 {object} models.SearchRp
// @Router /api/supplier/list [post]
func (pc *SupplierController) SearchSupplierList(c *gin.Context) {
	var req models.SupplierSearchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	code := req.Validate()
	if code != rs.OkCode {
		rs.FailResponseWithCode(c, code)
		return
	}
	suppliers, err := services.Service.SupplierService.SearchSupplierList(c, req)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	rs.SuccessResponse(c, suppliers)
}

// CreateSupplier creates a new supplier
// @Summary Create supplier
// @Description Creates a new supplier and returns the created supplier details
// @Tags Supplier
// @Accept  json
// @Produce  json
// @Param request body models.SupplierCreateReq true "Supplier creation details"
// @Success 201 {object} models.Supplier
// @Router /api/supplier/create [post]
func (pc *SupplierController) CreateSupplier(c *gin.Context) {
	var req models.SupplierCreateReq

	if err := c.ShouldBindJSON(&req); err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	code := req.Validate()
	if code != rs.OkCode {
		rs.FailResponseWithCode(c, code)
		return
	}

	supplier, err := services.Service.SupplierService.CreateSupplier(c, req)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	rs.SuccessResponse(c, supplier)
}
