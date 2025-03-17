package controller

import (
	"stock-management/internal/models"
	"stock-management/internal/services"
	rs "stock-management/pkgs/response"

	"github.com/gin-gonic/gin"
)

type SupplierController struct {
	supplierService services.SupplierService
}

func NewSupplierController(sv services.SupplierService) *SupplierController {
	return &SupplierController{
		supplierService: sv,
	}
}

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
	suppliers, err := pc.supplierService.SearchSupplierList(c, req)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	rs.SuccessResponse(c, suppliers)
}

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

	supplier, err := pc.supplierService.CreateSupplier(c, req)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	rs.SuccessResponse(c, supplier)
}
