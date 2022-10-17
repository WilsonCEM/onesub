package instance

import (
	commonModel "zimuzu/common/models"
	commonUtils "zimuzu/common/utils"

	"gorm.io/gorm"
)

var ZimuzuDB *gorm.DB

func InitDBInstance() {
	ZimuzuDB = commonUtils.ConnectDB()

	ZimuzuDB.AutoMigrate(&commonModel.UserModel{})
	ZimuzuDB.AutoMigrate(&commonModel.ResourceModel{})
	ZimuzuDB.AutoMigrate(&commonModel.FetchTMDBRecordModel{})
	ZimuzuDB.AutoMigrate(&commonModel.SeriesModel{})
	ZimuzuDB.AutoMigrate(&commonModel.SubGroupModel{})
	ZimuzuDB.AutoMigrate(&commonModel.DownloadLogModel{})
	ZimuzuDB.AutoMigrate(&commonModel.ActivationModel{})

}
