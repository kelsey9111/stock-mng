package response

type RespCode uint32

const (
	OkCode        RespCode = 200
	ErrInternal   RespCode = 500
	ErrBadRequest RespCode = 400

	Err                     RespCode = 3000
	ErrInvalidName          RespCode = 3001
	ErrInvalidCategory      RespCode = 3002
	ErrInvalidSupplier      RespCode = 3003
	ErrInvalidProduct       RespCode = 3004
	ErrInvalidStatus        RespCode = 3005
	ErrInvalidStockLocation RespCode = 3006
	ErrInvalidReference     RespCode = 3007
	ErrInvalidDate          RespCode = 2008
)

var msg = map[RespCode]string{
	OkCode:        "success",
	ErrInternal:   "Internal server error",
	ErrBadRequest: "Bad qequest",

	Err:                     "Error",
	ErrInvalidDate:          "Date is invalid",
	ErrInvalidName:          "Name is invalid",
	ErrInvalidCategory:      "Category is invalid",
	ErrInvalidSupplier:      "Supplier is invalid",
	ErrInvalidProduct:       "Product is invalid",
	ErrInvalidStatus:        "Status is invalid",
	ErrInvalidStockLocation: "Stock Location is invalid",
	ErrInvalidReference:     "Reference  is invalid",
}
