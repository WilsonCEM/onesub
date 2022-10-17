package commonMiddleware

import (
	"fmt"
	"strings"
	commonConfig "zimuzu/common/config"
	commonModel "zimuzu/common/models"

	"github.com/gin-gonic/gin"
)

func AuthContext(roleType commonModel.UserRoleType) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("authorization")
		handleAuthContext(c, authorization, roleType)
		c.Next()
	}
}

func handleAuthContext(c *gin.Context, token string, roleType commonModel.UserRoleType) {
	fmt.Println("*****************", token)
	if strings.TrimSpace(token) == "" {
		commonConfig.ErrorResponse(c, &commonConfig.ErrorStruct{ErrorCode: commonConfig.AUTH_INVALID, Message: "请登陆后再试"})
		c.Abort()
	}
	jwtModel, err := commonModel.JWTParse(token)
	if err != nil {
		commonConfig.ErrorResponse(c, &commonConfig.ErrorStruct{ErrorCode: commonConfig.TOKEN_INVALID, Message: "token无效"})
		c.Abort()
		return
	}
	if jwtModel.UserRole > roleType {
		commonConfig.ErrorResponse(c, &commonConfig.ErrorStruct{ErrorCode: commonConfig.AUTH_NOTEGNORE, Message: "权限不足"})
		c.Abort()
		return
	}
	c.Set("x-user", jwtModel.JWTPayLoad)
}
