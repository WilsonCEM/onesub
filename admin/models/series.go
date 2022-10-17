package models

import (
	commonModel "zimuzu/common/models"
)

type FetchDataByTMDBRequestBody struct {
	TMDB_URL string `json:"TMDB_URL" binding:"required"`
}

type ArchiveSeriesRequestBody struct {
	ID      uint `json:"seriesId"`
	Archive bool `json:"archive"`
}

type FindSeriesRequestBody struct {
	ID         uint                   `json:"id"`
	T          uint                   `json:"time"`
	SeriesType commonModel.SeriesType `json:"seriesType"`
	PageNumber uint                   `json:"pagenumber"`
}
type FindHotSeriesBody struct {
	T uint `json:"time"`
}

type ResponseSeriesDetail struct {
	SeriesDetail []commonModel.SeriesModel `json:"seriesDetail"`
	TVCount      int                       `json:"tvcount"`
	MovieCount   int                       `json:"moviecout"`
}

func InitSeriesModel(body commonModel.CreateSeriesRequestBody) commonModel.SeriesModel {
	return commonModel.SeriesModel{
		Series: commonModel.Series{
			CreateSeriesRequestBody: body,
			Views:                   0,
		},
	}
}

type FindIndexSeriesBody struct {
	SeriesId   uint   `json:"seriesid"`
	Text       string `json:"text"`
	PageNumber uint   `json:"pagenumber"`
	SeriesType uint   `json:"seriesType"`
}
