package errors

import (
	"net/http"
)

type TestServerError struct {
	ErrCode  int    `json:"err_code"`
	HttpCode int    `json:"-"`
	ErrMsg   string `json:"err_msg"`
}

func NewTestServerError(errorCode int, httpCode int, errMsg string) *TestServerError {
	return &TestServerError{
		ErrCode:  errorCode,
		HttpCode: httpCode,
		ErrMsg:   errMsg,
	}
}

func (e *TestServerError) Error() string {
	return e.ErrMsg
}

var (
	UnsupportedHandler      = NewTestServerError(3, http.StatusNotFound, "该服务器的此方法未支持")
	UploadFormBindErr       = NewTestServerError(101, http.StatusBadRequest, "上传谱面表单格式错误")
	UploadChartErr          = NewTestServerError(102, http.StatusBadRequest, "传入谱面格式错误")
	BadUidErr               = NewTestServerError(103, http.StatusBadRequest, "查询到多个UID或未找到UID")
	ConvertorUnexpectedBeat = NewTestServerError(201, http.StatusInternalServerError, "意料之外的beat格式")
	FailCreateFile          = NewTestServerError(301, http.StatusInternalServerError, "创建谱面相关文件失败或保存失败")
	FileTooBig              = NewTestServerError(302, http.StatusBadRequest, "上传的音乐文件过大")
	BadBGMType              = NewTestServerError(303, http.StatusBadRequest, "上传BGM格式有误或校验BGM格式出错")
	BGMProcessError         = NewTestServerError(304, http.StatusInternalServerError, "BGM处理过程出错")
	ChartConvertFail        = NewTestServerError(305, http.StatusInternalServerError, "谱面转码出错")
	FailInsertDatabase      = NewTestServerError(306, http.StatusInternalServerError, "储存至数据库出错")
	FailUploadToTencentCos  = NewTestServerError(401, http.StatusInternalServerError, "上传到Tencent Cos出错")
)
