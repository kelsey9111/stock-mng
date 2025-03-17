package response

type RespCode uint32

const (
	OkCode        RespCode = 200
	ErrInternal   RespCode = 500
	ErrBadRequest RespCode = 400

	ErrInvalidDate          RespCode = 20001
	ErrInvalidName          RespCode = 30001
	ErrInvalidCategory      RespCode = 30002
	ErrInvalidSupplier      RespCode = 30003
	ErrInvalidProduct       RespCode = 30004
	ErrInvalidStatus        RespCode = 30005
	ErrInvalidStockLocation RespCode = 30006
	ErrInvalidReference     RespCode = 30007
)

var msg = map[RespCode]string{
	OkCode:        "success",
	ErrInternal:   "Internal server error",
	ErrBadRequest: "Bad qequest",

	ErrInvalidDate:          "Date is invalid",
	ErrInvalidName:          "Name is invalid",
	ErrInvalidCategory:      "Category is invalid",
	ErrInvalidSupplier:      "Supplier is invalid",
	ErrInvalidProduct:       "Product is invalid",
	ErrInvalidStatus:        "Status is invalid",
	ErrInvalidStockLocation: "Stock Location is invalid",
	ErrInvalidReference:     "Reference  is invalid",
}
