package controller

import (
	"fmt"
	"stock-management/internal/models"
	"stock-management/internal/services"
	rs "stock-management/pkgs/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

var Product = new(ProductController)

type ProductController struct{}

// GetProductList dynamic filtering and incremental search
// @Summary dynamic filtering and incremental search
// @Description Returns a list of products matching the search request
// @Tags Product
// @Accept  json
// @Produce  json
// @Param request body models.ProductSearchReq true "Product search details"
// @Success 200 {object} models.ProductSearchRp
// @Router /api/product/list [post]
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

	if req.Limit == 0 {
		req.Limit = 20
	}

	products, err := services.Service.ProductService.GetProductList(c, req)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	rs.SuccessResponse(c, products)
}

// CreateProduct creates a new product
// @Summary Create new product
// @Description Adds a new product to the inventory
// @Tags Product
// @Accept  json
// @Produce  json
// @Param request body models.ProductCreateReq true "Product creation details"
// @Success 200 {object} models.Product
// @Router /api/product/create [post]
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

	product, err := services.Service.ProductService.CreateProduct(c, req)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	rs.SuccessResponse(c, product)
}

// UpdateProduct updates an existing product
// @Summary Update product details
// @Description Modifies the details of an existing product
// @Tags Product
// @Accept  json
// @Produce  json
// @Param request body models.ProductUpdateReq true "Product update details"
// @Success 200 {object} models.Product
// @Router /api/product/update [post]
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

	product, err := services.Service.ProductService.UpdateProduct(c, req)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	rs.SuccessResponse(c, product)
}

// GetProduct retrieves a product by its ID
// @Summary Retrieve product by ID
// @Description Returns details of a product specified by its ID
// @Tags Product
// @Accept  json
// @Produce  json
// @Param request body models.ProductByIdReq true "Product ID details"
// @Success 200 {object} models.Product
// @Router /api/product/get [post]
func (pc *ProductController) GetProduct(c *gin.Context) {
	var req models.ProductByIdReq

	if err := c.ShouldBindJSON(&req); err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		rs.FailResponseWithCode(c, rs.ErrInvalidProduct)
		return
	}
	product, err := services.Service.ProductService.GetProduct(c, productID)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}

	rs.SuccessResponse(c, product)
}

// ExportProductsToPDF exports the product list to a PDF file
// @Summary Export product list to PDF
// @Description Generates and returns a PDF file containing the list of products
// @Tags Product
// @Accept  json
// @Produce  application/pdf
// @Param request body models.ProductSearchReq true "Product search details"
// @Success 200 {file} application/pdf "Product List PDF"
// @Router /api/product/export [post]
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

	products, err := services.Service.ProductService.GetProductList(c, req)
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
	today := time.Now().Format("20060102")
	fileName := fmt.Sprintf("products_%s.pdf", today)

	err = pdf.OutputFileAndClose(fileName)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	c.File(fileName)
}

// GetProductDistance Calculate the distance
// @Summary Calculate the distance
// @Description Returns the distance of a product specified by its ID
// @Tags Product
// @Accept  json
// @Produce  json
// @Param request body models.ProductDistanceByIdReq true "Product distance request details"
// @Success 200 {object} models.ProductDistanceRp
// @Router /api/product/distance [post]
func (pc *ProductController) GetProductDistance(c *gin.Context) {
	var req models.ProductDistanceByIdReq

	if err := c.ShouldBindJSON(&req); err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		rs.FailResponseWithCode(c, rs.ErrInvalidProduct)
		return
	}

	clientIP := c.ClientIP()
	if req.Ip != "" {
		clientIP = req.Ip
	}

	product, err := services.Service.ProductService.GetProductDistance(c, productID, clientIP)
	if err != nil {
		rs.FailResponseWithMessage(c, err.Error())
		return
	}
	rs.SuccessResponse(c, product)
}
