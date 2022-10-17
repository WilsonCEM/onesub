package controller

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"
	"zimuzu/admin/service"
	commonConfig "zimuzu/common/config"
	commonModel "zimuzu/common/models"
	commonUtils "zimuzu/common/utils"

	"github.com/gin-gonic/gin"
)

func FindSubGroupAll(c *gin.Context) {
	// var gtr commonModel.QuerySubGroupBody
	var e commonConfig.ErrorStruct
	// if !commonUtils.Validator(c, &e, &gtr) {
	// 	return
	// }
	sgList := service.FindSubGroupAll(&e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, sgList)
}

func SubGroupByPost(c *gin.Context, e *commonConfig.ErrorStruct, gtr *commonModel.CreateSubGroupRequestBody) {

	gtr.GroupName = c.PostForm("groupname")
	gtr.Desc = c.PostForm("desc")
	gtr.Blog = c.PostForm("blog")
	gtr.Website = c.PostForm("website")

}
func SubGroupByPostID(c *gin.Context, e *commonConfig.ErrorStruct, gtr *commonModel.UpdateSubGroupRequestBody) {

	sid, err := strconv.Atoi(c.PostForm("groupId"))
	if err != nil {
		fmt.Println("**************************", err)
		commonConfig.ErrorError(e, commonConfig.PARAMETER_ERROR, err.Error())
		commonConfig.ErrorResponse(c, e)

	} else {
		gtr.GroupId = uint(sid)
	}
	gtr.GroupName = c.PostForm("groupname")
	gtr.Desc = c.PostForm("desc")
	gtr.Blog = c.PostForm("blog")
	gtr.Website = c.PostForm("website")

}

//创建字幕组
func CreateSubGroup(c *gin.Context) {
	var gtr commonModel.CreateSubGroupRequestBody
	var e commonConfig.ErrorStruct
	SubGroupByPost(c, &e, &gtr)
	hasSubGroup := service.HasSubGroup(gtr.GroupName)
	if hasSubGroup {
		commonConfig.ErrorError(&e, commonConfig.DB_ERROR, "字幕组已存在")
		commonConfig.ErrorResponse(c, &e)
		return
	}
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	//FormFile返回所提供的表单键的第一个文件
	logo, err := c.FormFile("logo_path")
	if err != nil {
		commonConfig.ErrorError(&e, commonConfig.CONDITION_ERROR, err.Error())
		commonConfig.ErrorResponse(c, &e)
		return
	}
	gtr.LogoPath = CreateUploadFile(logo.Filename)
	//SaveUploadedFile上传表单文件到指定的路径
	c.SaveUploadedFile(logo, "./static/icon/"+gtr.LogoPath)
	if err != nil {
		commonConfig.ErrorError(&e, commonConfig.CONDITION_ERROR, err.Error())
		commonConfig.ErrorResponse(c, &e)
		return
	}
	wechat, err := c.FormFile("wechat_path")
	if err != nil {
		commonConfig.ErrorError(&e, commonConfig.CONDITION_ERROR, err.Error())
		commonConfig.ErrorResponse(c, &e)
		return
	}
	gtr.WechatPath = CreateUploadFile(wechat.Filename)
	//SaveUploadedFile上传表单文件到指定的路径
	c.SaveUploadedFile(wechat, "./static/icon/"+gtr.WechatPath)
	if err != nil {
		commonConfig.ErrorError(&e, commonConfig.CONDITION_ERROR, err.Error())
		commonConfig.ErrorResponse(c, &e)
		return
	}
	service.CreateSubGroup(gtr, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, gtr)
}

//修改字幕组信息
func SubgroupUpdate(c *gin.Context) {
	var gtr commonModel.UpdateSubGroupRequestBody
	var e commonConfig.ErrorStruct
	SubGroupByPostID(c, &e, &gtr)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	authorization := c.Request.Header.Get("authorization")
	if authorization != "" {
		userModel := service.FindUserByToken(authorization, &e)
		if userModel.UserRole > commonModel.USER_ROLE_SUBGROUP {
			commonConfig.ErrorError(&e, commonConfig.CONDITION_ERROR, "您的权限不足！")
			commonConfig.ErrorResponse(c, &e)
			return
		}
		if userModel.UserRole == commonModel.USER_ROLE_SUBGROUP {
			if userModel.SubGroupId != gtr.GroupId {
				commonConfig.ErrorError(&e, commonConfig.CONDITION_ERROR, "您的权限不足！")
				commonConfig.ErrorResponse(c, &e)
				return
			}
		}
	} else {
		commonConfig.ErrorError(&e, commonConfig.CONDITION_ERROR, "请先登录！")
		commonConfig.ErrorResponse(c, &e)
		return
	}
	subGroup, hasSubGroup := service.HasSubGroupByID(gtr.GroupId)

	logoFile := subGroup.LogoPath
	wechatFile := subGroup.WechatPath
	if hasSubGroup {
		commonConfig.ErrorError(&e, commonConfig.DB_ERROR, "该字幕组不存在")
		commonConfig.ErrorResponse(c, &e)
		return
	}
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	subGroup.GroupName = gtr.GroupName
	subGroup.Desc = gtr.Desc
	subGroup.Blog = gtr.Blog
	subGroup.Website = gtr.Website
	if gtr.LogoPath != "" {
	}
	//FormFile返回所提供的表单键的第一个文件
	logo, err := c.FormFile("logo_path")
	if err == nil {
		subGroup.LogoPath = CreateUploadFile(logo.Filename)
		//SaveUploadedFile上传表单文件到指定的路径
		c.SaveUploadedFile(logo, "./static/icon/"+subGroup.LogoPath)
		if err != nil {
			commonConfig.ErrorError(&e, commonConfig.CONDITION_ERROR, err.Error())
			commonConfig.ErrorResponse(c, &e)
			return
		}
		os.Remove("./static/icon/" + logoFile)
	}

	wechat, err := c.FormFile("wechat_path")
	if err == nil {
		subGroup.WechatPath = CreateUploadFile(wechat.Filename)
		//SaveUploadedFile上传表单文件到指定的路径
		c.SaveUploadedFile(wechat, "./static/icon/"+subGroup.WechatPath)
		if err != nil {
			commonConfig.ErrorError(&e, commonConfig.CONDITION_ERROR, err.Error())
			commonConfig.ErrorResponse(c, &e)
			return
		}
		os.Remove("./static/icon/" + wechatFile)
	}

	fmt.Println("*******************", subGroup)
	service.UpdateSubGroup(&subGroup, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}

	commonConfig.SuccessResponse(c, gtr)
}

//热门字幕组
func HotSubGroup(c *gin.Context) {
	var e commonConfig.ErrorStruct
	seriesModels := service.FindSubGroup(&e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}

	commonConfig.SuccessResponse(c, seriesModels)
}

//字幕组列表
func SubgroupList(c *gin.Context) {
	var gtr commonModel.QuerySubGroupBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}
	sgList := service.SubGroupList(gtr, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, sgList)
}

//字幕组详情
func SubGroupDetail(c *gin.Context) {
	var gtr commonModel.QuerySubGroupBody
	var e commonConfig.ErrorStruct
	flag := false
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}
	authorization := c.Request.Header.Get("authorization")
	if authorization != "" {
		userModel := service.FindUserByToken(authorization, &e)
		if userModel.UserRole < commonModel.USER_ROLE_SUBGROUP {
			flag = true
		}
		if userModel.UserRole == commonModel.USER_ROLE_SUBGROUP {
			if userModel.SubGroupId == gtr.GroupId {
				flag = true
			}
		}
	}
	sgList := service.SubGroupDetailService(gtr, &e)

	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	sgList.Flag = flag
	commonConfig.SuccessResponse(c, sgList)
}

//字幕组主页影片
func SubGroupSeries(c *gin.Context) {
	var gtr commonModel.QuerySeriesSubGroupBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}
	sgList := service.SubGroupSeriesDetail(gtr, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, sgList)
}

//字幕组影片详情（英美剧）
func SubGroupSeriesTV(c *gin.Context) {
	var gtr commonModel.QuerySeriesSubGroupBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}
	sgList := service.SeriesListByTv(gtr, true, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, sgList)
}

//字幕组影片详情（电影）
func SubGroupSeriesMovie(c *gin.Context) {
	var gtr commonModel.QuerySeriesSubGroupBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}
	sgList := service.SeriesListByTv(gtr, false, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, sgList)
}

//上传图标文件名
func CreateSubGroupIcon(name string) string {
	// func CreateUploadFile(c *gin.Context) string {
	// var e commonConfig.ErrorStruct
	// file, err := c.FormFile("file")
	// if err != nil {
	// 	commonConfig.ErrorError(&e, commonConfig.FILEUPLOAD_ERROR, err.Error())
	// 	commonConfig.ErrorResponse(c, &e)
	// 	return ""
	// }

	// file_name := strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(rand.Intn(999999-100000)+100000) + path.Ext(file.Filename)
	file_name := strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(rand.Intn(999999-100000)+100000) + path.Ext(name)
	// c.SaveUploadedFile(file, basePath+file_name)
	return file_name
}
