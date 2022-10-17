package service

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"time"
	"zimuzu/admin/instance"
	"zimuzu/admin/models"
	commonConfig "zimuzu/common/config"
	commonModel "zimuzu/common/models"

	"gopkg.in/gomail.v2"
)

//注册用户
func CreateUser(gtr commonModel.CreateUserRequestBody, e *commonConfig.ErrorStruct) commonModel.UserModel {
	userModel := models.InitUserModel(gtr)
	result := instance.ZimuzuDB.Create(&userModel)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return userModel
	}
	return userModel
}

func FindUserByName(gtr commonModel.FindUerBody, user commonModel.UserModel, e *commonConfig.ErrorStruct) []commonModel.UserModel {
	userList := []commonModel.UserModel{}

	where := "username like '%" + gtr.UserName + "%' and user_role not in(0,1)"
	if user.UserRole > commonModel.USER_ROLE_SUBGROUP {
		return userList
	}

	if user.UserRole == commonModel.USER_ROLE_SUBGROUP {
		if gtr.SubGroupId == user.SubGroupId {
			where = where + " and sub_group_id=" + strconv.Itoa(int(gtr.SubGroupId))
		} else {
			return userList
		}
	}
	result := instance.ZimuzuDB.Where(where).Find(&userList)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
	}
	return userList
}

//添加字幕组成员
func UserToSubGroup(uts commonModel.UserToSubGroupBody, e *commonConfig.ErrorStruct) {
	result := instance.ZimuzuDB.Table("user_model").Where("id=? and  user_role not in(0,1)", uts.UserID).
		Updates(map[string]interface{}{"user_role": commonModel.USER_ROLE_SUBGROUP, "sub_group_id": uts.SubGroupId})
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
	}
}

func SubGroupRemoveUser(uts commonModel.UserToSubGroupBody, e *commonConfig.ErrorStruct) {
	result := instance.ZimuzuDB.Table("user_model").Where("id=? and  user_role not in(0,1)", uts.UserID).
		Updates(map[string]interface{}{"user_role": commonModel.USER_ROLE_NORMAL, "sub_group_id": 0})
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
	}
}

func SendEmail(mail models.MailBody) error {
	umail := gomail.NewMessage()
	// umail.From = commonConfig.EMAIL_USERNAME
	// umail.To = []string{userModel.Email}  "admin/user/activation/?t="
	password := commonConfig.EMAIL_CODE
	url := commonConfig.BASE_URL + mail.Router + mail.UserToken
	umailHTML := `<h1>` + mail.Text + `</h1><a href="` + url + `">请点击链接</a>`
	host := "smtp.qq.com"
	port := 465
	umail.SetHeader("From", commonConfig.EMAIL_USERNAME)
	umail.SetHeader("To", mail.MailTo) // 收件人，可以多个收件人，但必须使用相同的 SMTP 连接
	umail.SetHeader("Subject", mail.Title)
	umail.SetBody("text/html", umailHTML)
	d := gomail.NewDialer(
		host,
		port,
		commonConfig.EMAIL_USERNAME,
		password,
	)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := d.DialAndSend(umail)
	return err
}

//发送找回密码邮件
func SendRetrievePasswordService(gtr commonModel.UserLoginRequestBody, e *commonConfig.ErrorStruct) {
	mail := models.MailBody{
		Title:     "找回密码",
		MailTo:    gtr.Email,
		UserToken: gtr.UserToken,
		Router:    "admin/user/retrievepassword/?t=",
		Text:      "请点击下方链接修改新密码",
	}
	err := SendEmail(mail)
	if err != nil {
		fmt.Println("***************************", err.Error())
		commonConfig.ErrorError(e, commonConfig.EMAIL_SENDERROR, err.Error())
		return
	}
	currentTime := time.Now()
	m, _ := time.ParseDuration("1h")
	lm := currentTime.Add(m)
	ative := commonModel.Activation{UserId: gtr.ID, ActivToken: gtr.UserToken, FailureTime: lm, ActiveType: 1}
	ativeModel := &commonModel.ActivationModel{Activation: ative}
	result := instance.ZimuzuDB.Create(&ativeModel)
	if result.Error != nil {
		fmt.Println("***************", result.Error)
		commonConfig.ErrorError(e, commonConfig.EMAIL_SENDERROR, result.Error.Error())
		return
	}

}

//修改密码ChangePasswordService

func ChangePasswordService(gtr *commonModel.UserModel, e *commonConfig.ErrorStruct) {
	result := instance.ZimuzuDB.Model(&gtr).Update("password", gtr.Password)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return
	}
}

//找回密码RetrievePassword
func RetrievePasswordService(token string, id uint, e *commonConfig.ErrorStruct) {
	ativ := commonModel.Activation{UserId: id, ActivToken: token, ActiveType: 1}
	FindToken(ativ, e)
	if commonConfig.HasError(e) {
		commonConfig.ErrorError(e, commonConfig.EMAIL_SENDERROR, "该链接已失效！")
		return
	}

}

func ResetPasswordService(gtr *commonModel.UserModel, token string, e *commonConfig.ErrorStruct) {

	ativ := commonModel.Activation{UserId: gtr.ID, ActivToken: token, ActiveType: 1}
	var ativeModel = FindToken(ativ, e)
	if commonConfig.HasError(e) {
		commonConfig.ErrorError(e, commonConfig.EMAIL_SENDERROR, "该链接已失效！")
		return
	}
	result1 := instance.ZimuzuDB.Delete(&ativeModel)
	if result1.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result1.Error.Error())
		return
	}
	result := instance.ZimuzuDB.Model(&gtr).Update("password", gtr.Password)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return
	}

}

//发送验证邮件
func CreateEmail(token string, userModel commonModel.UserModel, e *commonConfig.ErrorStruct) bool {

	mail := models.MailBody{
		Title:     "用户激活",
		MailTo:    userModel.Email,
		UserToken: token,
		Router:    "admin/user/activation/?t=",
		Text:      "欢迎您加入一个字幕！请点击下方链接完成用户激活：",
	}
	err := SendEmail(mail)
	if err != nil {
		fmt.Println("***************************", err.Error())
		commonConfig.ErrorError(e, commonConfig.EMAIL_SENDERROR, err.Error())
		return false
	}
	currentTime := time.Now()
	m, _ := time.ParseDuration("1h")
	lm := currentTime.Add(m)
	ative := commonModel.Activation{UserId: userModel.ID, ActivToken: token, FailureTime: lm, ActiveType: 0}
	ativeModel := &commonModel.ActivationModel{Activation: ative}
	result := instance.ZimuzuDB.Create(&ativeModel)
	if result.Error != nil {
		fmt.Println("***************", result.Error)
		commonConfig.ErrorError(e, commonConfig.EMAIL_SENDERROR, result.Error.Error())
		return false
	}
	return true
}

//用户激活
func ActivationUser(token string, id uint, e *commonConfig.ErrorStruct) {
	ativ := commonModel.Activation{UserId: id, ActivToken: token, ActiveType: 0}
	var ativeModel = FindToken(ativ, e)
	if commonConfig.HasError(e) {
		commonConfig.ErrorError(e, commonConfig.EMAIL_SENDERROR, "该用户已激活！")
		return
	}
	result := instance.ZimuzuDB.Delete(&ativeModel)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return
	}
	result = instance.ZimuzuDB.Table("user_model").Where("id=?", ativeModel.UserId).Update("user_role", commonModel.USER_ROLE_NORMAL)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return
	}
}

//验证邮箱是否存在
func QueryUserByEmail(body *commonModel.UserLoginRequestBody, e *commonConfig.ErrorStruct) commonModel.UserModel {
	var user commonModel.UserModel
	result := instance.ZimuzuDB.Take(&user, "email = ?", body.Email)
	commonConfig.HandleDBError(result.Error, e, "用户不存在")
	return user
}

//验证用户名是否存在
func QueryUserByUserName(body *commonModel.UserLoginRequestBody, e *commonConfig.ErrorStruct) commonModel.UserModel {
	var user commonModel.UserModel
	result := instance.ZimuzuDB.Take(&user, "username = ?", body.Username)
	commonConfig.HandleDBError(result.Error, e, "用户不存在")
	return user
}

//判断用户是否存在
func QueryUserByID(body *commonModel.ChangePasswordBody, e *commonConfig.ErrorStruct) commonModel.UserModel {
	var user commonModel.UserModel
	result := instance.ZimuzuDB.Take(&user, "id = ?", body.UserID)
	commonConfig.HandleDBError(result.Error, e, "用户不存在")
	return user
}

//生成token
func JWTSign(userModel *commonModel.UserModel, e *commonConfig.ErrorStruct) string {
	token, err := commonModel.JWTSign(commonModel.JWTPayLoad{
		Uid:      userModel.ID,
		UserRole: userModel.UserRole,
		UserName: userModel.Username,
		GroupID:  userModel.SubGroupId,
	})
	if err != nil {
		commonConfig.ErrorError(e, commonConfig.SYSTEM_ERROR, "token sign error")
		return token
	}
	return token
}

func FindToken(activ commonModel.Activation, e *commonConfig.ErrorStruct) commonModel.ActivationModel {
	var activModel commonModel.ActivationModel
	result := instance.ZimuzuDB.Where("user_id=? and activ_token=? and active_type=?", activ.UserId, activ.ActivToken, activ.ActiveType).
		First(&activModel)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
		return activModel
	}
	return activModel
}

//根据token获取user
func FindUserByToken(token string, e *commonConfig.ErrorStruct) commonModel.UserModel {
	jwtModel, err := commonModel.JWTParse(token)
	var userModel commonModel.UserModel
	if err != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, err.Error())
		return userModel
	}
	result := instance.ZimuzuDB.Table("user_model").Where("id=?", jwtModel.JWTPayLoad.Uid).First(&userModel)
	if result.Error != nil {
		commonConfig.ErrorError(e, commonConfig.DB_ERROR, result.Error.Error())
	}
	return userModel

}
