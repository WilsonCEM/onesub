package controller

import (
	"fmt"
	"math/rand"
	"path"
	"strconv"
	"time"
	"zimuzu/admin/models"
	"zimuzu/admin/service"
	commonConfig "zimuzu/common/config"
	commonModel "zimuzu/common/models"
	commonUtils "zimuzu/common/utils"

	"github.com/gin-gonic/gin"
)

func ResourceByPost(c *gin.Context, e *commonConfig.ErrorStruct, gtr *commonModel.CreateResourceRequestBody) bool {

	sid, err := strconv.Atoi(c.PostForm("seriesId"))
	if err != nil {
		fmt.Println("**************************1", err)
		commonConfig.ErrorError(e, commonConfig.PARAMETER_ERROR, err.Error())
		commonConfig.ErrorResponse(c, e)

	} else {
		gtr.SeriesId = uint(sid)
	}
	gtr.SourceTitle = c.PostForm("sourcetitle")
	gtr.Format = c.PostForm("format")
	gtr.Origin = c.PostForm("origin")
	gtr.Language = c.PostForm("language")
	gtr.Translator = c.PostForm("translator")
	sn, err := strconv.Atoi(c.PostForm("seriesNo"))
	if err != nil {
		gtr.SeriesNo = 0
	} else {
		gtr.SeriesNo = uint16(sn)
	}

	gtr.Remarks = c.PostForm("remarks")

	// if err != nil {
	// 	fmt.Println("**************************2", err)
	// 	commonConfig.ErrorError(e, commonConfig.PARAMETER_ERROR, err.Error())

	// 	return false
	// }

	commonConfig.ErrorSuccess(e)
	return true
}

func DeleteResource(c *gin.Context) {
	var gtr commonModel.DeleteResourceBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}

	service.DeleteResource(gtr, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, gtr)
}

// func UpdateResourceByPost(c *gin.Context, e *commonConfig.ErrorStruct, gtr *commonModel.UpdateResourceRequestBody) bool {

// 	sid, err := strconv.Atoi(c.PostForm("seriesId"))
// 	if err != nil {
// 		fmt.Println("**************************1", err)
// 		commonConfig.ErrorError(e, commonConfig.PARAMETER_ERROR, err.Error())
// 		commonConfig.ErrorResponse(c, e)

// 	} else {
// 		gtr.SeriesId = uint(sid)
// 	}
// 	id, err := strconv.Atoi(c.PostForm("id"))
// 	if err != nil {
// 		fmt.Println("**************************2", err)
// 		commonConfig.ErrorError(e, commonConfig.PARAMETER_ERROR, err.Error())
// 		commonConfig.ErrorResponse(c, e)

// 	} else {
// 		gtr.ID = uint(id)
// 	}

// 	gtr.Format = c.PostForm("format")
// 	gtr.Origin = c.PostForm("origin")
// 	gtr.Language = c.PostForm("language")
// 	gtr.Translator = c.PostForm("translator")
// 	sn, err := strconv.Atoi(c.PostForm("episodecount"))
// 	if err != nil {
// 		gtr.SeriesNo = 0
// 	} else {
// 		gtr.SeriesNo = uint16(sn)
// 	}

// 	gtr.Remarks = c.PostForm("remarks")

// 	if err != nil {
// 		fmt.Println("**************************3", err)
// 		commonConfig.ErrorError(e, commonConfig.PARAMETER_ERROR, err.Error())

// 		return false
// 	}
// 	commonConfig.ErrorSuccess(e)
// 	return true
// }

func UpdateResource(c *gin.Context) {
	var gtr commonModel.CreateResourceRequestBody
	var e commonConfig.ErrorStruct

	// user, _ := c.Get("x-user")
	// userid := user.(commonModel.JWTPayLoad)
	// gtr.UserId = userid.Uid
	// gtr.UserId = 45

	if !ResourceByPost(c, &e, &gtr) {
		return
	}
	ids, err := strconv.Atoi(c.PostForm("id"))
	var rid uint
	if err != nil {
		fmt.Println("**************************2", err)
		commonConfig.ErrorError(&e, commonConfig.PARAMETER_ERROR, err.Error())
		commonConfig.ErrorResponse(c, &e)
		return

	} else {
		rid = uint(ids)
	}
	resModel := service.HasResource(rid, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorError(&e, commonConfig.PARAMETER_ERROR, err.Error())
		commonConfig.ErrorResponse(c, &e)
		return
	}

	//FormFile返回所提供的表单键的第一个文件
	f, err := c.FormFile("resourceFile")

	if err == nil {
		gtr.ResourceFile = f.Filename
		name := CreateUploadFile(gtr.ResourceFile)
		gtr.LoadFileName = name
		err1 := c.SaveUploadedFile(f, "./static/resource/"+name)
		if err1 != nil {
			commonConfig.ErrorError(&e, commonConfig.CONDITION_ERROR, err.Error())
			commonConfig.ErrorResponse(c, &e)
			return
		}
	} else {
		gtr.ResourceFile = resModel.ResourceFile
		gtr.LoadFileName = resModel.LoadFileName
	}
	gtr.UserId = resModel.UserId
	resModel.Resource.CreateResourceRequestBody = gtr
	//SaveUploadedFile上传表单文件到指定的路径

	service.UpdateResource(resModel, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}

	commonConfig.SuccessResponse(c, gtr)
}

func ResourceDetail(c *gin.Context) {
	var gtr models.FindSeriesRequestBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}
	seriesModels := service.ResourceDetail(gtr, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, seriesModels)
}

func CreateResource(c *gin.Context) {
	var gtr commonModel.CreateResourceRequestBody
	var e commonConfig.ErrorStruct

	user, _ := c.Get("x-user")
	userid := user.(commonModel.JWTPayLoad)
	gtr.UserId = userid.Uid
	// gtr.UserId = 45

	if !ResourceByPost(c, &e, &gtr) {
		return
	}
	//FormFile返回所提供的表单键的第一个文件
	f, err := c.FormFile("resourceFile")
	if err != nil {
		commonConfig.ErrorError(&e, commonConfig.CONDITION_ERROR, err.Error())
		commonConfig.ErrorResponse(c, &e)
		return
	}
	gtr.ResourceFile = f.Filename
	name := CreateUploadFile(gtr.ResourceFile)
	gtr.LoadFileName = name

	//SaveUploadedFile上传表单文件到指定的路径
	err = c.SaveUploadedFile(f, "./static/resource/"+name)
	if err != nil {
		commonConfig.ErrorError(&e, commonConfig.CONDITION_ERROR, err.Error())
		commonConfig.ErrorResponse(c, &e)
		return
	}
	service.CreateResource(gtr, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}

	commonConfig.SuccessResponse(c, gtr)
}

//字幕库
func SubLibrary(c *gin.Context) {
	var gtr models.FindSeriesRequestBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}
	seriesModels := service.FindSubLibrary(&gtr, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, seriesModels)
}

//英美剧字幕分组查询
func SeriesResourceByTV(c *gin.Context) {
	var gtr commonModel.FindSeriesResourceBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}

	seriesModels := service.FindSeriesResourceByTV(gtr, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, seriesModels)
}

func CreateUploadFile(name string) string {
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
	fmt.Println("******************", file_name)
	// c.SaveUploadedFile(file, basePath+file_name)
	return file_name

}

//下载字幕
func DownloadResource(c *gin.Context) {
	var gtr commonModel.DownloadLogModel
	// gtr.SeriesID = 152
	// gtr.ResourceId = 92
	var e commonConfig.ErrorStruct
	sid, err := strconv.Atoi(c.PostForm("resourceId"))
	if err != nil {
		fmt.Println("**************************", err)
		commonConfig.ErrorError(&e, commonConfig.PARAMETER_ERROR, err.Error())
		commonConfig.ErrorResponse(c, &e)
		return

	}
	gtr.ResourceId = uint(sid)
	// if !commonUtils.Validator(c, &e, &gtr) {
	// 	return
	// }
	r := service.HasResource(gtr.ResourceId, &e)
	if !commonConfig.HasError(&e) {
		gtr.SeriesID = r.SeriesId
	}
	authorization := c.Request.Header.Get("authorization")
	if authorization != "" {
		jmodel, _ := commonModel.JWTParse(authorization)
		gtr.UserId = jmodel.JWTPayLoad.Uid
	}
	service.DownloadResource(gtr, c)
}

func HotResSource(c *gin.Context) {
	hotResource := service.FindHotResource(c)
	commonConfig.SuccessResponse(c, hotResource)
}

func SeriesResource(c *gin.Context) {
	var gtr commonModel.FindSeriesResourceBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}
	resourceModels := service.FindSeriesResource(gtr, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, resourceModels)
}
