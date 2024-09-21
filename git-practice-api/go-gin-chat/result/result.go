// 结构数据定义
package result

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Result结构体
type Result struct {
	Code    uint        `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// success
func Success(data interface{}, c *gin.Context) {
	if data == nil {
		data = gin.H{}
	}
	result := Result{}
	result.Code = APIcode.Success
	result.Message = APIcode.GetMessage(APIcode.Success)
	result.Data = data
	c.JSON(http.StatusOK, result)
}

//failure

func Failture(code uint, message string, c *gin.Context, err *error) {
	result := Result{}
	result.Code = code
	result.Message = message
	if err != nil {
		result.Data = *err
	} else {
		result.Data = gin.H{}
	}
	c.JSON(http.StatusOK, result)
}
