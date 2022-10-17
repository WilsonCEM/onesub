package controller

import (
	"fmt"
	"regexp"
	"strings"
	"zimuzu/admin/models"
	"zimuzu/admin/service"
	commonConfig "zimuzu/common/config"
	commonModel "zimuzu/common/models"
	commonUtils "zimuzu/common/utils"

	"github.com/gin-gonic/gin"
)

func CreateSeries(c *gin.Context) {
	var gtr commonModel.CreateSeriesRequestBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}

	service.CreateSeries(gtr, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, gtr)
}

//删除影片
func DeleteSeries(c *gin.Context) {
	var gtr models.FindIndexSeriesBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}

	service.DeleteSeriesService(gtr, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, gtr)
}

//查询影片信息
func FindSeriesDetail(c *gin.Context) {
	var gtr models.FindIndexSeriesBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}
	seriesModels := service.QuerySeriesByID(gtr, &e)

	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	service.ChangeImageOne(&seriesModels)
	commonConfig.SuccessResponse(c, seriesModels)
}

//读取首页推荐
func FindRecommend(c *gin.Context) {
	var e commonConfig.ErrorStruct
	seriesModels := service.FindrecommendSeries(&e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, seriesModels)
}

//热门下载
func HotSeries(c *gin.Context) {
	var gtr models.FindHotSeriesBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}
	seriesModels := service.FindHotSeries(&e, gtr)

	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, seriesModels)
}

//最新发布
func NewSeries(c *gin.Context) {
	var e commonConfig.ErrorStruct
	seriesModels := service.FindNewSeries(&e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, seriesModels)
}

//搜索框搜索
func IndexSeries(c *gin.Context) {
	var gtr models.FindIndexSeriesBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}
	seriesModels := service.FindIndexSeries(gtr, &e)

	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, seriesModels)
}

//搜索详情页字幕
func SearchResourceDetail(c *gin.Context) {
	var gtr models.FindIndexSeriesBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}
	seriesModels := service.FindSearchResourceDetail(gtr, &e)
	commonConfig.SuccessResponse(c, seriesModels)
}

// //查找
// func FindSeriesDetail(c *gin.Context) {
// 	var gtr models.FindIndexSeriesBody
// 	var e commonConfig.ErrorStruct
// 	if !commonUtils.Validator(c, &e, &gtr) {
// 		return
// 	}
// 	seriesModels := service.FindSearchResourceDetail(gtr, &e)
// 	if commonConfig.HasError(&e) {
// 		commonConfig.ErrorResponse(c, &e)
// 		return
// 	}
// 	commonConfig.SuccessResponse(c, seriesModels)
// }

//英美剧详情
func SeriesListByTV(c *gin.Context) {
	var gtr models.FindIndexSeriesBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}
	seriesModels := service.SeriesListByTVService(gtr, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, seriesModels)
}

//搜索详情页
func SearcSeriveshDetail(c *gin.Context) {
	var gtr models.FindIndexSeriesBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}
	seriesModels := service.FindSearchSeriesDetail(gtr)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, seriesModels)
}

//影片归档
func ArchiveSeries(c *gin.Context) {
	var gtr models.ArchiveSeriesRequestBody
	text := "影片已撤销归档"

	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}
	if gtr.Archive {
		text = "影片归档成功"
	}
	service.TakeSeries(gtr, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	commonConfig.SuccessResponse(c, text)
}

// 从TMDB读取数据
func FetchDataByTMDB(c *gin.Context) {
	var gtr models.FetchDataByTMDBRequestBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}

	id, seriesType := getMovieIdByTMDBURL(gtr.TMDB_URL, &e)

	if id == "" {

		commonConfig.ErrorError(&e, commonConfig.PARAMETER_ERROR, "URL不符合规则 需包含规则：themoviedb.org/movie/xxx")
		commonConfig.ErrorResponse(c, &e)
		return
	}

	var record commonModel.FetchTMDBRecordModel
	fmt.Println("__________________________________________", id)
	if service.CanFetchSeriesByTMDB(id, &record, &e) {
		seriesList := service.SubmitFetchSeriesByTMDB(id, gtr.TMDB_URL, true, record, seriesType, &e)
		if commonConfig.HasError(&e) {
			commonConfig.ErrorResponse(c, &e)
			return
		}
		commonConfig.SuccessResponse(c, seriesList)
	} else {
		commonConfig.ErrorResponse(c, &e)
	}
}

func UpdateDataByTMDB(c *gin.Context) {
	var gtr models.FetchDataByTMDBRequestBody
	var e commonConfig.ErrorStruct
	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}

	id, seriesType := getMovieIdByTMDBURL(gtr.TMDB_URL, &e)

	if id == "" {

		commonConfig.ErrorError(&e, commonConfig.PARAMETER_ERROR, "URL不符合规则 需包含规则：themoviedb.org/movie/xxx")
		commonConfig.ErrorResponse(c, &e)
		return
	}

	var record commonModel.FetchTMDBRecordModel
	fmt.Println("__________________________________________", id)
	if service.CanFetchSeriesByTMDB(id, &record, &e) {
		seriesList := service.SubmitFetchSeriesByTMDB(id, gtr.TMDB_URL, false, record, seriesType, &e)
		if commonConfig.HasError(&e) {
			commonConfig.ErrorResponse(c, &e)
			return
		}
		commonConfig.SuccessResponse(c, seriesList)
	} else {
		commonConfig.ErrorResponse(c, &e)
	}
}

// https://www.themoviedb.org/movie/606402 -> 606402
func getMovieIdByTMDBURL(url string, e *commonConfig.ErrorStruct) (string, commonModel.SeriesType) {
	reg := regexp.MustCompile(`themoviedb\.org\/movie\/(?P<movie_id>[0-9]*)`)

	surl := strings.Split(url, "/")
	seriesType := commonModel.SERIES_TYPE_MOVIE
	fmt.Println("***********************", len(surl))
	if len(surl) < 4 || surl[3] != "tv" && surl[3] != "movie" {
		// commonConfig.ErrorError(e, commonConfig.PARAMETER_ERROR, "URL不符合规则 需包含规则：themoviedb.org/movie/xxx")
		return "", 1
	}
	if surl[3] == "tv" {
		reg = regexp.MustCompile(`themoviedb\.org\/tv\/(?P<movie_id>[0-9]*)`)
		seriesType = commonModel.SERIES_TYPE_SEASON
	}

	if reg == nil {

		return "", seriesType

	}

	var res = reg.FindStringSubmatch(url)
	if len(res) > 0 {
		return res[1], seriesType
	}
	return "", seriesType
}

func DouBanCrawler(c *gin.Context) {
	var gtr commonModel.DouBanBody
	var e commonConfig.ErrorStruct

	if !commonUtils.Validator(c, &e, &gtr) {
		return
	}
	sid := models.FindIndexSeriesBody{SeriesId: gtr.SeriesId}
	sModel := service.QuerySeriesByID(sid, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}

	gtr = service.GetDouBan(gtr.DouBanPath, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	sModel.Director = gtr.Director
	sModel.Writers = gtr.Writers
	sModel.Actor = gtr.Actor
	sModel.CnOriginName = gtr.CnOriginName
	sModel.Genres = gtr.Genres
	sModel.DouBanPath = gtr.DouBanPath
	if gtr.Score != 0 {
		sModel.Score = gtr.Score
	}
	if sModel.IMDBId != "" {
		sModel.IMDBId = gtr.IMDBId
	}
	service.UpdateSeriesByDouBan(sModel, &e)
	if commonConfig.HasError(&e) {
		commonConfig.ErrorResponse(c, &e)
		return
	}
	gtr.SeriesId = sModel.ID
	commonConfig.SuccessResponse(c, gtr)
}
