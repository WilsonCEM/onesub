package service

import (
	"fmt"
	"zimuzu/admin/instance"
	commonConfig "zimuzu/common/config"
	commonModel "zimuzu/common/models"
)

//查找字幕组
func FindSubGroup(e *commonConfig.ErrorStruct) []commonModel.SubGroupModel {
	subGroupModel := []commonModel.SubGroupModel{}

	result := instance.ZimuzuDB.Limit(20).Order("ID desc").Find(&subGroupModel)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return subGroupModel
	}
	ChangeImage(subGroupModel)
	return subGroupModel
}

//字幕组详情
func SubGroupDetailService(gtr commonModel.QuerySubGroupBody, e *commonConfig.ErrorStruct) commonModel.SubGroupByRoel {
	var subGroup commonModel.SubGroupModel

	result := instance.ZimuzuDB.Where("id=?", gtr.GroupId).Preload("Users").First(&subGroup)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
	}
	subGroup.WechatPath = commonConfig.BASE_SUBICONPATH + subGroup.WechatPath
	subGroup.LogoPath = commonConfig.BASE_SUBICONPATH + subGroup.LogoPath
	var subCount int64
	var dlCount commonModel.SumInt
	instance.ZimuzuDB.Table("resource_model as r").
		Joins("join user_model as u on r.user_id=u.id").
		Joins("join sub_group_model as sg on u.sub_group_id=sg.id").
		Where("sg.id=?", gtr.GroupId).
		Count(&subCount)

	instance.ZimuzuDB.Table("resource_model as r").
		Select("sum(r.downLoad_times) as sum_value").
		Joins("join user_model as u on r.user_id=u.id").
		Joins("join sub_group_model as sg on u.sub_group_id=sg.id").
		Where("sg.id=?", gtr.GroupId).
		Take(&dlCount)

	sgList := commonModel.SubGroupByRoel{
		SubGroup:      subGroup,
		SubCount:      int(subCount),
		DownloadCount: int(dlCount.SumValue),
	}
	return sgList
}

//字幕组主页影视详情
func SubGroupSeriesDetail(gtr commonModel.QuerySeriesSubGroupBody, e *commonConfig.ErrorStruct) commonModel.SeriesSubGroupList {
	tvList := []commonModel.SubLibraryBody{}
	movieList := []commonModel.SubLibraryBody{}
	limitTv := 12
	limitMovie := 6
	switch gtr.TvConditions {
	case 0:
		tvList = subGroupSeriesAll(gtr, true, 0, limitTv, e)
	case 1:
		tvList = subGroupSeriesUpdate(gtr, true, 0, limitTv, e)
	case 2:
		tvList = subGroupSeriesDownload(gtr, true, 0, limitTv, e)
	}
	switch gtr.MovieConditions {
	case 0:
		movieList = subGroupSeriesAll(gtr, false, 0, limitMovie, e)
	case 1:
		tvList = subGroupSeriesUpdate(gtr, false, 0, limitMovie, e)
	case 2:
		tvList = subGroupSeriesDownload(gtr, false, 0, limitMovie, e)
	}
	seriesList := commonModel.SeriesSubGroupList{
		TvList:    tvList,
		MovieList: movieList,
	}
	return seriesList

}

//字幕组影视详情电视剧列表
func SeriesListByTv(gtr commonModel.QuerySeriesSubGroupBody, tvOrMovie bool, e *commonConfig.ErrorStruct) []commonModel.SubLibraryBody {
	tvList := []commonModel.SubLibraryBody{}
	limit := 5
	ofset := (gtr.PageNumber - 1) * 5
	switch gtr.TvConditions {
	case 0:
		tvList = subGroupSeriesAll(gtr, tvOrMovie, ofset, limit, e)
	case 1:
		tvList = subGroupSeriesUpdate(gtr, tvOrMovie, ofset, limit, e)
	case 2:
		tvList = subGroupSeriesDownload(gtr, tvOrMovie, ofset, limit, e)
	}

	return tvList

}

//字幕组全部影片
func subGroupSeriesAll(gtr commonModel.QuerySeriesSubGroupBody, tvOrMove bool, offset int, limit int, e *commonConfig.ErrorStruct) []commonModel.SubLibraryBody {
	var series []commonModel.SubLibraryBody
	result1 := instance.ZimuzuDB.Table("series_model as s").
		Select("s.id as series_id,MAX(r.created_at) as ltime,cnname,origin_name,genres,series_type,screen_time,score,"+
			"s.desc,s.banner,count(r.id) as cnumber,cover,archive,tmdb_id,imdb_id,number_seasons").
		Joins("join resource_model as r on r.series_id=s.id").
		Joins("join user_model as u on r.user_id=u.id").
		Joins("join sub_group_model as sg on u.sub_group_id=sg.id").
		Where("sg.id=?", gtr.GroupId).
		Group("s.id").
		Order("s.created_at desc").Offset(offset)
	if tvOrMove {
		result1.Where("series_type = ?", 1).Limit(limit)
	} else {
		result1.Where("series_type = ?", 2).Limit(limit)
	}
	result1.Debug().Find(&series)
	for i, _ := range series {
		series[i].Cover = commonConfig.IMAGE_PATH + series[i].Cover
	}
	if result1.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result1.Error.Error())
		return nil
	}
	return series
}

//字幕组全部最近更新
func subGroupSeriesUpdate(gtr commonModel.QuerySeriesSubGroupBody, tvOrMove bool, offset int, limit int, e *commonConfig.ErrorStruct) []commonModel.SubLibraryBody {
	var series []commonModel.SubLibraryBody
	result1 := instance.ZimuzuDB.Table("series_model as s").
		Select("s.id as series_id,MAX(r.created_at) as ltime,cnname,origin_name,genres,series_type,screen_time,score,"+
			"s.desc,s.banner,count(r.id) as cnumber,cover,archive,tmdb_id,imdb_id,number_seasons").
		Joins("join resource_model as r on r.series_id=s.id").
		Joins("join user_model as u on r.user_id=u.id").
		Joins("join sub_group_model as sg on u.sub_group_id=sg.id").
		Where("sg.id=?", gtr.GroupId).
		Group("s.id").
		Order("ltime desc").Offset(offset)
	if tvOrMove {
		result1.Where("series_type = ?", 1).Limit(limit)
	} else {
		result1.Where("series_type = ?", 2).Limit(limit)
	}
	result1.Debug().Find(&series)
	for i, _ := range series {
		series[i].Cover = commonConfig.IMAGE_PATH + series[i].Cover
	}
	if result1.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result1.Error.Error())
		return nil
	}
	return series
}

//字幕组下载次数
func subGroupSeriesDownload(gtr commonModel.QuerySeriesSubGroupBody, tvOrMove bool, offset int, limit int, e *commonConfig.ErrorStruct) []commonModel.SubLibraryBody {
	var series []commonModel.SubLibraryBody
	result1 := instance.ZimuzuDB.Table("series_model as s").
		Select("s.id as series_id,MAX(r.created_at) as ltime,cnname,origin_name,genres,series_type,screen_time,score,"+
			"s.desc,s.banner,count(r.id) as cnumber,cover,archive,tmdb_id,imdb_id,number_seasons,sum(r.download_times) as dcount").
		Joins("join resource_model as r on r.series_id=s.id").
		Joins("join user_model as u on r.user_id=u.id").
		Joins("join sub_group_model as sg on u.sub_group_id=sg.id").
		Where("sg.id=?", gtr.GroupId).
		Group("s.id").
		Order("dcount desc").Offset(offset)
	if tvOrMove {
		result1.Where("series_type = ?", 1).Limit(limit)
	} else {
		result1.Where("series_type = ?", 2).Limit(limit)
	}
	result1.Debug().Find(&series)
	for i, _ := range series {
		series[i].Cover = commonConfig.IMAGE_PATH + series[i].Cover
	}
	if result1.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result1.Error.Error())
		return nil
	}
	return series
}

//字幕组列表
func SubGroupList(gtr commonModel.QuerySubGroupBody, e *commonConfig.ErrorStruct) []commonModel.SubGroupListBody {
	subGroupList := make([]commonModel.SubGroupListBody, 0, 20)
	subGroups := []commonModel.SubGroupModel{}
	result := instance.ZimuzuDB.Limit(commonConfig.Pageing_rows_subgroup).
		Offset((int(gtr.PageNumber) - 1) * commonConfig.Pageing_rows_subgroup).
		Find(&subGroups)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return nil
	}

	ChangeImage(subGroups)
	for _, sg := range subGroups {

		// result1 := instance.ZimuzuDB.Table("series_model as s").Select("distinct(s.id) series_id,cnname,cover,r.created_at").
		// 	Joins("join resource_model as r on s.id=r.series_id").
		// 	Joins("join user_model as u on r.user_id=u.id").
		// 	Where("u.sub_group_id=?", sg.ID).
		// 	Order("r.created_at desc")
		var series []commonModel.SeriesList

		result1 := instance.ZimuzuDB.Table("series_model as s").
			Select("s.id as series_id,MAX(r.created_at) as ltime,cnname,cover").
			Joins("join resource_model as r on r.series_id=s.id").
			Joins("join user_model as u on r.user_id=u.id").
			Joins("join sub_group_model as sg on u.sub_group_id=sg.id").
			Where("sg.id=?", sg.ID).
			Group("s.id").
			Order("ltime desc").
			Limit(5)
		result1.Debug().Find(&series)
		for i, _ := range series {
			series[i].Cover = commonConfig.IMAGE_PATH + series[i].Cover
		}

		// result2 := instance.ZimuzuDB.Table("(?) as v", result1).Select("series_id,cnname,cover").Group("series_id").
		// 	Find(&series)
		if result1.Error != nil {
			commonConfig.ErrorError(e, commonConfig.DB_ERROR, result1.Error.Error())
			return nil
		}
		var subCount int64
		var dlCount commonModel.SumInt
		instance.ZimuzuDB.Table("resource_model as r").
			Joins("join user_model as u on r.user_id=u.id").
			Joins("join sub_group_model as sg on u.sub_group_id=sg.id").
			Where("sg.id=?", sg.ID).
			Count(&subCount)

		instance.ZimuzuDB.Table("resource_model as r").
			Select("sum(r.downLoad_times) as sum_value").
			Joins("join user_model as u on r.user_id=u.id").
			Joins("join sub_group_model as sg on u.sub_group_id=sg.id").
			Where("sg.id=?", sg.ID).
			Take(&dlCount)

		sgList := commonModel.SubGroupListBody{
			SubGroups:     sg,
			Series:        series,
			SubCount:      int(subCount),
			DownloadCount: int(dlCount.SumValue),
		}
		subGroupList = append(subGroupList, sgList)
		// subGroupList = append(subGroupList, &sg)

	}
	return subGroupList
}

//创建字幕组
func CreateSubGroup(gtr commonModel.CreateSubGroupRequestBody, e *commonConfig.ErrorStruct) commonModel.SubGroupModel {
	var sg = commonModel.SubGroupModel{
		CreateSubGroupRequestBody: gtr,
	}

	result := instance.ZimuzuDB.Create(&sg)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return sg
	}
	return sg
}

func UpdateSubGroup(gtr *commonModel.SubGroupModel, e *commonConfig.ErrorStruct) {

	result := instance.ZimuzuDB.Save(gtr)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())

	}

}

func FindSubGroupAll(e *commonConfig.ErrorStruct) []commonModel.SubGroupModel {
	var subGroupModel = []commonModel.SubGroupModel{}
	instance.ZimuzuDB.Find(&subGroupModel)
	return subGroupModel
}

func FindSubGroupByName(name string) []commonModel.SubGroupModel {
	var subGroupModel = []commonModel.SubGroupModel{}
	instance.ZimuzuDB.Where("group_name = '%?%'", name).Find(&subGroupModel)
	return subGroupModel

}

//查找字幕组（按名称）
func HasSubGroup(name string) bool {
	var subGroupModel = commonModel.SubGroupModel{}
	instance.ZimuzuDB.Where("group_name = ?", name).Find(&subGroupModel)
	return subGroupModel.ID != 0

}
func HasSubGroupByID(id uint) (commonModel.SubGroupModel, bool) {
	var subGroupModel = commonModel.SubGroupModel{}
	instance.ZimuzuDB.Debug().Where("id = ?", id).Find(&subGroupModel)
	fmt.Println("*****************", subGroupModel)
	return subGroupModel, subGroupModel.ID == 0

}

//更改图片
func ChangeImage(seriesModels []commonModel.SubGroupModel) {
	for i, _ := range seriesModels {
		seriesModels[i].LogoPath = commonConfig.BASE_SUBICONPATH + seriesModels[i].LogoPath
		seriesModels[i].WechatPath = commonConfig.BASE_SUBICONPATH + seriesModels[i].WechatPath
	}

}
