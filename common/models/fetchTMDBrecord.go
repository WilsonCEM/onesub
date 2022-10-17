package commonModel

type FetchTMDBStatus uint8

const (
	FetchTMDBStatusProcess FetchTMDBStatus = 0 // 进行中
	FetchTMDBStatusSuccess FetchTMDBStatus = 1 // 成功
	FetchTMDBStatusFail    FetchTMDBStatus = 2 // 失败
)

type FetchTMDBRecordModel struct {
	BaseModel
	FetchTMDBRecord
}

type FetchTMDBRecord struct {
	TMDBId        string          `gorm:"not null=true;" json:"TMDB_ID"`  // TMDBID
	FetchStatus   FetchTMDBStatus `gorm:"default 0;" json:"fetchStatus"`  // 目前状态
	ErrorReson    string          `json:"errorReson"`                     // 错误原因
	NumberSeasons uint            `gorm:"not null" json:"number_seasons"` //当前季
}

type GenresBody struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
type SeasonsBody struct {
	Name         string `json:"name"`
	PosterPath   string `json:"poster_path"`
	SeasonNumber uint   `json:"season_number"`
	AirDate      string `json:"air_date"`
	Desc         string `json:"overview"`
	EpisodeCount uint   `json:"episode_count"`
}

type FetchTMDBResponse struct {
	TMDBId           uint          `json:"id"`                // TMDBID tv
	BackdropPath     string        `json:"backdrop_path"`     // 背景图片 tv
	ErrorReson       string        `json:"errorReson"`        // 错误原因
	Genres           []GenresBody  `json:"genres"`            //类型 tv
	OriginalLanguage string        `json:"original_language"` // 原始语言 tv
	OriginalTitle    string        `json:"original_title"`    // 原始标题
	Overview         string        `json:"overview"`          // 简介 tv
	ReleaseDate      string        `json:"release_date"`      // 发布日 tv
	Revenue          int           `json:"revenue"`           // 收入
	Status           string        `json:"status"`            // 放映状态
	Title            string        `json:"title"`             // 标题
	VoteAverage      float32       `json:"vote_average"`      // 评分 tv
	PosterPath       string        `json:"poster_path"`       // 封面 tv
	IMDBID           string        `json:"imdb_id"`
	Name             string        `json:"name"`              //剧集名 tv
	OriginalName     string        `json:"original_name"`     //剧集别名 tv
	NumberSeasons    uint          `json:"number_of_seasons"` //当前季数 tv
	Seasons          []SeasonsBody `json:"seasons"`
}
