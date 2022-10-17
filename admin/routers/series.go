package routers

import (
	"zimuzu/admin/controller"
	commonMiddleware "zimuzu/common/middleware"
	commonModel "zimuzu/common/models"

	"github.com/gin-gonic/gin"
)

func RegisterSeriesRouter(router *gin.Engine) {
	adminSeries := router.Group("/admin/series")
	/**
	* @api { POST } /admin/series/create 创建系列
	* @apiGroup Series
	*
	* @apiBody {String} cnname 中文名
	* @apiBody {String} [originname] 原名
	* @apiBody {String} [aliasname] 别名
	* @apiBody {String} [area] 国家
	* @apiBody {Number=1,2} seriesType 1 电影 2 剧集
	* @apiBody {String} [desc] 描述
	* @apiBody {String} [score] 评分
	* @apiBody {String} [banner] 横幅
	* @apiBody {String} cover 封面
	* @apiSampleRequest /admin/series/create
	 */
	adminSeries.POST("/create", commonMiddleware.AuthContext(commonModel.USER_ROLE_ADMIN), controller.CreateSeries)

	/**
	* @api {POST} /admin/series/recommend 首页轮播推荐
	* @apiGroup Series
	 */
	adminSeries.POST("/recommend", controller.FindRecommend)

	/**
	* @api {POST} /admin/series/hotseries 热门
	* @apiGroup Series
	*
	* @apiBody {number=0,1} [time] 0=24小时 1=7天
	 */
	adminSeries.POST("/hotseries", controller.HotSeries)

	/**
	* @api {POST} /admin/series/newseries 最新
	* @apiGroup Series
	 */
	adminSeries.POST("/newseries", controller.NewSeries)

	//adminSeries.POST("/sublibrary", controller.SubLibrary)

	/**
	* @api {POST} /admin/series/all 字幕库瀑布流查询
	* @apiGroup Series
	*
	* @api {number=ID} [id] 瀑布流最后一条记录的id
	* @api {number=0,1,2,3} [time] 0=24小时 1=其他 2=最新更新 3=最新开播
	* @api {number=0,1,2} [seriesType] 0=全部 1=英美剧 2=电影
	* @api {number} [pagenumber] 当前页
	 */
	//adminSeries.POST("/all", controller.FindSeries)
	/**
	* @api { POST } /admin/series/searchdetail 系列详情
	* @apiGroup Series
	*
	* @apiBody {String} id 按照id查询系列详情
	* @apiSampleRequest /admin/series/searchdetail
	 */
	adminSeries.POST("/searchdetail", controller.SearcSeriveshDetail)

	/**
	* @api { POST } /admin/series/searchresourcedetail 系列详情
	* @apiGroup Series
	*
	* @apiBody {String} id 按照id查询系列详情
	* @apiSampleRequest /admin/series/searchresourcedetail
	 */
	adminSeries.POST("/searchresourcedetail", controller.SearchResourceDetail)
	/**
	* @api { POST } /admin/series/archive 系列归档
	* @apiGroup Series
	*
	* @apiBody {String} id 按照id归档系列
	* @apiBody {String} archive 是否归档
	* @apiSampleRequest /admin/series/archive
	 */
	adminSeries.POST("/archive", commonMiddleware.AuthContext(commonModel.USER_ROLE_ADMIN), controller.ArchiveSeries)
	/**
	* @api { POST } /admin/series/bytmdb 从TMDB_API读取数据
	* @apiGroup Series
	*
	* @apiBody {String} TMDB_URL 按照TMDB的链接读取
	* @apiSampleRequest /admin/series/bytmdb
	 */
	adminSeries.POST("/bytmdb", commonMiddleware.AuthContext(commonModel.USER_ROLE_ADMIN), controller.FetchDataByTMDB)

	/**
	* @api { POST } /admin/series/bytmdb 爬取豆瓣信息
	* @apiGroup Series
	*
	* @apiBody {String} url 按照豆瓣的链接读取
	* @apiSampleRequest /admin/series/bytmdb
	 */
	adminSeries.POST("/doubancrawler", commonMiddleware.AuthContext(commonModel.USER_ROLE_ADMIN), controller.DouBanCrawler)

	/**
	* @api { POST } /admin/series/findindexseries 搜索框查询
	* @apiGroup Series
	*
	* @apiBody {String} text 模糊查询内容
	*
	 */
	adminSeries.POST("/findindexseries", controller.IndexSeries)

	/**
	* @api { POST } /admin/series/serieslistbytv 搜索框查询
	* @apiGroup Series
	*
	* @apiBody {String} seriesId 英美剧id
	*
	 */
	adminSeries.POST("/serieslistbytv", controller.SeriesListByTV)

	/**
	* @api { POST } /admin/series/deleteseries 删除影片
	* @apiGroup Series
	*
	* @apiBody {String} seriesId 电影id
	*
	 */
	adminSeries.POST("/deleteseries", commonMiddleware.AuthContext(commonModel.USER_ROLE_ADMIN), controller.DeleteSeries)

	/**
	* @api { POST } /admin/series/serieslistbytv 搜索框查询
	* @apiGroup Series
	*
	* @apiBody {String} seriesId 电影id
	*
	 */
	adminSeries.POST("/seriesdetail", controller.FindSeriesDetail)

	/**
	* @api { POST } /admin/series/bytmdb 从TMDB_API读取数据
	* @apiGroup Series
	*
	* @apiBody {String} TMDB_URL 按照TMDB的链接读取
	* @apiSampleRequest /admin/series/bytmdb
	 */
	adminSeries.POST("/updatebytmdb", commonMiddleware.AuthContext(commonModel.USER_ROLE_ADMIN), controller.UpdateDataByTMDB)

}
