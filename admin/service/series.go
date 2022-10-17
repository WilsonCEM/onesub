package service

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"zimuzu/admin/instance"
	"zimuzu/admin/models"
	commonConfig "zimuzu/common/config"
	commonModel "zimuzu/common/models"
	commonUtils "zimuzu/common/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"gorm.io/gorm"
)

//创建剧集
func CreateSeries(gtr commonModel.CreateSeriesRequestBody, e *commonConfig.ErrorStruct) {
	seriesModel := models.InitSeriesModel(gtr)
	result := instance.ZimuzuDB.Save(&seriesModel)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return
	}

}

//删除剧集
func DeleteSeriesService(gtr models.FindIndexSeriesBody, e *commonConfig.ErrorStruct) {
	result := instance.ZimuzuDB.Begin()
	result.Where("resource_id=?", gtr.SeriesId).Delete(new(commonModel.DownloadLogModel))
	result.Where("series_id=?", gtr.SeriesId).Delete(new(commonModel.ResourceModel))
	result.Where("id=?", gtr.SeriesId).Delete(new(commonModel.SeriesModel))
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		result.Rollback()
	}
	result.Commit()
}

//豆瓣链接爬取
func GetDouBan(url string, e *commonConfig.ErrorStruct) commonModel.DouBanBody {
	var d commonModel.DouBanBody
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.Flag("hide-scrollbars", false),
		chromedp.Flag("mute-audio", false),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	//创建chrome窗口
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	var html1 string
	timeoutCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	if err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second),
		chromedp.WaitVisible(`#content > div.grid-16-8.clearfix > div.article > div.indent.clearfix > div.subjectwrap.clearfix`),
		chromedp.OuterHTML(`document.querySelector('div.subjectwrap.clearfix')`, &html1, chromedp.ByJSPath),
		//在这里加上你需要的后续操作，如Navigate，SendKeys，Click等
	); err != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, err.Error())
		return d
	}

	hl := ""
	doc, err1 := goquery.NewDocumentFromReader(strings.NewReader(html1))
	if err1 != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, err1.Error())
		return d
	}
	doc.Find("#info").Each(func(i int, selection *goquery.Selection) {
		hl = selection.Text()
	})
	list := strings.Split(hl, "\n")
	for i, _ := range list {
		if strings.TrimSpace(list[i]) != "" {
			sl := strings.Split(list[i], ":")

			sl[1] = strings.Replace(sl[1], "更多...", "", -1)
			title := strings.TrimSpace(sl[0])
			text := strings.TrimSpace(sl[1])
			// fmt.Println("---------------------", sl[0], sl[1])
			// data = append(data, commonModel.Excessive{Title: strings.TrimSpace(sl[0]), Text: strings.TrimSpace(sl[1])})
			switch title {
			case "导演":
				d.Director = text

			case "编剧":
				d.Writers = text

			case "主演":
				d.Actor = text

			case "类型":
				d.Genres = text

			case "又名":
				d.CnOriginName = text

			case "IMDb":
				d.IMDBId = text
			}
		}
	}

	doc.Find("#interest_sectl>div>div.rating_self>strong").
		Each(func(i int, selection *goquery.Selection) {
			// data = append(data, Data{Pl: selection.Text()})

			s, err := strconv.ParseFloat(selection.Text(), 32)
			if err != nil {
				d.Score = 0
			} else {
				d.Score = float32(s)
			}
		})
	d.DouBanPath = url
	return d
}

//查找影片
func QuerySeriesByID(gtr models.FindIndexSeriesBody, e *commonConfig.ErrorStruct) commonModel.SeriesModel {
	var seriesModel commonModel.SeriesModel
	result := instance.ZimuzuDB.Where("id=?", gtr.SeriesId).First(&seriesModel)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())

	}

	return seriesModel

}

//英美剧详情查询
func SeriesListByTVService(gtr models.FindIndexSeriesBody, e *commonConfig.ErrorStruct) []commonModel.SeriesModel {
	var seriesList []commonModel.SeriesModel
	seriesModel := QuerySeriesByID(gtr, e)
	result := instance.ZimuzuDB.Where("tmdb_id=?", seriesModel.TMDBId).Find(&seriesList)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
	}
	changeImage(seriesList)
	return seriesList
}

//验证剧集是否存在
func HasSeries(seriesId uint) bool {
	var seriesModel = commonModel.SeriesModel{}
	result := instance.ZimuzuDB.First(&seriesModel, seriesId)
	if result.Error != nil {
		return false
	}
	return true
}

//搜索页详情
func FindSearchResourceDetail(seriesBody models.FindIndexSeriesBody, e *commonConfig.ErrorStruct) []commonModel.ResourceModel {
	var resourceModels = []commonModel.ResourceModel{}
	findResource(seriesBody, &resourceModels, e)
	return resourceModels
}

func FindSearchSeriesDetail(seriesBody models.FindIndexSeriesBody) models.ResponseSeriesDetail {
	var seriesModel = []commonModel.SeriesModel{}
	var tv int
	var movie int

	tv, movie = takeSeries(seriesBody, &seriesModel)

	return models.ResponseSeriesDetail{
		SeriesDetail: seriesModel,
		TVCount:      tv,
		MovieCount:   movie,
	}
}

func takeSeries(gtr models.FindIndexSeriesBody, seriesModel *[]commonModel.SeriesModel) (int, int) {
	gtr.Text = "%" + gtr.Text + "%"
	var tv int64
	var movies int64
	result := instance.ZimuzuDB.
		Order("created_at desc").
		Where("cnname like ? or origin_name like ? or tmdb_id like ? or imdb_id like ?", gtr.Text, gtr.Text, gtr.Text, gtr.Text).
		Limit(20)

	result1 := result
	result2 := result
	result1.Debug().Table("series_model").Where("series_type =1").Count(&tv)
	result2.Debug().Table("series_model").Where("series_type =2").Count(&movies)
	if gtr.SeriesType != 0 {
		result.Where("series_type =?", gtr.SeriesType)
	}
	result.Debug().Find(&seriesModel)
	changeImage(*seriesModel)
	return int(tv), int(movies)
}

func findResource(gtr models.FindIndexSeriesBody, resourceModels *[]commonModel.ResourceModel, e *commonConfig.ErrorStruct) {
	gtr.Text = "%" + gtr.Text + "%"
	result := instance.ZimuzuDB.Table("series_model").Select("id").
		Order("created_at desc").
		Where("cnname like ? or origin_name like ? or tmdb_id like ? or imdb_id like ?", gtr.Text, gtr.Text, gtr.Text, gtr.Text)
	if gtr.SeriesType != 0 {
		result.Where("series_type =?", gtr.SeriesType)
	}
	result1 := instance.ZimuzuDB.Where("series_id in (?)", result).
		Preload("User").
		Preload("User.SubGroups").
		Order("updated_at desc").
		Limit(commonConfig.Pageing_rows).
		Offset((int(gtr.PageNumber) - 1) * commonConfig.Pageing_rows).
		Debug().
		Find(&resourceModels)
	if result1.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return
	}
	changeSubGroupIcon(*resourceModels)

}

func changeSubGroupIcon(subModels []commonModel.ResourceModel) {
	for i, _ := range subModels {
		subModels[i].User.SubGroups.LogoPath = commonConfig.BASE_SUBICONPATH + subModels[i].User.SubGroups.LogoPath

	}
}

//首页推荐轮播
func FindrecommendSeries(e *commonConfig.ErrorStruct) []commonModel.SeriesModel {
	seriesModels := []commonModel.SeriesModel{}

	result := instance.ZimuzuDB.Debug().Where("recommend=?", 1).Limit(5).Order("updated_at desc").Find(&seriesModels)

	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return seriesModels
	}
	changeImage(seriesModels)
	return seriesModels
}

//搜索框
func FindIndexSeries(gtr models.FindIndexSeriesBody, e *commonConfig.ErrorStruct) []commonModel.SubLibraryBody {
	resourceModels := []commonModel.SubLibraryBody{}
	gtr.Text = "%" + gtr.Text + "%"
	result := instance.ZimuzuDB.Table("series_model").
		Select("series_model.id as series_id,MAX(resource_model.created_at) as ltime,cnname,origin_name,genres,series_type,screen_time,score,"+
			"series_model.desc,series_model.banner,count(resource_model.id) as cnumber,cover,archive,tmdb_id,imdb_id,number_seasons").
		Joins("left join resource_model on resource_model.series_id=series_model.id").
		Group("series_model.id").
		Order("ltime desc").
		Limit(commonConfig.Pageing_rows).
		Where("cnname like ? or origin_name like ? or tmdb_id like ? or imdb_id like ?", gtr.Text, gtr.Text, gtr.Text, gtr.Text)
	result.Debug().Find(&resourceModels)
	ChangeSubImage(resourceModels, e)

	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return resourceModels
	}
	return resourceModels
}

func TakeSeries(as models.ArchiveSeriesRequestBody, e *commonConfig.ErrorStruct) {
	result := instance.ZimuzuDB.Table("series_model").Where("id=?", as.ID).Update("archive", as.Archive)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return
	}
}

//热门下载
func FindHotSeries(e *commonConfig.ErrorStruct, t models.FindHotSeriesBody) []commonModel.SeriesModel {
	seriesModels := []commonModel.SeriesModel{}
	var result *gorm.DB
	now := time.Now()
	old := now.AddDate(0, 0, -1)
	old2 := now.AddDate(0, 0, -7)
	fmt.Println("************", old, old2)
	// var idList commonModel.QueryID

	if t.T == 0 {
		result1 := instance.ZimuzuDB.Table("download_log_model").
			Select("distinct(series_id) as sid,max(created_at) as dtime,count(*) as dcount").
			Group("sid")
		result2 := instance.ZimuzuDB.Table("(?) as t", result1).Where("dtime>=?", old)
		result = instance.ZimuzuDB.Debug().Joins("left join (?) as l on id=l.sid", result2).
			Limit(20).
			Order("dcount desc,updated_at desc").
			Find(&seriesModels)
	} else {
		result1 := instance.ZimuzuDB.Table("download_log_model").
			Select("distinct(series_id) as sid,max(created_at) as dtime,count(*) as dcount").
			Group("sid")
		result2 := instance.ZimuzuDB.Table("(?) as t", result1).Where("dtime>=?", old2)
		result = instance.ZimuzuDB.Debug().Joins("left join (?) as l on id=l.sid", result2).
			Limit(20).
			Order("dcount desc,updated_at desc").
			Find(&seriesModels)
	}
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return seriesModels
	}
	changeImage(seriesModels)
	return seriesModels
}

//最新发布
func FindNewSeries(e *commonConfig.ErrorStruct) []commonModel.SeriesModel {
	seriesModels := []commonModel.SeriesModel{}
	result := instance.ZimuzuDB.Limit(20).Order("created_at desc").Find(&seriesModels)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return seriesModels
	}
	changeImage(seriesModels)
	return seriesModels
}

//建立tmdb信息
func CanFetchSeriesByTMDB(tmdbId string, tmdbRecord *commonModel.FetchTMDBRecordModel, e *commonConfig.ErrorStruct) bool {
	fmt.Println("__________________________________________")
	tmdbRecord.TMDBId = tmdbId
	result := instance.ZimuzuDB.Where("tmdb_id = ?", tmdbId).First(&tmdbRecord)
	fmt.Println("+++++++++++++++++", tmdbRecord.TMDBId, result.Error)
	if result.Error != nil {
		// 记录不存在 可以请求
		if result.Error == gorm.ErrRecordNotFound {
			return true
		}
		// DB异常
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return false
	}
	if tmdbRecord.FetchStatus == commonModel.FetchTMDBStatusProcess {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, "拉取进行中")
		return false
	}
	// 记录存在并且拉取状态为成功和失败
	return true
}

func HasSeriesTV(tmdbId uint, numberSeasons uint) bool {
	var seriesModel = commonModel.SeriesModel{}
	var e commonConfig.ErrorStruct
	result := instance.ZimuzuDB.Table("series_model").Where("tmdb_id=? and number_seasons=?", tmdbId, numberSeasons).First(&seriesModel)
	if result.Error != nil {
		commonConfig.ErrorError(&e, commonConfig.CONDITION_ERROR, result.Error.Error())
		return false
	}
	return true
}

//解析tmdb
func FetchResourceByTMDB(id string, url string, e *commonConfig.ErrorStruct, seriesType commonModel.SeriesType) commonModel.FetchTMDBResponse {
	var path = ""
	var result commonModel.FetchTMDBResponse
	if seriesType == commonModel.SERIES_TYPE_SEASON {
		path = commonConfig.TMDB_BASEURL + "/3/tv/" + id + "?api_key=" + commonConfig.TMDB_API_KEY + "&language=zh-CN"
		// res, err := instance.TMDB_Request.SetResult(&result).ForceContentType("application/json").Get(path)
		// fmt.Println(res, err)
		// fmt.Println("++++", result, path)
	} else {
		path = commonConfig.TMDB_BASEURL + "/3/movie/" + id + "?api_key=" + commonConfig.TMDB_API_KEY + "&language=zh-CN"

	}
	res, err := instance.TMDB_Request.SetResult(&result).ForceContentType("application/json").Get(path)
	fmt.Println(res, err)
	fmt.Println("++++", result, path)
	CreateBaseTMDB(id, seriesType, e)
	if err != nil {
		commonConfig.ErrorError(e, commonConfig.CONDITION_ERROR, err.Error())
		return result
	}
	return result
}

// 提交系列需要
// 1. 写入请求记录
// 2. 写入系列记录
// 3. 更新请求记录
func SubmitFetchSeriesByTMDB(tmdbId string, url string, cOru bool, record commonModel.FetchTMDBRecordModel, seriesType commonModel.SeriesType, e *commonConfig.ErrorStruct) []commonModel.SeriesModel {
	var wg sync.WaitGroup
	var tmdbRespose commonModel.FetchTMDBResponse

	wg.Add(2)
	// 请求系列
	go (func() {
		defer wg.Done()
		tmdbRespose = FetchResourceByTMDB(tmdbId, url, e, seriesType)
	})()
	// 处理记录
	go (func() {
		defer wg.Done()
		if record.ID != 0 {
			record.FetchStatus = commonModel.FetchTMDBStatusProcess
			instance.ZimuzuDB.Model(&record).Update("fetch_status", commonModel.FetchTMDBStatusProcess)
		} else {
			record = models.InitTMDBRecordModel(tmdbId)
			result := instance.ZimuzuDB.Create(&record)
			if result.Error != nil {
				commonConfig.ErrorError(e, commonConfig.CONDITION_ERROR, result.Error.Error())
				return
			}
		}
	})()
	wg.Wait()

	if commonConfig.HasError(e) {
		fmt.Printf("****************************")
		fmt.Println(e.Message)
		// instance.ZimuzuDB.Model(record).Update("fetch_status", commonModel.FetchTMDBStatusFail)
		return nil
	}

	var gtr = commonModel.CreateSeriesRequestBody{
		Banner:     tmdbRespose.BackdropPath,
		Cnname:     tmdbRespose.Title,
		OriginName: tmdbRespose.OriginalTitle,
		ScreenTime: tmdbRespose.ReleaseDate,
		TMDBId:     tmdbRespose.TMDBId,
		Desc:       tmdbRespose.Overview,
		Score:      tmdbRespose.VoteAverage,
		Cover:      tmdbRespose.PosterPath,
		IMDBId:     tmdbRespose.IMDBID,
		SeriesType: seriesType,
		TMDBPath:   url,
	}
	if tmdbRespose.Genres != nil {
		genres := ""
		for i, _ := range tmdbRespose.Genres {
			genres = genres + tmdbRespose.Genres[i].Name + ","
		}
		gtr.Genres = strings.TrimRight(genres, ",")
	}

	if seriesType == commonModel.SERIES_TYPE_SEASON {
		if cOru {

			gtr.Cnname = tmdbRespose.Name
			gtr.OriginName = tmdbRespose.OriginalName
			commonUtils.DownImg(tmdbRespose.BackdropPath)

			for i, _ := range tmdbRespose.Seasons {
				gtr.NumberSeasons = tmdbRespose.Seasons[i].SeasonNumber
				gtr.Cover = tmdbRespose.Seasons[i].PosterPath
				gtr.ScreenTime = tmdbRespose.Seasons[i].AirDate
				gtr.EpisodeCount = tmdbRespose.Seasons[i].EpisodeCount
				if tmdbRespose.Seasons[i].Desc != "" {
					gtr.Desc = tmdbRespose.Seasons[i].Desc
				}
				// if HasSeriesTV(gtr.TMDBId, gtr.NumberSeasons) {
				// 	continue
				// }
				var s commonModel.SeriesModel
				restmdbid := instance.ZimuzuDB.Debug().Where("tmdb_id=? and number_seasons=?", gtr.TMDBId, gtr.NumberSeasons).First(&s)
				if restmdbid.Error == nil {
					commonConfig.ErrorError(e, commonConfig.DB_ERROR, "该记录已存在！！")
					instance.ZimuzuDB.Model(record).Update("fetch_status", commonModel.FetchTMDBStatusFail)
					return nil
				}
				CreateSeries(gtr, e)
				if commonConfig.HasError(e) {
					instance.ZimuzuDB.Model(record).Update("fetch_status", commonModel.FetchTMDBStatusFail)
				} else {
					instance.ZimuzuDB.Model(record).Update("fetch_status", commonModel.FetchTMDBStatusSuccess)
				}

				commonUtils.DownImg(tmdbRespose.Seasons[i].PosterPath)
			}
		} else {
			gtr.Cnname = tmdbRespose.Name
			gtr.OriginName = tmdbRespose.OriginalName
			commonUtils.DownImg(tmdbRespose.BackdropPath)
			var series []commonModel.SeriesModel
			res1 := instance.ZimuzuDB.Where("tmdb_id=? and series_type=?", gtr.TMDBId, gtr.SeriesType).Find(&series)

			if res1.Error != nil {
				commonConfig.ErrorError(e, commonConfig.DB_ERROR, res1.Error.Error())
				return nil
			}
			os.Remove("./static/" + series[0].Banner)
			// var series commonModel.SeriesModel
			// res1 = instance.ZimuzuDB.Where("tmdb_id=? and series_type=?", gtr.TMDBId, gtr.SeriesType).Find(&series)
			for i, _ := range tmdbRespose.Seasons {
				gtr.NumberSeasons = tmdbRespose.Seasons[i].SeasonNumber
				gtr.Cover = tmdbRespose.Seasons[i].PosterPath
				gtr.ScreenTime = tmdbRespose.Seasons[i].AirDate
				gtr.EpisodeCount = tmdbRespose.Seasons[i].EpisodeCount
				if tmdbRespose.Seasons[i].Desc != "" {
					gtr.Desc = tmdbRespose.Seasons[i].Desc
				}
				// if HasSeriesTV(gtr.TMDBId, gtr.NumberSeasons) {
				// 	continue
				// }

				//res1 := instance.ZimuzuDB.Where("tmdb_id=? and series_type=?", gtr.TMDBId, gtr.SeriesType).First(&series)
				if res1.Error != nil {
					commonConfig.ErrorError(e, commonConfig.DB_ERROR, res1.Error.Error())
					return nil
				}

				if res1.Error != nil {
					commonConfig.ErrorError(e, commonConfig.DB_ERROR, res1.Error.Error())
					return nil
				}
				fmt.Println("*****************", gtr)
				UpdateSeries(gtr, series[i], e)

				if commonConfig.HasError(e) {
					commonConfig.ErrorError(e, commonConfig.DB_ERROR, "更新失败！！！")
					instance.ZimuzuDB.Model(record).Update("fetch_status", commonModel.FetchTMDBStatusFail)
					return nil
				} else {
					instance.ZimuzuDB.Model(record).Update("fetch_status", commonModel.FetchTMDBStatusSuccess)
				}
				os.Remove("./static/" + series[i].Cover)
				commonUtils.DownImg(tmdbRespose.Seasons[i].PosterPath)
			}
		}

	} else {
		if cOru {
			var s commonModel.SeriesModel
			restmdbid := instance.ZimuzuDB.Where("tmdb_id=?", gtr.TMDBId).First(&s)
			if restmdbid.Error == nil {
				commonConfig.ErrorError(e, commonConfig.DB_ERROR, "该记录已存在！！")
				instance.ZimuzuDB.Model(record).Update("fetch_status", commonModel.FetchTMDBStatusFail)
				return nil
			}
			commonUtils.DownImg(tmdbRespose.BackdropPath)
			commonUtils.DownImg(tmdbRespose.PosterPath)
			CreateSeries(gtr, e)
			if commonConfig.HasError(e) {
				instance.ZimuzuDB.Model(record).Update("fetch_status", commonModel.FetchTMDBStatusFail)
			} else {
				instance.ZimuzuDB.Model(record).Update("fetch_status", commonModel.FetchTMDBStatusSuccess)
			}
		} else {
			var series commonModel.SeriesModel
			res1 := instance.ZimuzuDB.Where("tmdb_id=? and series_type=?", gtr.TMDBId, gtr.SeriesType).First(&series)
			if res1.Error != nil {
				commonConfig.ErrorError(e, commonConfig.DB_ERROR, res1.Error.Error())
				return nil
			}

			commonUtils.DownImg(tmdbRespose.BackdropPath)
			commonUtils.DownImg(tmdbRespose.PosterPath)
			os.ReadDir("./static/" + series.Banner)
			os.ReadDir("./static/" + series.Cover)
			UpdateSeries(gtr, series, e)
			if commonConfig.HasError(e) {
				commonConfig.ErrorError(e, commonConfig.DB_ERROR, "更新失败！！！")
				instance.ZimuzuDB.Model(record).Update("fetch_status", commonModel.FetchTMDBStatusFail)
				return nil
			} else {
				instance.ZimuzuDB.Model(record).Update("fetch_status", commonModel.FetchTMDBStatusSuccess)
			}
		}

	}
	var seriesList []commonModel.SeriesModel
	result := instance.ZimuzuDB.Where("tmdb_id=?", gtr.TMDBId).Find(&seriesList)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.CONDITION_ERROR, result.Error.Error())
		return nil
	}
	changeImage(seriesList)
	return seriesList

}

//更新剧集
func UpdateSeries(gtr commonModel.CreateSeriesRequestBody, s commonModel.SeriesModel, e *commonConfig.ErrorStruct) {
	seriesModel := models.InitSeriesModel(gtr)
	seriesModel.CreatedAt = s.CreatedAt
	seriesModel.ID = s.ID
	result := instance.ZimuzuDB.Debug().Save(&seriesModel)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return
	}

}

func UpdateSeriesByDouBan(s commonModel.SeriesModel, e *commonConfig.ErrorStruct) {

	result := instance.ZimuzuDB.Debug().Save(&s)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return
	}
}

//完善图片地址
func changeImage(seriesModels []commonModel.SeriesModel) {
	for i, _ := range seriesModels {
		seriesModels[i].Banner = commonConfig.IMAGE_PATH + seriesModels[i].Banner
		seriesModels[i].Cover = commonConfig.IMAGE_PATH + seriesModels[i].Cover
		if i != 0 {
			seriesModels[i].HotValue = seriesModels[i-1].HotValue - 1543
		} else {
			seriesModels[i].HotValue = 8575066
		}
	}
}

func ChangeImageOne(seriesModels *commonModel.SeriesModel) {

	seriesModels.Banner = commonConfig.IMAGE_PATH + seriesModels.Banner
	seriesModels.Cover = commonConfig.IMAGE_PATH + seriesModels.Cover

}
