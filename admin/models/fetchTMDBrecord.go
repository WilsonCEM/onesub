package models

import (
	commonModel "zimuzu/common/models"
)

func InitTMDBRecordModel(tmdbid string) commonModel.FetchTMDBRecordModel {
	return commonModel.FetchTMDBRecordModel{
		FetchTMDBRecord: commonModel.FetchTMDBRecord{
			TMDBId:      tmdbid,
			FetchStatus: commonModel.FetchTMDBStatusProcess,
			ErrorReson:  "",
		},
	}
}
