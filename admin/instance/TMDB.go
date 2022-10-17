package instance

import (
	commonUtils "zimuzu/common/utils"

	"github.com/go-resty/resty/v2"
)

var TMDB_Request *resty.Request

func InitTMDBInstance() {
	TMDB_Request = commonUtils.InitTMDBHTTPClient()
}
