package controller

import (
	"fmt"
	"stock-management/internal/models"
	"stock-management/internal/services"
	rs "stock-management/pkgs/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

type ProductController struct {
	productService services.ProductService
}

func NewProductController(sv services.ProductService) *ProductController {
	return &ProductController{
		productService: sv,
	}
}

func (pc *ProductController) GetProductList(c *gin.Context) {
	var req models.ProductSearchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	code := req.Validate()
	if code != rs.OkCode {
		rs.FailResponseWithCode(c, code)
		return
	}

	// if req.Limit == 0 {
	// 	req.Limit = 20
	// }

	products, err := pc.productService.GetProductList(c, req)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	rs.SuccessResponse(c, products)
}

func (pc *ProductController) CreateProduct(c *gin.Context) {
	var req models.ProductCreateReq

	if err := c.ShouldBindJSON(&req); err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	code := req.Validate()
	if code != rs.OkCode {
		rs.FailResponseWithCode(c, code)
		return
	}

	product, err := pc.productService.CreateProduct(c, req)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	rs.SuccessResponse(c, product)
}

func (pc *ProductController) UpdateProduct(c *gin.Context) {
	var req models.ProductUpdateReq

	if err := c.ShouldBindJSON(&req); err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	code := req.Validate()
	if code != rs.OkCode {
		rs.FailResponseWithCode(c, code)
		return
	}

	product, err := pc.productService.UpdateProduct(c, req)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	rs.SuccessResponse(c, product)
}

func (pc *ProductController) GetProduct(c *gin.Context) {
	var req models.ProductByIdReq

	if err := c.ShouldBindJSON(&req); err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		rs.FailResponseWithCode(c, rs.ErrInvalidProduct)
	}
	product, err := pc.productService.GetProduct(c, productID)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}

	rs.SuccessResponse(c, product)
}

func (pc *ProductController) GetProductPerCategory(c *gin.Context) {
	results, err := pc.productService.GetProductPerCategory(c)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
	}
	rs.SuccessResponse(c, results)
}

func (pc *ProductController) GetProductPerSupplier(c *gin.Context) {
	results, err := pc.productService.GetProductPerSupplier(c)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
	}
	rs.SuccessResponse(c, results)
}

func (pc *ProductController) ExportProductsToPDF(c *gin.Context) {
	var req models.ProductSearchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	code := req.Validate()
	if code != rs.OkCode {
		rs.FailResponseWithCode(c, code)
		return
	}
	req.Offset = 0
	req.Limit = 0

	products, err := pc.productService.GetProductList(c, req)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Product List")
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(30, 10, "Product Name")
	pdf.Cell(20, 10, "Price")
	pdf.Cell(20, 10, "Quantity")
	pdf.Cell(30, 10, "Stock Location")
	pdf.Cell(30, 10, "Supplier")
	pdf.Cell(30, 10, "Category")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	for _, product := range products.Data {
		pdf.Cell(30, 10, product.ProductName)
		pdf.Cell(20, 10, fmt.Sprintf("%d", product.Price))
		pdf.Cell(20, 10, fmt.Sprintf("%d", product.Quantity))
		pdf.Cell(30, 10, product.StockLocation)
		pdf.Cell(30, 10, product.Supplier.SupplierName)
		pdf.Cell(30, 10, product.ProductCategory.ProductCategoryName)
		pdf.Ln(10)
	}
	err = pdf.OutputFileAndClose("products.pdf")
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	c.File("products.pdf")
}

func (pc *ProductController) GetProductDistance(c *gin.Context) {
	var req models.ProductDistanceByIdReq

	if err := c.ShouldBindJSON(&req); err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		rs.FailResponseWithCode(c, rs.ErrInvalidProduct)
	}

	clientIP := c.ClientIP()

	product, err := pc.productService.GetProductDistance(c, productID, clientIP)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	rs.SuccessResponse(c, product)
}
