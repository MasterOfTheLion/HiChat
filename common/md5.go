package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

func Md5Encoder(code string) string {
	m := md5.New()
	io.WriteString(m, code)
	return hex.EncodeToString(m.Sum(nil))
}

func Md5StrToUpper(code string) string {
	return strings.ToUpper(Md5Encoder(code))
}

func SaltPassWord(pw string, salt string) string {
	saltPW := fmt.Sprintf("%s$%s", Md5Encoder(pw), salt)
	return saltPW
}

func CheckPassWord(rpw, salt, pw string) bool {
	return pw == SaltPassWord(rpw, salt)
}
