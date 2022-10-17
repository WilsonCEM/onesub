package service

import (
	"fmt"
	"time"
	"zimuzu/admin/instance"
	"zimuzu/admin/models"
	commonConfig "zimuzu/common/config"
	commonModel "zimuzu/common/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateResource(gtr commonModel.CreateResourceRequestBody, e *commonConfig.ErrorStruct) commonModel.ResourceModel {
	resourceModel := models.InitResourceModel(gtr)
	hasSeries := HasSeries(gtr.SeriesId)
	if !hasSeries {
		commonConfig.ErrorError(e, commonConfig.NOTFOUNT_ERROR, "系列剧集不存在")
		return resourceModel
	}
	if commonConfig.HasError(e) {
		return resourceModel
	}
	result := instance.ZimuzuDB.Debug().Create(&resourceModel)
	fmt.Println("********************", resourceModel.ResourceFile)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return resourceModel
	}
	return resourceModel
}

func UpdateResource(gtr commonModel.ResourceModel, e *commonConfig.ErrorStruct) {
	hasSeries := HasSeries(gtr.SeriesId)
	if !hasSeries {
		commonConfig.ErrorError(e, commonConfig.NOTFOUNT_ERROR, "系列剧集不存在")
		return
	}
	result := instance.ZimuzuDB.Debug().Save(&gtr)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
	}

}

//删除字幕
func DeleteResource(gtr commonModel.DeleteResourceBody, e *commonConfig.ErrorStruct) {
	result := instance.ZimuzuDB.Begin()
	result.Debug().Where("id=?", gtr.ID).Delete(new(commonModel.ResourceModel))
	result.Debug().Where("resource_id=?", gtr.ID).Delete(new(commonModel.DownloadLogModel))
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		result.Rollback()
	}
	result.Commit()
}

func DownloadResource(gtr commonModel.DownloadLogModel, c *gin.Context) {
	var e commonConfig.ErrorStruct
	res := HasResource(gtr.ResourceId, &e)
	if commonConfig.HasError(&e) {

		commonConfig.ErrorResponse(c, &e)
		return
	}
	res.DownloadTimes += 1
	result := instance.ZimuzuDB.Create(&gtr)
	if result.Error != nil {
		commonConfig.ErrorError(&e, commonConfig.DB_ERROR, result.Error.Error())
		return
	}
	fmt.Println("********************", res.DownloadTimes)
	instance.ZimuzuDB.Model(res).Update("download_times", res.DownloadTimes)
	filePath := commonConfig.RESOURCE_PATH + res.LoadFileName
	commonConfig.SuccessResponseFile(c, gtr.ID, filePath, res.ResourceFile)

}

//字幕是否存在
func HasResource(resourceId uint, e *commonConfig.ErrorStruct) commonModel.ResourceModel {
	var resourceModel = commonModel.ResourceModel{}
	result := instance.ZimuzuDB.First(&resourceModel, resourceId)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.NOTFOUNT_ERROR, "字幕不存在")
		return resourceModel
	}

	return resourceModel
}

//热门字幕
func FindHotResource(c *gin.Context) []commonModel.ResourceModel {
	var resourceModel = []commonModel.ResourceModel{}
	now := time.Now()
	old := now.AddDate(0, 0, -1)
	result1 := instance.ZimuzuDB.Table("download_log_model").
		Select("distinct(series_id) as sid,max(created_at) as dtime,count(*) as dcount").
		Group("sid")
	result2 := instance.ZimuzuDB.Table("(?) as t", result1).Where("dtime>=?", old)
	result := instance.ZimuzuDB.Limit(5).Order("created_at desc").Preload("Series").Preload("User").
		Joins("left join (?) as l on series_id=l.sid", result2).Order("dcount desc,dtime desc").
		Find(&resourceModel)

	var e *commonConfig.ErrorStruct
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.NOTFOUNT_ERROR, "字幕不存在")
		return resourceModel
	}

	return resourceModel
}

//字幕库查询
func FindSubLibrary(gtr *models.FindSeriesRequestBody, e *commonConfig.ErrorStruct) []commonModel.SubLibraryBody {
	resourceModels := []commonModel.SubLibraryBody{}
	// var result1 *gorm.DB
	var result *gorm.DB
	fmt.Println("******************", gtr.PageNumber)
	now := time.Now()
	old := now.AddDate(0, 0, -1)
	old2 := now.AddDate(0, 0, -7)
	result = instance.ZimuzuDB.Table("resource_model").
		Select("series_id,MAX(resource_model.created_at) as ltime,cnname,origin_name,count(resource_model.id) as cnumber" +
			",cover,archive,number_seasons,series_model.desc,series_model.score,series_model.series_type,series_model.genres,series_model.banner" +
			",series_model.screen_time").
		Joins("join series_model on resource_model.series_id=series_model.id").
		Group("series_id").
		Limit(commonConfig.Pageing_rows).
		Offset((int(gtr.PageNumber) - 1) * commonConfig.Pageing_rows)
	if gtr.SeriesType != 0 {
		result.Where("series_model.series_type=?", gtr.SeriesType)
	}
	switch gtr.T {
	case 2:
		result1 := instance.ZimuzuDB.Table("download_log_model").
			Select("distinct(series_id) as sid,max(created_at) as dtime,count(*) as dcount").
			Group("sid")
		result2 := instance.ZimuzuDB.Table("(?) as t", result1).Where("dtime>=?", old)
		result.Joins("left join (?) as l on series_model.id=l.sid", result2).Order("dcount desc,ltime desc")

	case 3:
		result1 := instance.ZimuzuDB.Table("download_log_model").
			Select("distinct(series_id) as sid,max(created_at) as dtime,count(*) as dcount").
			Group("sid")
		result2 := instance.ZimuzuDB.Table("(?) as t", result1).Where("dtime>=?", old2)
		result.Joins("left join (?) as l on series_model.id=l.sid", result2).Order("dcount desc")
	case 1:
		result.Order("series_model.screen_time desc")
	case 0:
		result.Order("ltime desc")
	}

	result.Debug().Find(&resourceModels)
	ChangeSubImage(resourceModels, e)

	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return resourceModels
	}
	return resourceModels
}

func ResourceDetail(gtr models.FindSeriesRequestBody, e *commonConfig.ErrorStruct) commonModel.ResourceModel {
	var res commonModel.ResourceModel
	result := instance.ZimuzuDB.Where("id=?", gtr.ID).Preload("User").Preload("User.SubGroups").First(&res)
	res.User.SubGroups.LogoPath = commonConfig.BASE_SUBICONPATH + res.User.SubGroups.LogoPath
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return res
	}
	return res
}

//英美剧字幕查询
func FindSeriesResourceByTV(gtr commonModel.FindSeriesResourceBody, e *commonConfig.ErrorStruct) []commonModel.ResponseSeriesResourceBody {
	resourceModels := []commonModel.ResponseSeriesResourceBody{}
	result := instance.ZimuzuDB.Debug().Table("resource_model").
		Select("resource_model.id,resource_file,language,origin,translator,download_times,"+
			"username,sub_group_id,group_name,source_title,user_id,series_no,source_title,"+
			"series_id,s.cnname,s.series_type,s.number_seasons,format,remarks").
		Joins("join user_model on resource_model.user_id=user_model.id").
		Joins("join sub_group_model on user_model.sub_group_id=sub_group_model.id").
		Joins("join series_model as s on resource_model.series_id=s.id").
		// Joins("join series_model on resource_model.series_id=series_model.id").
		Where("series_id = ? and resource_model.deleted_at is null", gtr.SeriesId).
		Group("series_no,resource_model.id").
		Order("series_no desc,resource_model.updated_at desc")

	// if gtr.SubGroupId != 0 {
	// 	result.Where("sub_group_id=?", gtr.SubGroupId)
	// }
	result.Find(&resourceModels)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return resourceModels
	}
	return resourceModels
}

func FindSeriesResource(gtr commonModel.FindSeriesResourceBody, e *commonConfig.ErrorStruct) []commonModel.ResponseSeriesResourceBody {
	resourceModels := []commonModel.ResponseSeriesResourceBody{}
	result := instance.ZimuzuDB.Table("resource_model").
		Select("resource_model.id,resource_file,language,origin,translator,download_times,"+
			"username,sub_group_id,group_name,source_title,user_id,series_no,source_title,"+
			"series_id,s.cnname,s.series_type,s.number_seasons,format,remarks").
		Joins("join series_model as s on resource_model.series_id=s.id").
		Joins("join user_model on resource_model.user_id=user_model.id").
		Joins("join sub_group_model on user_model.sub_group_id=sub_group_model.id").
		Where("series_id = ?", gtr.SeriesId).
		Order("resource_model.updated_at desc")

	// if gtr.SubGroupId != 0 {
	// 	result.Where("sub_group_id=?", gtr.SubGroupId)
	// }
	result.Find(&resourceModels)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return resourceModels
	}
	return resourceModels
}
func ChangeSubImage(seriesModels []commonModel.SubLibraryBody, e *commonConfig.ErrorStruct) {
	for i, _ := range seriesModels {
		seriesModels[i].Cover = commonConfig.IMAGE_PATH + seriesModels[i].Cover
		if seriesModels[i].Banner != "" {
			seriesModels[i].Banner = commonConfig.IMAGE_PATH + seriesModels[i].Banner
		}

		instance.ZimuzuDB.Table("resource_model").Select("group_name").
			Joins("join user_model on resource_model.user_id= user_model.id").
			Joins("join sub_group_model on user_model.sub_group_id=sub_group_model.id").
			Where("series_id=?", seriesModels[i].SeriesId).Order("resource_model.updated_at desc").
			First(&seriesModels[i])
		// if result1.Error != nil {
		// 	commonConfig.ErrorError(e, commonConfig.DB_ERROR, result1.Error.Error())
		// }
	}
}
