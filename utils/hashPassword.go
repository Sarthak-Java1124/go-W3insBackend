package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "Error in hashing  the password"
	}
	return string(hashedPassword)
}

func VerifyPassword(hashedPassword string , password string) bool {
	verify := bcrypt.CompareHashAndPassword([]byte(hashedPassword) , []byte(password))
	return verify==nil
}
