package models

import (
	commonModel "zimuzu/common/models"
)

func InitUserModel(body commonModel.CreateUserRequestBody) commonModel.UserModel {
	return commonModel.UserModel{
		User: commonModel.User{
			CreateUserRequestBody: body,
			UserRole:              commonModel.USER_ROLE_NOACTIVR,
		},
	}
}

type MailBody struct {
	Title     string
	MailTo    string
	UserToken string
	Router    string
	Text      string
}
