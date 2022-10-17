package commonModel

type SubGroupModel struct {
	BaseModel
	CreateSubGroupRequestBody
}

type CreateSubGroupRequestBody struct {
	GroupName  string      `json:"groupname" binding:"required"` //字幕组名称
	LogoPath   string      `json:"logo_path"`                    //logo路径
	Desc       string      `json:"desc"`                         //详细介绍
	Blog       string      `json:"blog"`                         //微博地址
	Website    string      `json:"website"`                      //官网
	WechatPath string      `json:"wechat_path"`                  //公众号二维码路径
	Users      []UserModel `gorm:"foreignKey:SubGroupId"`        //用户关联外键
}

type UpdateSubGroupRequestBody struct {
	GroupId    uint        `json:"groupId" binding:"required"`
	GroupName  string      `json:"groupname" binding:"required"` //字幕组名称
	LogoPath   string      `json:"logo_path"`                    //logo路径
	Desc       string      `json:"desc"`                         //详细介绍
	Blog       string      `json:"blog"`                         //微博地址
	Website    string      `json:"website"`                      //官网
	WechatPath string      `json:"wechat_path"`                  //公众号二维码路径
	Users      []UserModel `gorm:"foreignKey:SubGroupId"`        //用户关联外键
}

type SeriesList struct {
	SeriesId uint   `json:"seriesId"`
	Cnname   string `json:"cnname"`
	Cover    string `json:"cover"`
}

type SumInt struct {
	SumValue int64 `json:"sum_value"`
}

type QuerySubGroupBody struct {
	GroupId    uint `json:"groupId"`
	PageNumber int  `json:"pagenumber"`
}

type SubGroupListBody struct {
	SubGroups     SubGroupModel
	Series        []SeriesList
	SubCount      int `json:"subcount"`
	DownloadCount int `json:"downloadcount"`
}

type SeriesSubGroupList struct {
	TvList    []SubLibraryBody
	MovieList []SubLibraryBody
}

type SubGroupDetailBody struct {
	SubGroup      SubGroupModel
	SubCount      int `json:"subcount"`
	DownloadCount int `json:"downloadcount"`
}

type QuerySeriesSubGroupBody struct {
	GroupId         uint `json:"groupId"`
	TvConditions    int  `json:"tvconditions"`    //0 全部 1 最近更新 2最多下载
	MovieConditions int  `json:"movieconditions"` //0 全部 1 最近更新 2最多下载
	PageNumber      int  `json:"pagenumber"`
}

type SubGroupByRoel struct {
	SubGroup      SubGroupModel
	SubCount      int  `json:"subcount"`
	DownloadCount int  `json:"downloadcount"`
	Flag          bool `json:"flag"`
}
