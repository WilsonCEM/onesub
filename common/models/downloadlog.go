package commonModel

type DownloadLogModel struct {
	BaseModel
	ResourceId uint `gorm:"not null" json:"resourceId" binding:"required"`
	UserId     uint `grom:"default 0" json:"userId"`
	SeriesID   uint `gorm:"not null" json:"seriesId"`
}
