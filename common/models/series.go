package commonModel

import "time"

type SeriesModel struct {
	BaseModel
	Series
}

type SeriesType uint8

const (
	SERIES_TYPE_SEASON SeriesType = 1 // 剧集
	SERIES_TYPE_MOVIE  SeriesType = 2 // 电影
)

type CreateSeriesRequestBody struct {
	Cnname        string     `json:"cnname" binding:"required"` //中文名
	OriginName    string     `json:"originname"`                //别名
	Director      string     `json:"director"`                  //导演
	Writers       string     `json:"writers"`                   //编剧
	Actor         string     `json:"actor"`                     //演员
	CnOriginName  string     `json:"cnoriginname"`              //演员
	Area          string     `json:"area"`                      //地区
	TMDBId        uint       `gorm:"not null" json:"tmdb_id"`
	Genres        string     `gorm:"not null" json:"genres"` //类型
	ScreenTime    string     `json:"screenTime"`             // 放映日
	SeriesType    SeriesType `gorm:"not null" json:"seriesType" binding:"required"`
	Desc          string     `json:"desc"`                                     //简介
	Score         float32    `json:"score"`                                    //评分
	Banner        string     `json:"banner"`                                   //竖版海报
	Cover         string     `gorm:"not null" json:"cover" binding:"required"` //横版海报
	IMDBId        string     `json:"imdbid"`
	NumberSeasons uint       `json:"numberseasons" gorm:"default 0"` //当前季数
	EpisodeCount  uint       `json:"episodecount"`                   //当季多少集
	TMDBPath      string     `json:"tmdb_path"`                      //TMDB目录
	DouBanPath    string     `json:"douban_path"`
}

type Series struct {
	CreateSeriesRequestBody
	Views     uint            `gorm:"default=0" json:"views"`
	Archive   bool            `gorm:"default=false" json:"archive" binding:"required"`
	Recommend bool            `gorm:"default=false" json:"recommend" binding:"required"` //是否推荐
	Resources []ResourceModel `gorm:"foreignKey:SeriesId"`
	HotValue  int             `json:"hotvalue"`
}

type PageingSerices struct {
	BaseModel
	Series
	PageNumber int `json:"pagenumber"`
}

type HotSeriesModel struct {
	BaseModel
	Banner string `json:"banner"`                    //竖版海报
	Cnname string `json:"cnname" binding:"required"` //中文名
}

type SubLibraryBody struct {
	SeriesId      uint       `json:"seriesId"`                      //影片id
	Ltime         time.Time  `json:"ltime"`                         //最后更新时间
	Cnname        string     `json:"cnname"`                        //中文名
	OriginName    string     `json:"originname"`                    //英文名or别名
	Cnumber       int        `json:"cnumber"`                       //字幕数量
	Cover         string     `json:"cover"`                         //海报
	GroupName     string     `json:"groupname"`                     //字幕组名称
	Archive       bool       `json:"archive"`                       //播出状态 false 未完结 true 已完结
	Desc          string     `json:"desc"`                          //简介
	Banner        string     `json:"banner"`                        //横版海报
	Genres        string     `json:"genres"`                        //类型
	NumberSeasons uint       `json:"numberseasons"`                 //当前季数
	SeriesType    SeriesType `json:"seriesType" binding:"required"` //影片类型
	ScreenTime    string     `json:"screenTime"`                    // 放映日
	Score         float32    `json:"score"`                         //评分
}
type QueryID struct {
	SeriesId uint `json:"seriesid"`
}

type DouBanBody struct {
	SeriesId     uint    `json:"seriesId"`
	Director     string  `json:"director"`     //导演
	Writers      string  `json:"writers"`      //编剧
	Actor        string  `json:"actor"`        //演员
	CnOriginName string  `json:"cnoriginname"` //演员
	Genres       string  `json:"genres"`       //类型
	DouBanPath   string  `json:"url"`          //豆瓣链接
	Score        float32 `json:"score"`        //评分
	IMDBId       string  `json:"imdbid"`
}
type Excessive struct {
	Title string
	Text  string
}
