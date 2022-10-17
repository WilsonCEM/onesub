package routers

import (
	"zimuzu/admin/controller"
	commonMiddleware "zimuzu/common/middleware"
	commonModel "zimuzu/common/models"

	"github.com/gin-gonic/gin"
)

func RegisterResourceRouter(router *gin.Engine) {
	adminSeries := router.Group("/admin/resource")
	/**
	* @api { POST } /admin/resource/create 上传字幕
	* @apiGroup Resource
	*
	* @apiBody {String} seriesId 关联系列
	* @apiBody {String} [resourceFile] 资源文件
	* @apiBody {String} language 语言
	* @apiBody {String} format 格式
	* @apiBody {String} translator 译者
	* @apiBody {String} origin 来源
	* @apiBody {String} seriesNo 系列集号
	* @apiBody {String} remarks 备注
	* @apiSampleRequest /admin/resource/create
	 */
	adminSeries.POST("/create", commonMiddleware.AuthContext(commonModel.USER_ROLE_NORMAL), controller.CreateResource)
	/**
	* @api { POST } /admin/resource/update 更新字幕
	* @apiGroup Resource
	*
	* @apiBody {String} id 字幕id
	* @apiBody {String} seriesId 关联系列
	* @apiBody {String} [resourceFile] 资源文件
	* @apiBody {String} language 语言
	* @apiBody {String} format 格式
	* @apiBody {String} translator 译者
	* @apiBody {String} origin 来源
	* @apiBody {String} seriesNo 系列集号
	* @apiBody {String} remarks 备注
	* @apiSampleRequest /admin/resource/create
	 */
	adminSeries.POST("/update", commonMiddleware.AuthContext(commonModel.USER_ROLE_NORMAL), controller.UpdateResource)

	/**
	* @api { POST } /admin/resource/delete 更新字幕
	* @apiGroup Resource
	*
	* @apiBody {uint} id 字幕id
	 */
	adminSeries.POST("/delete", commonMiddleware.AuthContext(commonModel.USER_ROLE_ADMIN), controller.DeleteResource)

	/**
	* @api{POST} /admin/resource/download 下载字幕
	* @apiGroup Resource
	*
	* @apiBody {String} seriesId 关联剧集
	* @apiBody {String} userId 下载用户 默认为0
	* @apiBody {String} ResourceId 下载字幕
	* @apiSampleRequest /admin/resource/download
	 */
	adminSeries.POST("/download", controller.DownloadResource)

	adminSeries.POST("/resourcedetail", controller.ResourceDetail)

	/**
	* @api{POST} /admin/resource/hotresource 字幕排行
	 */
	adminSeries.POST("/hotresource", controller.HotResSource)

	// adminSeries.POST("/uploadfile", controller.CreateUploadFile)

	/**
	* @api {POST} /admin/series/sublibrary 字幕库瀑布流查询
	* @apiGroup Series
	*
	* @api {number=ID} [id] 瀑布流最后一条记录的id
	* @api {number=0,1,2,3} [time] 0=24小时 1=其他 2=最新更新 3=最新开播
	* @api {number=0,1,2} [seriesType] 0=全部 1=英美剧 2=电影
	* @api {number} [pagenumber] 当前页
	 */
	adminSeries.POST("/sublibrary", controller.SubLibrary)

	/**
	* @api {POST} /admin//admin/resource/seriesresource 按series查询
	* @api {number} [seriesid] 电影/英美剧id
	 */
	adminSeries.POST("/seriesresource", controller.SeriesResource)

	/**
	* @api {POST} /admin//admin/resource/seriesresourcebytv 英美剧字幕分组查询
	* @api {number} [seriesid] 英美剧id
	 */
	adminSeries.POST("/seriesresourcebytv", controller.SeriesResourceByTV)

}
