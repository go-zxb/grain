package jwtx

import (
	"errors"
	"github.com/go-grain/grain/pkg/response"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Jwt struct {
	res response.Response
}

// Claims 结构体，包含 ID 和 Username 字段，以及标准的声明
type Claims struct {
	AppID string `json:"appId,omitempty"`
	Uid   string `json:"uid,omitempty"`
	Type  uint   `json:"type,omitempty"` //登录方式
	Role  string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

func getSecretKey(secretKey string) []byte {
	return []byte(secretKey)
}

// GenerateToken 生成token
func (j Jwt) GenerateToken(uid, role, secretKey string, exp int64) (tokenString string, err error) {
	claim := Claims{
		Uid:  uid,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(exp))), // 过期时间在配置文件设置
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                       // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                       // 生效时间
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) // 使用HS256算法
	tokenString, err = token.SignedString([]byte(secretKey))
	return tokenString, err
}

// Secret 获取秘钥
func Secret(secretKey string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil // 这是我的secret
	}
}

// ParseToken 解析token 并验证是否有效
func (j Jwt) ParseToken(tokenStr, secretKey string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, Secret(secretKey))
	if err != nil {
		var errMsg string
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch {
			case ve.Errors&jwt.ValidationErrorMalformed != 0:
				errMsg = "无效令牌"
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				errMsg = "令牌已过期"
			case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
				errMsg = "令牌尚未激活"
			default:
				errMsg = "无法处理此令牌"
			}
		} else {
			errMsg = "无法处理此令牌"
		}
		return nil, errors.New(errMsg)
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("无法处理此令牌")
}
