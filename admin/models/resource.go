package models

import (
	commonModel "zimuzu/common/models"
)

func InitResourceModel(body commonModel.CreateResourceRequestBody) commonModel.ResourceModel {

	return commonModel.ResourceModel{
		Resource: commonModel.Resource{
			CreateResourceRequestBody: body,
			DownloadTimes:             0,
		},
	}
}
