package commonUtils

import (
	"fmt"
	commonConfig "zimuzu/common/config"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// 参数校验
func Validator(c *gin.Context, e *commonConfig.ErrorStruct, json interface{}) bool {
	err := c.ShouldBindBodyWith(&json, binding.JSON)
	for k, v := range c.Request.PostForm {
		fmt.Printf("k:%v\n", k)

		fmt.Printf("v:%v\n", v)

	}
	if err != nil {
		fmt.Println("**************************", err)
		commonConfig.ErrorError(e, commonConfig.PARAMETER_ERROR, err.Error())
		commonConfig.ErrorResponse(c, e)

		return false
	}

	commonConfig.ErrorSuccess(e)
	return true
}
