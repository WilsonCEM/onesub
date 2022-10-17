package service

import (
	"fmt"
	"strconv"
	"zimuzu/admin/instance"
	commonConfig "zimuzu/common/config"
	commonModel "zimuzu/common/models"

	"github.com/ryanbradynd05/go-tmdb"
)

var tmdbAPI *tmdb.TMDb

func TakeRecordByTMDBID(tmdbId string, tmdbRecord *commonModel.FetchTMDBRecord, e *commonConfig.ErrorStruct) {
	result := instance.ZimuzuDB.Take(&tmdbRecord, tmdbId)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
	}
}
func CreateBaseTMDB(tmdbId string, seriesType commonModel.SeriesType, e *commonConfig.ErrorStruct) {
	config := tmdb.Config{
		APIKey:   "2f192a534054783f241dd4ebc24da30c",
		Proxies:  nil,
		UseProxy: false,
	}

	tmdbAPI = tmdb.Init(config)
	m := make(map[string]string)
	m["language"] = "zh-CN"
	id, _ := strconv.Atoi(tmdbId)
	var move1 *tmdb.Movie
	var tv1 *tmdb.TV
	var err error
	var err2 error
	var j string
	if seriesType == commonModel.SERIES_TYPE_SEASON {
		tv1, err = tmdbAPI.GetTvInfo(id, m)
		if err != nil {
			commonConfig.ErrorError(e, commonConfig.CONDITION_ERROR, err.Error())
			// return record
		}
		j, err2 = tmdb.ToJSON(tv1)
	} else {
		move1, err = tmdbAPI.GetMovieInfo(id, m)
		if err != nil {
			commonConfig.ErrorError(e, commonConfig.CONDITION_ERROR, err.Error())
			// return record
		}
		j, err2 = tmdb.ToJSON(move1)
	}

	if err2 != nil {
		commonConfig.ErrorError(e, commonConfig.CONDITION_ERROR, err2.Error())
		// return record
	}
	fmt.Println(j)
	// return record
}

// func CreateBaseSeries(imdbid string) {

// 	opts := append(chromedp.DefaultExecAllocatorOptions[:],
// 		chromedp.NoDefaultBrowserCheck,                   //不检查默认浏览器
// 		chromedp.Flag("headless", false),                 //开启图像界面
// 		chromedp.Flag("ignore-certificate-errors", true), //忽略错误
// 		chromedp.Flag("disable-web-security", true),      //禁用网络安全标志
// 		chromedp.Flag("disable-extensions", false),       //开启插件支持
// 		chromedp.Flag("disable-default-apps", false),
// 		chromedp.NoFirstRun, //不是首次运行
// 		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36"),
// 	)

// 	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
// 	//	defer cancel()
// 	//标签1
// 	ctx, _ := chromedp.NewContext(
// 		allocCtx,
// 		chromedp.WithLogf(log.Printf),
// 	)

// 	// ch := addNewTabListener(ctx)
// 	// run task list
// 	err := chromedp.Run(ctx,
// 		chromedp.Navigate("https://buyin.jinritemai.com/dashboard/live/control"),
// 		// chromedp.Sleep(3*time.Second),
// 		chromedp.WaitReady(`._1M2DJMjmtJZhkk7ei5nC9F`, chromedp.ByQuery),
// 		// chromedp.Sleep(5*time.Second),
// 		chromedp.Click(`._1M2DJMjmtJZhkk7ei5nC9F`, chromedp.ByQuery, chromedp.NodeNotVisible),
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// }
