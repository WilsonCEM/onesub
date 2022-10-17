package controller

import (
	"fmt"
	"zimuzu/admin/service"
	commonConfig "zimuzu/common/config"
	commonModel "zimuzu/common/models"
	commonUtils "zimuzu/common/utils"

	"github.com/gin-gonic/gin"
)

//注册
func CreateUser(c *gin.Context) {
	var gtr commonModel.CreateUserRequestBody
	var e commonConfig.ErrorStruct

	if !commonUtils.Validator(c, &e, &gtr) {
		fmt.Println("1231231231231")
		return
	}
	if gtr.Email == "" || gtr.Username == "" {
		commonConfig.ErrorResponse(c, &commonConfig.ErrorStruct{ErrorCode: commonConfig.PARAMETER_ERROR, Message: "用户以及邮箱必填"})
		return
	}
	var userModel = service.CreateUser(gtr, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	var token = service.JWTSign(&userModel, &e)
	flag := service.CreateEmail(token, userModel, &e)
	if !flag {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}

	commonConfig.SuccessResponse(c, token)
}

//查找用户
func FindUserByName(c *gin.Context) {
	var gtr commonModel.FindUerBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		fmt.Println("1231231231231")
		return
	}
	if gtr.UserName == "" {
		//commonConfig.ErrorResponse(c, &commonConfig.ErrorStruct{ErrorCode: commonConfig.PARAMETER_ERROR, Message: "用户以及邮箱必填"})
		return
	}
	authorization := c.Request.Header.Get("authorization")
	if authorization != "" {
		userModel := service.FindUserByToken(authorization, &e)
		userList := service.FindUserByName(gtr, userModel, &e)
		if commonConfig.HasError(&e) {
			commonConfig.ErrorResponse(c, &e)
			return
		}
		commonConfig.SuccessResponse(c, userList)
	} else {
		commonConfig.ErrorError(&e, commonConfig.CONDITION_ERROR, "请先登录！")
		commonConfig.ErrorResponse(c, &e)
		return
	}

}

//添加用户到字幕组
func UserToSubGroup(c *gin.Context) {
	var gtr commonModel.UserToSubGroupBody
	var e commonConfig.ErrorStruct

	if !commonUtils.Validator(c, &e, &gtr) {
		fmt.Println("1231231231231")
		return
	}

	service.UserToSubGroup(gtr, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}

	commonConfig.SuccessResponse(c, gtr)
}

//将用户移除字幕组
func SubGroupRemoveUser(c *gin.Context) {
	var gtr commonModel.UserToSubGroupBody
	var e commonConfig.ErrorStruct
	flag := false
	if !commonUtils.Validator(c, &e, &gtr) {
		fmt.Println("1231231231231")
		return
	}

	authorization := c.Request.Header.Get("authorization")
	if authorization != "" {
		userModel := service.FindUserByToken(authorization, &e)
		if userModel.UserRole < commonModel.USER_ROLE_SUBGROUP {
			flag = true
		}
		if userModel.UserRole == commonModel.USER_ROLE_SUBGROUP {
			if userModel.SubGroupId == gtr.SubGroupId {
				flag = true
			}
		}
	} else {
		commonConfig.ErrorError(&e, commonConfig.CONDITION_ERROR, "请先登录！")
		commonConfig.ErrorResponse(c, &e)
	}
	if flag {
		service.SubGroupRemoveUser(gtr, &e)
		if commonConfig.HasError(&e) {
			commonConfig.ErrorResponse(c, &e)
			return
		}
		commonConfig.SuccessResponse(c, gtr)
	} else {
		commonConfig.ErrorError(&e, commonConfig.CONDITION_ERROR, "该用户操作权限不足")
		commonConfig.ErrorResponse(c, &e)
	}

}

//发送找回密码邮件
func SendRetrievePassword(c *gin.Context) {
	var gtr commonModel.UserLoginRequestBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}
	var userModel commonModel.UserModel
	if gtr.LoginBy == commonModel.LOGIN_BY_EMAIL {
		userModel = service.QueryUserByEmail(&gtr, &e)
	} else {
		userModel = service.QueryUserByUserName(&gtr, &e)
	}
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	var token = service.JWTSign(&userModel, &e)
	gtr.UserToken = token
	gtr.ID = userModel.ID
	gtr.Email = userModel.Email
	service.SendRetrievePasswordService(gtr, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, "邮件发送成功")
}

//找回密码验证
func RetrievePassword(c *gin.Context) {
	var e commonConfig.ErrorStruct
	token := c.Query("t")
	fmt.Println("**********************", token)
	jmodel, err := commonModel.JWTParse(token)
	if err != nil {
		commonConfig.ErrorError(&e, commonConfig.EMAIL_SENDERROR, err.Error())
		commonConfig.ErrorResponse(c, &e)
		return
	}
	service.RetrievePasswordService(token, jmodel.JWTPayLoad.Uid, &e)

	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, token)
}

//修改密码
func ChangePassword(c *gin.Context) {
	var e commonConfig.ErrorStruct
	authorization := c.Request.Header.Get("authorization")
	var gtr commonModel.ChangePasswordBody
	var userModel commonModel.UserModel
	jmodel, err := commonModel.JWTParse(authorization)
	if err != nil {
		commonConfig.ErrorError(&e, commonConfig.EMAIL_SENDERROR, err.Error())
		commonConfig.ErrorResponse(c, &e)
		return
	}
	if !commonUtils.Validator(c, &e, &gtr) {
		fmt.Println("1231231231231")
		return
	}
	gtr.UserID = jmodel.JWTPayLoad.Uid
	fmt.Println("****************", jmodel, authorization)
	userModel = service.QueryUserByID(&gtr, &e)
	comparePasswords := commonUtils.CompareHashPassword(userModel.Password, []byte(gtr.Password))
	if comparePasswords {
		newpass, _ := commonUtils.HashAndSalt(gtr.NewPassword)

		userModel.Password = newpass
		service.ChangePasswordService(&userModel, &e)

		if commonConfig.HasError(&e) {
			commonConfig.ErrorResponse(c, &e)
			return
		}
		var token = service.JWTSign(&userModel, &e)
		commonConfig.SuccessResponse(c, token)
	} else {
		e.ErrorCode = commonConfig.CONDITION_ERROR
		e.Message = "用户密码异常"
		commonConfig.ErrorResponse(c, &e)
	}

}

//重置密码
func ResetPassword(c *gin.Context) {
	var gtr commonModel.ChangePasswordBody
	var e commonConfig.ErrorStruct
	var userModel commonModel.UserModel
	authorization := c.Request.Header.Get("authorization")
	if authorization == "" {
		commonConfig.ErrorResponse(c, &commonConfig.ErrorStruct{ErrorCode: commonConfig.PARAMETER_ERROR, Message: "参数错误"})
		return
	}
	if !commonUtils.Validator(c, &e, &gtr) {
		fmt.Println("1231231231231")
		return
	}
	userModel = service.QueryUserByID(&gtr, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	if gtr.UserID == 0 || gtr.NewPassword == "" {
		commonConfig.ErrorResponse(c, &commonConfig.ErrorStruct{ErrorCode: commonConfig.PARAMETER_ERROR, Message: "找不到该用户或新密码为空！"})
		return
	}
	newpass, _ := commonUtils.HashAndSalt(gtr.NewPassword)

	userModel.Password = newpass
	service.ResetPasswordService(&userModel, authorization, &e)
	var token = service.JWTSign(&userModel, &e)

	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}

	commonConfig.SuccessResponse(c, token)

}

//登录
func Login(c *gin.Context) {
	var gtr commonModel.UserLoginRequestBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}
	var userModel commonModel.UserModel
	if gtr.LoginBy == commonModel.LOGIN_BY_EMAIL {
		userModel = service.QueryUserByEmail(&gtr, &e)
	} else {
		userModel = service.QueryUserByUserName(&gtr, &e)
	}
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	comparePasswords := commonUtils.CompareHashPassword(userModel.Password, []byte(gtr.Password))
	if comparePasswords {
		var token = service.JWTSign(&userModel, &e)
		commonConfig.SuccessResponse(c, token)
	} else {
		e.ErrorCode = commonConfig.CONDITION_ERROR
		e.Message = "用户密码异常"
		commonConfig.ErrorResponse(c, &e)
	}
}

//激活用户
func UserActivation(c *gin.Context) {
	var e commonConfig.ErrorStruct
	token := c.Query("t")
	fmt.Println("**********************", token)
	jmodel, err := commonModel.JWTParse(token)
	if err != nil {
		commonConfig.ErrorError(&e, commonConfig.EMAIL_SENDERROR, err.Error())
		commonConfig.ErrorResponse(c, &e)
		return
	}
	service.ActivationUser(token, jmodel.JWTPayLoad.Uid, &e)

	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, token)
}
