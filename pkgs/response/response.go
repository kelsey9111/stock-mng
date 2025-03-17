package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code    RespCode    `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResponse(c *gin.Context, code RespCode, data interface{}) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    code,
		Message: msg[code],
		Data:    data,
	})
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    OkCode,
		Message: msg[OkCode],
		Data:    data,
	})
}
func FailResponseWithCode(c *gin.Context, code RespCode) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    code,
		Message: msg[code],
		Data:    nil,
	})
}
func FailResponseWithMessage(c *gin.Context, msgs string) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    ErrInternal,
		Message: msg[ErrInternal] + " " + msgs,
		Data:    nil,
	})
}

// func FailResponse(c *gin.Context, err error) {
// 	c.JSON(http.StatusOK, ResponseData{
// 		Code:    code,
// 		Message: msg[code],
// 		Data:    nil,
// 	})
// }
