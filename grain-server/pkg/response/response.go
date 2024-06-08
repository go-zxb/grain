package response

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-grain/grain/pkg/encrypt"
	"github.com/go-pay/gopay/pkg/xlog"
	"time"
)

var (
	ReqOk ErrCode = 2000
)

type ErrCode int

type IResponse interface {
	New() IResponse
	WithCode(code ErrCode) IResponse
	WithSuccess(success bool) IResponse
	WithMessage(message string) IResponse
	WithData(data interface{}) IResponse
	WithData2(data interface{}) IResponse
	WithTotal(total int64) IResponse
	WithPage(page int) IResponse
	WithPageSize(pageSize int) IResponse
	aesEncrypt(encryptKey string) error
	Success(ctx *gin.Context)
	Fail(ctx *gin.Context)
}

type Response struct {
	Code          ErrCode     `json:"code"`
	IsSuccess     bool        `json:"success"`
	Message       string      `json:"message"`
	IV            string      `json:"iv,omitempty"`
	Encrypted     bool        `json:"encrypted,omitempty"`
	EncryptedType string      `json:"encrypted_type,omitempty"`
	Data          interface{} `json:"data,omitempty"`
	Data2         interface{} `json:"data2,omitempty"`
	Time          int64       `json:"time,omitempty"`
	Total         int64       `json:"total,omitempty"`
	Page          int         `json:"page,omitempty"`
	PageSize      int         `json:"page_size,omitempty"`
}

func New() IResponse {
	return &Response{}
}

func (r *Response) New() IResponse {
	return &Response{}
}

func (r *Response) WithCode(code ErrCode) IResponse {
	r.Code = code
	return r
}

func (r *Response) WithSuccess(success bool) IResponse {
	r.IsSuccess = success
	return r
}

func (r *Response) WithMessage(message string) IResponse {
	r.Message = message
	return r
}

func (r *Response) WithData(data interface{}) IResponse {
	r.Data = data
	return r
}

func (r *Response) WithData2(data interface{}) IResponse {
	r.Data2 = data
	return r
}

func (r *Response) WithTotal(total int64) IResponse {
	r.Total = total
	return r
}

func (r *Response) WithPage(page int) IResponse {
	r.Page = page
	return r
}

func (r *Response) WithPageSize(pageSize int) IResponse {
	r.PageSize = pageSize
	return r
}

func (r *Response) Success(ctx *gin.Context) {
	if r.Data != nil {
		if ctx.GetString("encrypt") == "ok" {
			r.EncryptedType = ctx.GetString("encryptType")
			switch r.EncryptedType {
			case "aes":
				err := r.aesEncrypt(ctx.GetString("encryptKey"))
				if err != nil {
					r.WithMessage("加密数据失败").Fail(ctx)
					return
				}
			default:
				r.WithMessage("加密数据失败").Fail(ctx)
				return
			}
		}
	}

	if r.Code == 0 {
		r.Code = ReqOk
	}
	s := &Response{
		IsSuccess: true,
		Code:      r.Code,
		Message:   r.Message,
		Data:      r.Data,
		Data2:     r.Data2,
		IV:        r.IV,
		Encrypted: r.Encrypted,
		Total:     r.Total,
		PageSize:  r.PageSize,
		Page:      r.Page,
		Time:      time.Now().UnixMilli(),
	}
	ctx.JSON(200, s)
}

func (r *Response) Fail(ctx *gin.Context) {
	if r.Data != nil {
		if ctx.GetString("encrypt") == "ok" {
			var err error
			r.EncryptedType = ctx.GetString("encryptType")
			switch r.EncryptedType {
			case "aes":
				err = r.aesEncrypt(ctx.GetString("encryptKey"))
			default:
				err = errors.New("")
			}
			if err != nil {
				s := &Response{
					IsSuccess: false,
					Code:      r.Code,
					Message:   "数据加密传输失败",
					Time:      time.Now().UnixMilli(),
				}
				ctx.JSON(200, s)
				return
			}
		}
	}
	if r.Code == 0 {
		r.Code = 500
	}
	s := &Response{
		IsSuccess: false,
		Code:      r.Code,
		Message:   r.Message,
		Data:      r.Data,
		IV:        r.IV,
		Encrypted: r.Encrypted,
		Time:      time.Now().UnixMilli(),
	}
	ctx.JSON(200, s)
}

// 加密返回数据
func (r *Response) aesEncrypt(encryptKey string) error {
	marshal, err := json.Marshal(r.Data)
	if err != nil {
		r.WithData(nil)
		return err
	}
	key := []byte(encryptKey)
	ivStr := encrypt.GenerateAesIV(8)
	iv := []byte(ivStr)
	aesEncrypt, err := encrypt.AesEncrypt(marshal, key, iv)
	if err != nil {
		xlog.Info("加密失败", err)
		r.WithData(nil)
		return err
	}

	r.IV = ivStr
	r.Encrypted = true
	r.WithData(aesEncrypt)
	return nil
}
