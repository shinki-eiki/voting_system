/*
使用bcrypt加密和解密
*/
package util

import (
	"golang.org/x/crypto/bcrypt"
)

/* 随机生成n位的字符串 */
func RandomString(n int) string {
	// var letter=
	// a:=srtconv.Itoa
	return "114514"
}

/* 给字符串加密后返回 */
func EncryptedString(s string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	return string(res), err
}

/* 比对传入的密码是s否与数据库里的加密过的密码tar一致 */
func DecryptString(hashPassword, s string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(s))
}
