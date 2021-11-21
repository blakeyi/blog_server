package api

import (
	"net/http"
)

var Router = map[string]func(w http.ResponseWriter, req *http.Request){
	"/articleList":   ArticleList,
	"/articleQuery": ArticleQuery,
	"/articleCreate": ArticleCreate,
	"/articleUpdate": ArticleUpdate,
	"/articleDelete": ArticleDelete,

	"/userLogin": UserLogin,
}

type Response struct {
	RetCode    int         `json:"ret_code"`
	RetMsg     string      `json:"ret_msg"`
	RetContent interface{} `json:"ret_content"`
}

type Errno struct {
	Code    int
	Message string
}

// 错误码定义
var (
	Succeed = &Errno{Code: 0, Message: "success"}
	ErrParam = &Errno{Code:10001, Message: "param error"}
	ErrInsertFail = &Errno{Code: 10002, Message: "insert articles failed"}
	ErrDeleteFail = &Errno{Code: 10003, Message: "delete articles failed"}
	ErrUpdateFail = &Errno{Code: 10004, Message: "update articles failed"}
)

func (r *Response) SetErrorCode(e *Errno, errs ...error) {
	r.RetCode = e.Code
	r.RetMsg = e.Message
	if len(errs) > 0 {
		for _, err := range errs {
			r.RetContent = err.Error()
		}
	}
}
