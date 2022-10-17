package commonModel

import (
	"time"
	commonUtils "zimuzu/common/utils"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type UserModel struct {
	BaseModel
	User
}

type UserRoleType uint8
type LoginMethods uint8

const (
	USER_ROLE_ROOT     UserRoleType = 0 // 超级管理员
	USER_ROLE_ADMIN    UserRoleType = 1 // 管理员
	USER_ROLE_SUBGROUP UserRoleType = 2 //字幕组成员
	USER_ROLE_NORMAL   UserRoleType = 3 // 注册用户
	USER_ROLE_NOACTIVR UserRoleType = 4 // 未激活用户

)

const (
	LOGIN_BY_USERNAME LoginMethods = 1
	LOGIN_BY_EMAIL    LoginMethods = 2
)

type UserLoginRequestBody struct {
	ID        uint
	LoginBy   LoginMethods `json:"loginBy" binding:"required"`
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	Password  string       `json:"password"`
	UserToken string
}

type CreateUserRequestBody struct {
	Username   string `gorm:"not null;unique" json:"username" binding:"required,min=2,max=12"`
	Email      string `gorm:"not null;unique" json:"email" binding:"required,email"`
	Password   string `gorm:"not null" json:"password" binding:"required,min=6,max=15"`
	SubGroupId uint   `gorm:"default 0"`
}

type ChangePasswordBody struct {
	UserID      uint   `json:"userid"`
	Password    string `json:"password"`
	NewPassword string `json:"newpassword"`
}

type UserToSubGroupBody struct {
	UserID     uint   `json:"userId" binding:"required"`
	SubGroupId uint   `json:"groupId" binding:"required"`
	UserName   string `json:"username"`
	UserRole   UserRoleType
}

type FindUerBody struct {
	UserID     uint   `json:"userId"`
	SubGroupId uint   `json:"groupId" binding:"required"`
	UserName   string `json:"username" binding:"required"`
	UserRole   UserRoleType
}

type User struct {
	CreateUserRequestBody
	UserRole  UserRoleType    `gorm:"not null;default=2" json:"userRole"`
	Resources []ResourceModel `gorm:"foreignKey:UserId"`
	SubGroups SubGroupModel   `gorm:"foreignKey:SubGroupId"`
}

func (u *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	hashPwd, e := commonUtils.HashAndSalt(u.Password)
	if e != nil {
		return e
	}
	u.Password = hashPwd
	return nil
}

type JWTPayLoad struct {
	Uid      uint
	UserRole UserRoleType
	UserName string
	GroupID  uint
}

type JWTModel struct {
	*jwt.StandardClaims
	TokenType string
	JWTPayLoad
}

// JWT签名相关
var (
	JWTPrivateKey   []byte        = []byte("whosyourdaddy")
	JWTValidityTime time.Duration = time.Hour * 24
)

func JWTSign(payload JWTPayLoad) (string, error) {
	t := jwt.New(jwt.SigningMethodHS256)

	t.Claims = &JWTModel{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(JWTValidityTime).Unix(),
		},
		"level1",
		payload,
	}
	jwt, err := t.SignedString(JWTPrivateKey)
	if err != nil {
		return jwt, err
	}
	return jwt, err
}

func JWTParse(token string) (JWTModel, error) {
	var jModel JWTModel
	tt, err := jwt.ParseWithClaims(token, &jModel, func(t *jwt.Token) (interface{}, error) {
		return JWTPrivateKey, nil
	})

	// 有效的
	if tt.Valid {
		return jModel, err
	}
	return jModel, err
}
