package commonModel

import (
	"time"
)

type ActivationModel struct {
	BaseModel
	Activation
}

type Activation struct {
	UserId      uint      `gorm:"not null" json:"userid"`
	ActivToken  string    `gorm:"not null" json:"activtoken"`
	FailureTime time.Time `gorm:"failuretime"`
	ActiveType  uint      `gorm:"not null" json:"activetype"`
}
