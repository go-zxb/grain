package uuidx

import (
	"crypto/rand"
	"github.com/go-pay/gopay/pkg/xlog"
	"github.com/gofrs/uuid"
	"github.com/oklog/ulid"
	"strings"
)

func UID() string {
	uid, _ := uuid.NewV4()
	return strings.ReplaceAll(uid.String(), "-", "")
}

func ULID() string {
	// 生成ULID
	id, err := ulid.New(ulid.Now(), ulid.Monotonic(rand.Reader, 0))
	if err != nil {
		xlog.Info("Error generating ULID:", err)
		return ""
	}
	return id.String()
}
