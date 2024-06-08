package encrypt

import (
	"github.com/go-pay/gopay/pkg/xlog"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		xlog.Info(err)
	}

	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		xlog.Info(err)
		return false
	}

	return true
}
