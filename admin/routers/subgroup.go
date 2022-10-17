package routers

import (
	"zimuzu/admin/controller"
	commonMiddleware "zimuzu/common/middleware"
	commonModel "zimuzu/common/models"

	"github.com/gin-gonic/gin"
)

func RegisterSubGroupRouter(router *gin.Engine) {
	apiSubGroup := router.Group("/admin/subgroup")

	/**
	* @api{POST} /admin/subgroup/hotsubgroup 字幕排行
	 */
	apiSubGroup.POST("/hotsubgroup", controller.HotSubGroup)

	/**
	* @api{POST} /admin/subgroup/create 创建字幕组
	 */
	apiSubGroup.POST("/create", commonMiddleware.AuthContext(commonModel.USER_ROLE_ADMIN), controller.CreateSubGroup)

	/**
	* @api{POST} /admin/subgroup/list 字幕组列表
	 */
	apiSubGroup.POST("/list", controller.SubgroupList)

	/**
	* @api{POST} /admin/subgroup/detail 字幕组详情
	 */
	apiSubGroup.POST("/detail", controller.SubGroupDetail)

	/**
	* @api{POST} /admin/subgroup/subgroupseries 字幕组影视详情
	 */
	apiSubGroup.POST("/subgroupseries", controller.SubGroupSeries)

	/**
	* @api{POST} /admin/subgroup/subgroupseriestV 字幕组影片详情（英美剧）
	 */
	apiSubGroup.POST("/subgroupseriestv", controller.SubGroupSeriesTV)
	/**
	* @api{POST} /admin/subgroup/subgroupseriesmovie 字幕组影片详情（电影）
	 */
	apiSubGroup.POST("/subgroupseriesmovie", controller.SubGroupSeriesMovie)
	/**
	* @api{POST} /admin/subgroup/update 修改字幕组
	 */
	apiSubGroup.POST("/update", commonMiddleware.AuthContext(commonModel.USER_ROLE_SUBGROUP), controller.SubgroupUpdate)
	/**
	 * @api {POST} /admin//admin/subgroup/findsubgroupall 查询所有字幕组
	 * @api {number} [seriesid] 英美剧id
	 */
	apiSubGroup.POST("/findsubgroupall", controller.FindSubGroupAll)

	// /**
	//  * @api {POST} /admin//admin/subgroup/findsubgroupall 查询所有字幕组
	//  * @api {number} [seriesid] 英美剧id
	//  */
	// apiSubGroup.POST("/findsubgroupall", controller.FindSubGroupAll)

}
