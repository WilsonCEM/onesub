package commonConfig

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ErrorCode uint16

const (
	ERRORCODE_OK     ErrorCode = 0000
	PARAMETER_ERROR  ErrorCode = 1111 // 参数校验异常
	NOTFOUNT_ERROR   ErrorCode = 1002 // 数据缺失
	DB_ERROR         ErrorCode = 1003 // 新增数据异常
	CONDITION_ERROR  ErrorCode = 9999 // 业务异常
	THIRDPARTY_ERROR ErrorCode = 9998 // 三方异常
	SYSTEM_ERROR     ErrorCode = 9997 // 系统异常
	TOKEN_TIMEOUT    ErrorCode = 1112 // token过期
	TOKEN_INVALID    ErrorCode = 1113 // token无效
	AUTH_INVALID     ErrorCode = 1114 // 登陆后访问
	AUTH_NOTEGNORE   ErrorCode = 1115 // 权限不足
	FILEUPLOAD_ERROR ErrorCode = 1116 //上传文件失败
	EMAIL_SENDERROR  ErrorCode = 1117 //邮件发送失败
)

// 程序内的交互结构
type ErrorStruct struct {
	ErrorCode ErrorCode
	Message   string
}

type CustomContext struct {
	*gin.Context
	errorStruct ErrorStruct
}

const (
	ERRORMESSAGE_GOODSTEMPLATE_OK string = "success"
)

func HasError(e *ErrorStruct) bool {
	return e.ErrorCode != ERRORCODE_OK
}

// 注意 函数间交互统一用包装好的错误
func ErrorSuccess(e *ErrorStruct) {
	e = &ErrorStruct{
		ErrorCode: ERRORCODE_OK,
		Message:   ERRORMESSAGE_GOODSTEMPLATE_OK,
	}
}

// 注意 函数间交互统一用包装好的错误
func ErrorError(e *ErrorStruct, code ErrorCode, msg string) {
	e.ErrorCode = code
	e.Message = msg
}

func ErrorResponse(c *gin.Context, Error *ErrorStruct) {
	c.JSON(http.StatusOK, gin.H{
		"code": Error.ErrorCode,
		"msg":  Error.Message,
	})
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": ERRORCODE_OK,
		"data": data,
		"msg":  ERRORMESSAGE_GOODSTEMPLATE_OK,
	})
}
func SuccessResponseFile(c *gin.Context, data interface{}, filePath string, fileName string) {
	// c.JSON(http.StatusOK, gin.H{
	// 	"code": ERRORCODE_OK,
	// 	"data": data,
	// 	"msg":  ERRORMESSAGE_GOODSTEMPLATE_OK,
	// })
	fmt.Println("*****************", fileName)
	c.Header("Content-Type", "application/octet-stream")
	// c.Header("Content-Disposition", "attachment; filename="+fileName) // 用来指定下载下来的文件名
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(fileName))) // 用来指定下载下来的文件名

	c.Header("Content-Transfer-Encoding", "binary")
	c.File(filePath)
}

func HandleDBError(err error, e *ErrorStruct, msg string) {
	if err != nil {
		switch err.Error() {
		case gorm.ErrRecordNotFound.Error():
			ErrorError(e, NOTFOUNT_ERROR, msg)
			break
		default:
			ErrorError(e, SYSTEM_ERROR, "dberror")
			break
		}
	}
}
