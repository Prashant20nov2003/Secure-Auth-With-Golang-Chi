package function

import "golang.org/x/crypto/bcrypt"

func PasswordBcrypt(password string) []byte {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 11)
	return hashedPassword
}