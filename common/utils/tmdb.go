package commonUtils

import (
	commonConfig "zimuzu/common/config"

	"github.com/go-resty/resty/v2"
)

func InitTMDBHTTPClient() *resty.Request {
	client := resty.New()
	client.SetHeader("Authorization", "Bearer "+commonConfig.TMDB_API_TOKEN)
	client.SetHeader("Content-Type", "application/json;charset=utf-8")
	return client.R()
}
