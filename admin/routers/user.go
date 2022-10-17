package routers

import (
	"zimuzu/admin/controller"
	commonMiddleware "zimuzu/common/middleware"
	commonModel "zimuzu/common/models"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(router *gin.Engine) {
	apiUser := router.Group("/admin/user")
	/**
	 * @api { POST } /api/user/register 用户注册
	 * @apiGroup User
	 *
	 * @apiBody {String} username 邮箱
	 * @apiBody {String} email 邮箱
	 * @apiBody {String} password 密码
	 * @apiSampleRequest /admin/user/register
	 */
	apiUser.POST("/register", controller.CreateUser)
	/**
	 * @api { POST } /api/user/login 用户登陆
	 * @apiGroup User
	 *
	 * @apiBody {String} username 邮箱
	 * @apiBody {String} email 邮箱
	 * @apiBody {String} password 密码
	 * @apiBody {number=1,2} loginBy 登陆通过 1 用户名登陆 2 邮箱登陆
	 * @apiSampleRequest /admin/user/login
	 */
	apiUser.POST("/login", controller.Login)
	/**
	* 用户激活
	 */
	apiUser.GET("/activation", controller.UserActivation)

	/**
	 * @api { POST } admin/user/sendretrievepassword/ 发送找回密码邮件
	 * @apiGroup User
	 *
	 * @apiBody {String} username 邮箱
	 * @apiBody {String} email 邮箱
	 * @apiBody {String} password 密码
	 * @apiBody {number=1,2} loginBy 登陆通过 1 用户名 2 邮箱
	 * @apiSampleRequest admin/user/sendretrievepassword/
	 */
	apiUser.POST("/sendretrievepassword", controller.SendRetrievePassword)

	/*
	 * @api { GET } admin/user/retrievepassword/ 找回密码
	 * @apiGroup User
	 *
	 * @apiBody {String} token 邮箱
	 * @apiSampleRequest admin/user/retrievepassword/
	 */
	apiUser.GET("/retrievepassword", controller.RetrievePassword)

	/*
	 * @api { POST } admin/user/changepassword/ 修改密码
	 * @apiGroup User
	 *
	 * @apiBody {String} token 邮箱
	 * @apiSampleRequest admin/user/changepassword/
	 */
	apiUser.POST("/changepassword", commonMiddleware.AuthContext(commonModel.USER_ROLE_NOACTIVR), controller.ChangePassword)
	/*
	 * @api { GET } admin/user/newpassword/ 重置密码
	 * @apiGroup User
	 *
	 * @apiBody {String} token
	 * @apiSampleRequest admin/user/newpassword/
	 */
	apiUser.POST("/resetpassword", commonMiddleware.AuthContext(commonModel.USER_ROLE_NOACTIVR), controller.ResetPassword)

	/*
	 * @api { GET } admin/user/findusername/ 查找用户
	 * @apiGroup User
	 *
	 * @apiBody {String} text 用户名
	 * @apiSampleRequest admin/user/findusername/
	 */
	apiUser.POST("/findusername", commonMiddleware.AuthContext(commonModel.USER_ROLE_SUBGROUP), controller.FindUserByName)

	/*
	 * @api { GET } admin/user/usertosubgroup/ 查找用户
	 * @apiGroup User
	 *
	 * @apiBody {int} userId 用户id
	 * @apiBody {int} groupId 字幕组id
	 * @apiSampleRequest admin/user/usertosubgroup/
	 */
	apiUser.POST("/usertosubgroup", commonMiddleware.AuthContext(commonModel.USER_ROLE_ADMIN), controller.UserToSubGroup)

	/*
	 * @api { GET } admin/user/usertosubgroup/ 查找用户
	 * @apiGroup User
	 *
	 * @apiBody {String} text 用户名
	 * @apiSampleRequest admin/user/subgroupremoveuser/
	 */
	apiUser.POST("/subgroupremoveuser", commonMiddleware.AuthContext(commonModel.USER_ROLE_SUBGROUP), controller.SubGroupRemoveUser)
}
