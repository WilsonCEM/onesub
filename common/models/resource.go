package commonModel

type ResourceModel struct {
	BaseModel
	Resource
}

type CreateResourceRequestBody struct {
	SeriesId     uint   `gorm:"not null" json:"seriesId" binding:"required"` //关联影片id
	LoadFileName string `gorm:"not null" json:"loadfilename"`
	ResourceFile string `gorm:"not null" json:"resourceFile" ` //文件名
	SourceTitle  string `gorm:"not null" json:"sourcetitle"`   //片源
	Language     string `gorm:"not null" json:"language" binding:"required"`
	Format       string `gorm:"not null" json:"format" binding:"required"`     // 字幕格式
	Translator   string `gorm:"not null" json:"translator" binding:"required"` // 译者
	Origin       string `gorm:"not null" json:"origin" binding:"required"`     // 字幕来源
	SeriesNo     uint16 `gorm:"not null" json:"seriesNo" `                     // 系列中第几集
	Remarks      string `json:"remarks"`                                       // 备注
	UserId       uint   `gorm:"not null" json:"userid"`
}

type UpdateResourceRequestBody struct {
	ID           uint   `json:"id"`
	SeriesId     uint   `json:"seriesId" binding:"required"` //关联影片id
	LoadFileName string `json:"loadfilename"`
	ResourceFile string `json:"resourceFile" ` //文件名
	SourceTitle  string `json:"sourcetitle"`   //片源
	Language     string `json:"language" binding:"required"`
	Format       string `json:"format" binding:"required"`     // 字幕格式
	Translator   string `json:"translator" binding:"required"` // 译者
	Origin       string `json:"origin" binding:"required"`     // 字幕来源
	SeriesNo     uint16 `json:"seriesNo" `                     // 系列中第几集
	Remarks      string `json:"remarks"`                       // 备注
	UserId       uint   `json:"userid"`
}

type Resource struct {
	CreateResourceRequestBody
	DownloadTimes uint `gorm:"default=0" json:"downloadTimes"` // 下载次数

	User   UserModel
	Series SeriesModel `gorm:"foreignKey:ID;references:SeriesId"`
	// 	Cnname        string      `json:"cnname"`
	// 	GroupName     string      `json:"groupname"`
	// 	Ltime         time.Time   `json:time`
}

type FindResourceBody struct {
	BaseModel
	Resource
}

type FindSeriesResourceBody struct {
	SeriesId uint `json:"seriesid"`
}

type DeleteResourceBody struct {
	ID uint `json:"id"`
}

type ResponseSeriesResourceBody struct {
	ID            uint   `json:"id"`
	SeriesId      uint   `json:"seriesId"`
	ResourceFile  string `json:"resourcefile"`
	Language      string `json:"language"`
	Origin        string `json:"origin"`
	Translator    string `json:"translator"`
	DownloadTimes uint   `json:"downloadtimes"`
	UserId        uint   `json:"userid"`
	Username      string `json:"username"`
	SubGroupId    uint   `json:"subgroupid"`
	GroupName     string `json:"groupname"`
	SourceTitle   string `json:"sourcetitle"`
	SeriesNo      uint16 `json:"seriesNo" `
	Cnname        string `json:"cnname"`
	SeriesType    int    `json:"seriestype"`
	NumberSeasons int    `json:"number_seasons"`
	Format        string `json:"format" binding:"required"`
	Remarks       string `json:"remarks"` // 备注
}
