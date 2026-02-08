package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(p string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(p), 14)
	return string(hash)
}

func CheckPassword(hash, p string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p)) == nil
}
